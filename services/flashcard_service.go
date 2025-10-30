package services

import (
	"context"
	"errors"
	"log"

	"github.com/akolybelnikov/flashcards/db"
	"github.com/akolybelnikov/flashcards/models"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
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
}

type FlashcardService struct {
	repo db.FlashcardRepository
	llm  *openai.LLM
}

func NewFlashcardService(repo db.FlashcardRepository) *FlashcardService {
	llm, err := openai.New()
	if err != nil {
		log.Printf("AI hint generation disabled: %v", err)
		return nil
	}
	return &FlashcardService{repo: repo, llm: llm}
}

func (s *FlashcardService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	if req.Question == "" {
		return nil, errors.New("question is required")
	}
	if req.Answer == "" {
		return nil, errors.New("answer is required")
	}

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

// GenerateAIHint attempts to generate a short hint using OpenAI. It returns nil if generation fails.
func (s *FlashcardService) GenerateAIHint(flashcard *models.Flashcard, lang string) *string {
	ctx := context.Background()
	prompt := "Provide a one-sentence hint or translation for the following flashcard."
	if lang != "" {
		prompt += " Target language: " + lang + "."
	}
	prompt += "\nQuestion: " + flashcard.Question + "\nAnswer: " + flashcard.Answer

	response, llmErr := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
	if llmErr != nil {
		log.Printf("AI hint generation failed: %v", llmErr)
		return nil
	}

	return &response
}
