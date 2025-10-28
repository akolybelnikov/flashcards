package services

import (
	"errors"

	"github.com/akolybelnikov/flashcards/db"
	"github.com/akolybelnikov/flashcards/models"
)

type FlashcardService struct {
	repo db.FlashcardRepository
}

func NewFlashcardService(repo db.FlashcardRepository) *FlashcardService {
	return &FlashcardService{repo: repo}
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
