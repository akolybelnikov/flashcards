package services

import (
	"database/sql"
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
	mockLLM := &MockLLMClient{}
	svc := NewFlashcardService(&mockRepo{}, mockLLM)

	// Both fields present - no translation needed
	fc, aiUsed, field, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question:     "hello",
		Answer:       "γεια σας",
		QuestionLang: "en",
		AnswerLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if aiUsed {
		t.Fatalf("expected aiUsed to be false when both fields provided")
	}
	if field != "" {
		t.Fatalf("expected empty translatedField when both fields provided")
	}
	if fc.Question != "hello" || fc.Answer != "γεια σας" {
		t.Fatalf("flashcard fields don't match input")
	}
}

func TestCreateFlashcardWithTranslation(t *testing.T) {
	mockLLM := &MockLLMClient{}
	svc := NewFlashcardService(&mockRepo{}, mockLLM)

	// Only the question provided - should translate to answer
	fc, aiUsed, field, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question:     "hello",
		Answer:       "",
		QuestionLang: "en",
		AnswerLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !aiUsed {
		t.Fatalf("expected aiUsed to be true")
	}
	if field != "answer" {
		t.Fatalf("expected translatedField to be 'answer', got '%s'", field)
	}
	if fc.Answer == "" {
		t.Fatalf("expected answer to be translated")
	}

	// Only answer provided - should translate to question
	fc, aiUsed, field, err = svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question:     "",
		Answer:       "γεια σας",
		QuestionLang: "en",
		AnswerLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !aiUsed {
		t.Fatalf("expected aiUsed to be true")
	}
	if field != "question" {
		t.Fatalf("expected translatedField to be 'question', got '%s'", field)
	}
	if fc.Question == "" {
		t.Fatalf("expected question to be translated")
	}
}

func TestCreateFlashcardWithoutLLMClient(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{}, nil)

	// Should fail when translation is needed but no LLM client
	_, _, _, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question:     "hello",
		Answer:       "",
		QuestionLang: "en",
		AnswerLang:   "el",
	})
	if err == nil {
		t.Fatalf("expected error when LLM client is nil and translation needed")
	}
}

func TestUpdateFlashcardValidation(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{}, nil)

	// Both fields nil
	_, err := svc.UpdateFlashcard(1, &models.UpdateFlashcardRequest{})
	if err == nil {
		t.Fatalf("expected error when update request has no fields")
	}
}

func TestGetRandomFlashcardReturnsFlashcard(t *testing.T) {
	svc := NewFlashcardService(&mockRepo{}, nil)

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

func TestGenerateAIHintWithoutLLMClientReturnsNil(t *testing.T) {
	// Service without LLM client
	svc := NewFlashcardService(&mockRepo{}, nil)
	fc := &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σας"}

	hint := svc.GenerateAIHint(fc, "el")
	if hint != nil {
		t.Fatalf("expected nil hint when LLM client is nil, got '%s'", *hint)
	}
}

func TestGenerateAIHintWithLLMClient(t *testing.T) {
	// Service with mock LLM client
	mockLLM := &MockLLMClient{}
	svc := NewFlashcardService(&mockRepo{}, mockLLM)
	fc := &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σας"}

	hint := svc.GenerateAIHint(fc, "el")
	if hint == nil {
		t.Fatalf("expected hint when LLM client is available")
	}
	if *hint == "" {
		t.Fatalf("expected non-empty hint")
	}
}
