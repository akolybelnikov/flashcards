package services

//go:generate mockgen -destination=mocks/mock_flashcard_service.go -package=mocks github.com/akolybelnikov/flashcards/services FlashcardServiceInterface

import (
	"context"
	"errors"
	"log"

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
}

func NewFlashcardService(repo db.FlashcardRepository, llmClient LLMClient, translationCache TranslationCache) *FlashcardService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &FlashcardService{
		repo:             repo,
		llmClient:        llmClient,
		translationCache: translationCache,
	}
}

func (s *FlashcardService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	// Validate that both question and answer are provided
	if req.Question == "" || req.Answer == "" {
		return nil, errors.New("both question and answer must be provided")
	}

	// Create the flashcard with AI flags if provided (default to false)
	return s.repo.Create(req)
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

// GenerateTranslation generates a translation using AI, with caching to prevent redundant API calls
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

	// Generate cache key
	cacheKey := ""
	if s.translationCache != nil {
		cacheKey = s.translationCache.GenerateKey(req.Content, req.FromLang, req.ToLang)

		// Check cache first
		if cached, found := s.translationCache.Get(cacheKey); found {
			return &models.GenerateTranslationResponse{
				Translation: cached.Translation,
				Cached:      true,
				CacheKey:    cacheKey,
			}, nil
		}
	}

	// Cache miss - call LLM
	ctx := context.Background()
	translation, err := s.llmClient.Translate(ctx, req.Content, req.FromLang, req.ToLang)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if s.translationCache != nil {
		s.translationCache.Set(cacheKey, &CachedTranslation{
			Translation: translation,
			FromLang:    req.FromLang,
			ToLang:      req.ToLang,
		})
	}

	return &models.GenerateTranslationResponse{
		Translation: translation,
		Cached:      false,
		CacheKey:    cacheKey,
	}, nil
}
