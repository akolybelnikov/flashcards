package services

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
	CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, bool, string, error)
	GetAllFlashcards() ([]*models.Flashcard, error)
	GetFlashcardByID(id int) (*models.Flashcard, error)
	UpdateFlashcard(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error)
	DeleteFlashcard(id int) error
	GetRandomFlashcard() (*models.Flashcard, error)
	GenerateAIHint(flashcard *models.Flashcard, lang string) *string
}

type FlashcardService struct {
	repo      db.FlashcardRepository
	llmClient LLMClient
}

func NewFlashcardService(repo db.FlashcardRepository, llmClient LLMClient) *FlashcardService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &FlashcardService{
		repo:      repo,
		llmClient: llmClient,
	}
}

func (s *FlashcardService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, bool, string, error) {
	// Case 1: Both question and answer provided - no translation needed
	if req.Question != "" && req.Answer != "" {
		fc, err := s.repo.Create(req)
		return fc, false, "", err
	}

	if s.llmClient == nil {
		return nil, false, "", errors.New("AI translation not available: API key not configured")
	}

	translatedField := ""

	// Case 2: Only question provided - translate to answer
	if req.Question != "" && req.Answer == "" {
		translation, err := s.llmClient.Translate(context.Background(), req.Question, req.QuestionLang, req.AnswerLang)
		if err != nil {
			return nil, false, "", errors.New("failed to translate question to answer: " + err.Error())
		}

		req.Answer = translation
		translatedField = "answer"
	}

	// Case 3: Only answer provided - translate to question
	if req.Answer != "" && req.Question == "" {
		translation, err := s.llmClient.Translate(context.Background(), req.Answer, req.AnswerLang, req.QuestionLang)
		if err != nil {
			return nil, false, "", errors.New("failed to translate answer to question: " + err.Error())
		}

		req.Question = translation
		translatedField = "question"
	}

	flashcard, err := s.repo.Create(req)
	return flashcard, true, translatedField, err
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
