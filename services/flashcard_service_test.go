package services

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/akolybelnikov/flashcards/models"
)

// mockRepo is a small in-memory implementation of db.FlashcardRepository for tests.
type mockRepo struct{}

func (m *mockRepo) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{ID: 1, Question: req.Question, Answer: req.Answer, CreatedAt: now, UpdatedAt: now}, nil
}

func (m *mockRepo) GetAll() ([]*models.Flashcard, error) {
	now := time.Now()
	return []*models.Flashcard{{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}}, nil
}

func (m *mockRepo) GetByID(id int) (*models.Flashcard, error) {
	if id == 1 {
		now := time.Now()
		return &models.Flashcard{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepo) Update(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	if id != 1 {
		return nil, sql.ErrNoRows
	}
	q := "q"
	a := "a"
	if req.Question != nil {
		q = *req.Question
	}
	if req.Answer != nil {
		a = *req.Answer
	}
	now := time.Now()
	return &models.Flashcard{ID: id, Question: q, Answer: a, CreatedAt: now.Add(-time.Hour), UpdatedAt: now}, nil
}

func (m *mockRepo) Delete(id int) error {
	if id != 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *mockRepo) GetRandom() (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σασ", CreatedAt: now, UpdatedAt: now}, nil
}

func TestCreateFlashcardValidation(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{})

	// Missing question
	_, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{Question: "", Answer: "ans"})
	if err == nil {
		t.Fatalf("expected error when question is empty")
	}

	// Missing answer
	_, err = svc.CreateFlashcard(&models.CreateFlashcardRequest{Question: "q", Answer: ""})
	if err == nil {
		t.Fatalf("expected error when answer is empty")
	}
}

func TestUpdateFlashcardValidation(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{})

	// Both fields nil
	_, err := svc.UpdateFlashcard(1, &models.UpdateFlashcardRequest{})
	if err == nil {
		t.Fatalf("expected error when update request has no fields")
	}
}

func TestGetRandomFlashcardReturnsFlashcard(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{})

	fc, err := svc.GetRandomFlashcard()
	if err != nil {
		t.Fatalf("unexpected error getting random flashcard: %v", err)
	}
	if fc == nil {
		t.Fatalf("expected a flashcard, got nil")
	}
	if fc.Question != "hello" {
		t.Fatalf("expected question 'hello', got '%s'", fc.Question)
	}
}

func TestGenerateAIHintWithoutAPIKeyReturnsNil(t *testing.T) {
	// Ensure OPENAI_API_KEY is not set for this test to avoid network calls.
	old := os.Getenv("OPENAI_API_KEY")
	_ = os.Unsetenv("OPENAI_API_KEY")
	defer func() {
		_ = os.Setenv("OPENAI_API_KEY", old)
	}()

	svc := NewFlashcardService(&mockRepo{})
	fc := &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σασ"}

	hint := svc.GenerateAIHint(fc, "el")
	if hint != nil {
		t.Fatalf("expected nil hint when OPENAI_API_KEY is unset, got '%s'", *hint)
	}
}
