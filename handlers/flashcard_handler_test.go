package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akolybelnikov/flashcards/models"
	"github.com/gorilla/mux"
)

// mockService implements the FlashcardServiceInterface for handler tests and returns deterministic values.
// Note: This could be replaced with gomock generated mocks, but inline mocks are simpler for handler tests.
type mockService struct{}

func (m *mockService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	now := time.Now()

	// Set defaults for pointers
	var questionLang, answerLang *string
	if req.QuestionLang != nil {
		questionLang = req.QuestionLang
	}
	if req.AnswerLang != nil {
		answerLang = req.AnswerLang
	}

	aiTranslatedQuestion := false
	if req.AITranslatedQuestion != nil {
		aiTranslatedQuestion = *req.AITranslatedQuestion
	}

	aiTranslatedAnswer := false
	if req.AITranslatedAnswer != nil {
		aiTranslatedAnswer = *req.AITranslatedAnswer
	}

	fc := &models.Flashcard{
		ID:                   1,
		Question:             req.Question,
		Answer:               req.Answer,
		QuestionLang:         questionLang,
		AnswerLang:           answerLang,
		AITranslatedQuestion: aiTranslatedQuestion,
		AITranslatedAnswer:   aiTranslatedAnswer,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	return fc, nil
}

func (m *mockService) GetAllFlashcards() ([]*models.Flashcard, error) {
	now := time.Now()
	return []*models.Flashcard{{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}}, nil
}

func (m *mockService) GetFlashcardByID(id int) (*models.Flashcard, error) {
	if id == 1 {
		now := time.Now()
		return &models.Flashcard{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}, nil
	}
	return nil, errors.New("flashcard with id not found")
}

func (m *mockService) UpdateFlashcard(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	if id != 1 {
		return nil, errors.New("flashcard with id not found")
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

func (m *mockService) DeleteFlashcard(id int) error {
	if id != 1 {
		return errors.New("flashcard with id not found")
	}
	return nil
}

func (m *mockService) GetRandomFlashcard() (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σασ", CreatedAt: now, UpdatedAt: now}, nil
}

func (m *mockService) GenerateAIHint(_ *models.Flashcard, _ string) *string {
	h := "hint"
	return &h
}

func (m *mockService) GenerateTranslation(req *models.GenerateTranslationRequest) (*models.GenerateTranslationResponse, error) {
	// Simple mock translation
	translation := "translated: " + req.Content
	if req.FromLang == "en" && req.ToLang == "el" && req.Content == "hello" {
		translation = "γεια σας"
	}

	return &models.GenerateTranslationResponse{
		Translation: translation,
		Cached:      false,
		CacheKey:    "mock-key-123",
	}, nil
}

func TestCreateFlashcardHandler(t *testing.T) {
	// use a mock service that provides deterministic results
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	// use gorilla/mux so path variables are parsed correctly
	r := mux.NewRouter()
	h.RegisterRoutes(r)

	payload := map[string]string{"question": "hello", "answer": "γεια σας"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/flashcards", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var flashcard models.Flashcard
	if err := json.NewDecoder(rr.Body).Decode(&flashcard); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if flashcard.Question != "hello" {
		t.Fatalf("expected question 'hello', got '%s'", flashcard.Question)
	}
	if flashcard.Answer != "γεια σας" {
		t.Fatalf("expected answer 'γεια σας', got '%s'", flashcard.Answer)
	}
}

func TestCreateFlashcardBothFieldsEmpty(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	payload := map[string]string{"question": "", "answer": ""}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/flashcards", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var errResp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if errResp["error"] != "Both question and answer must be provided" {
		t.Fatalf("unexpected error message: %s", errResp["error"])
	}
}

func TestGetAllFlashcardsHandler(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest("GET", "/flashcards", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

func TestGetFlashcardByIDNotFound(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest("GET", "/flashcards/2", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rr.Code)
	}
}

func TestUpdateFlashcardInvalidID(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	payload := map[string]string{"answer": "x"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/flashcards/abc", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	// gorilla/mux will return 404 for routes that don't match the {id:[0-9]+} pattern
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for invalid id pattern, got %d", rr.Code)
	}
}

func TestDeleteFlashcardHandler(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest("DELETE", "/flashcards/1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rr.Code)
	}
}

func TestGetRandomFlashcardHandler(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest("GET", "/flashcards/random", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp models.RandomFlashcardResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Flashcard == nil {
		t.Fatalf("expected flashcard in response")
	}
	if resp.AIHint == nil || *resp.AIHint != "hint" {
		t.Fatalf("expected ai_hint 'hint', got %v", resp.AIHint)
	}
}

func TestGenerateTranslationHandler(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	payload := map[string]string{
		"content":   "hello",
		"from_lang": "en",
		"to_lang":   "el",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/flashcards/translate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp models.GenerateTranslationResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Translation != "γεια σας" {
		t.Fatalf("expected translation 'γεια σας', got '%s'", resp.Translation)
	}
	if resp.CacheKey == "" {
		t.Fatalf("expected cache_key to be set")
	}
}

func TestGenerateTranslationHandlerEmptyContent(t *testing.T) {
	svc := &mockService{}
	h := NewFlashcardHandler(svc)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	payload := map[string]string{
		"content":   "",
		"from_lang": "en",
		"to_lang":   "el",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/flashcards/translate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
	}

	var errResp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if errResp["error"] != "Content cannot be empty" {
		t.Fatalf("unexpected error message: %s", errResp["error"])
	}
}
