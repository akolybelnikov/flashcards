package services

//go:generate mockgen -destination=mocks/mock_flashcard_service.go -package=mocks github.com/akolybelnikov/flashcards/services FlashcardServiceInterface

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/singleflight"

	"github.com/akolybelnikov/flashcards/db"
	"github.com/akolybelnikov/flashcards/models"
)

// FlashcardServiceInterface defines the methods the handlers depend on. This allows tests to
// provide a mock service implementation without depending on the concrete type.
type FlashcardServiceInterface interface {
	CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error)
	GetAllFlashcards() ([]*models.Flashcard, error)
	GetFlashcardByID(id int) (*models.Flashcard, error)
	UpdateFlashcard(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error)
	DeleteFlashcard(id int) error
	GetRandomFlashcard() (*models.Flashcard, error)
	GenerateAIHint(flashcard *models.Flashcard, lang string) *string
	GenerateTranslation(req *models.GenerateTranslationRequest) (*models.GenerateTranslationResponse, error)
}

type FlashcardService struct {
	repo             db.FlashcardRepository
	llmClient        LLMClient
	translationCache TranslationCache

	// translateGroup deduplicates concurrent translation requests for the same key.
	translateGroup *singleflight.Group
}

func NewFlashcardService(repo db.FlashcardRepository, llmClient LLMClient, translationCache TranslationCache) *FlashcardService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &FlashcardService{
		repo:             repo,
		llmClient:        llmClient,
		translationCache: translationCache,
		translateGroup:   &singleflight.Group{},
	}
}

func (s *FlashcardService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	// Validate: at least one of question or answer must be provided
	if req.Question == "" && req.Answer == "" {
		return nil, errors.New("at least one of question or answer must be provided")
	}

	// Track which fields were AI-generated
	aiTranslatedQuestion := false
	aiTranslatedAnswer := false

	// Case 1: Both fields provided - no translation needed
	if req.Question != "" && req.Answer != "" {
		// Just create the flashcard
		return s.repo.Create(req)
	}

	// Case 2: Translation needed - ensure LLM client is available
	if s.llmClient == nil {
		return nil, errors.New("AI translation not available: API key not configured")
	}

	// Case 3: Only question provided - translate to answer
	if req.Question != "" && req.Answer == "" {
		translation, err := s.generateTranslationWithCache(req.Question, req.FromLang, req.ToLang)
		if err != nil {
			return nil, errors.New("failed to generate answer translation: " + err.Error())
		}
		req.Answer = translation
		aiTranslatedAnswer = true
	}

	// Case 4: Only answer provided - translate to question
	if req.Answer != "" && req.Question == "" {
		translation, err := s.generateTranslationWithCache(req.Answer, req.FromLang, req.ToLang)
		if err != nil {
			return nil, errors.New("failed to generate question translation: " + err.Error())
		}
		req.Question = translation
		aiTranslatedQuestion = true
	}

	// Create the flashcard (now both fields are populated)
	flashcard, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Update AI flags in the returned flashcard
	flashcard.AITranslatedQuestion = aiTranslatedQuestion
	flashcard.AITranslatedAnswer = aiTranslatedAnswer

	return flashcard, nil
}

// generateTranslationWithCache is a helper that uses the cache+singleflight translation flow
func (s *FlashcardService) generateTranslationWithCache(content, fromLang, toLang string) (string, error) {
	resp, err := s.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  content,
		FromLang: fromLang,
		ToLang:   toLang,
	})
	if err != nil {
		return "", err
	}
	return resp.Translation, nil
}

func (s *FlashcardService) GetAllFlashcards() ([]*models.Flashcard, error) {
	return s.repo.GetAll()
}

func (s *FlashcardService) GetFlashcardByID(id int) (*models.Flashcard, error) {
	return s.repo.GetByID(id)
}

func (s *FlashcardService) UpdateFlashcard(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	if req.Question == nil && req.Answer == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	return s.repo.Update(id, req)
}

func (s *FlashcardService) DeleteFlashcard(id int) error {
	return s.repo.Delete(id)
}

func (s *FlashcardService) GetRandomFlashcard() (*models.Flashcard, error) {
	return s.repo.GetRandom()
}

// GenerateAIHint attempts to generate a short hint using OpenAI. It returns nil if the generation fails
// or if the OpenAI client was not initialized.
func (s *FlashcardService) GenerateAIHint(flashcard *models.Flashcard, lang string) *string {
	if s == nil || s.llmClient == nil {
		log.Printf("AI hint generation not available: llm not initialized")
		return nil
	}

	// For now, just use translation as a hint
	// In the future this could be expanded to generate more sophisticated hints
	ctx := context.Background()

	// Determine source and target language based on lang parameter
	sourceLang := "en"
	targetLang := lang
	if targetLang == "" {
		targetLang = "el" // default to Greek
	}

	hint, err := s.llmClient.Translate(ctx, flashcard.Question, sourceLang, targetLang)
	if err != nil {
		log.Printf("AI hint generation failed: %v", err)
		return nil
	}

	return &hint
}

func (s *FlashcardService) getTranslation(term, sourceLang, targetLang string) (string, error) {
	ctx := context.Background()
	return s.llmClient.Translate(ctx, term, sourceLang, targetLang)
}

// GenerateTranslation generates a translation using AI, with caching and singleflight deduplication
// to prevent redundant API calls for concurrent identical requests.
func (s *FlashcardService) GenerateTranslation(req *models.GenerateTranslationRequest) (*models.GenerateTranslationResponse, error) {
	// Validate input
	if req.Content == "" {
		return nil, errors.New("content cannot be empty")
	}
	if req.FromLang == "" || req.ToLang == "" {
		return nil, errors.New("both from_lang and to_lang must be provided")
	}

	// Check if LLM client is available
	if s.llmClient == nil {
		return nil, errors.New("AI translation not available: API key not configured")
	}

	// Generate a cache key and check cache first if available
	cacheKey := ""
	if s.translationCache != nil {
		cacheKey = s.translationCache.GenerateKey(req.Content, req.FromLang, req.ToLang)
		if cached, found := s.translationCache.Get(cacheKey); found {
			return &models.GenerateTranslationResponse{
				Translation: cached.Translation,
				Cached:      true,
				CacheKey:    cacheKey,
			}, nil
		}
	}

	// Build singleflight key (use cacheKey when available, otherwise deterministic composite key)
	sfKey := cacheKey
	if sfKey == "" {
		sfKey = req.Content + "|" + req.FromLang + "->" + req.ToLang
	}

	// Ensure translateGroup exists
	if s.translateGroup == nil {
		s.translateGroup = &singleflight.Group{}
	}

	// Use singleflight to deduplicate concurrent calls
	v, err, _ := s.translateGroup.Do(sfKey, func() (interface{}, error) {
		ctx := context.Background()
		translation, err := s.llmClient.Translate(ctx, req.Content, req.FromLang, req.ToLang)
		if err != nil {
			return "", err
		}
		// Store in cache if available, and we have a cache key
		if s.translationCache != nil && cacheKey != "" {
			s.translationCache.Set(cacheKey, &CachedTranslation{
				Translation: translation,
				FromLang:    req.FromLang,
				ToLang:      req.ToLang,
			})
		}
		return translation, nil
	})
	if err != nil {
		return nil, err
	}

	translation, _ := v.(string)

	// Return the translation obtained from the LLM; if a caller wants a cached result, subsequent
	// calls will hit the cache at the top of this function.
	return &models.GenerateTranslationResponse{
		Translation: translation,
		Cached:      false,
		CacheKey:    cacheKey,
	}, nil
}
