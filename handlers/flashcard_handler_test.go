package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akolybelnikov/flashcards/handlers"
	"github.com/akolybelnikov/flashcards/models"
	"github.com/akolybelnikov/flashcards/services/mocks"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestCreateFlashcardHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	// Set up expectation for CreateFlashcard
	mockSvc.EXPECT().
		CreateFlashcard(gomock.Any()).
		DoAndReturn(func(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
			now := time.Now()
			return &models.Flashcard{
				ID:                   1,
				Question:             req.Question,
				Answer:               req.Answer,
				QuestionLang:         req.QuestionLang,
				AnswerLang:           req.AnswerLang,
				AITranslatedQuestion: false,
				AITranslatedAnswer:   false,
				CreatedAt:            now,
				UpdatedAt:            now,
			}, nil
		})

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	now := time.Now()
	mockSvc.EXPECT().
		GetAllFlashcards().
		Return([]*models.Flashcard{{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}}, nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	mockSvc.EXPECT().
		GetFlashcardByID(2).
		Return(nil, errors.New("flashcard with id not found"))

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	mockSvc.EXPECT().
		DeleteFlashcard(1).
		Return(nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	now := time.Now()
	hint := "hint"

	mockSvc.EXPECT().
		GetRandomFlashcard().
		Return(&models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σας", CreatedAt: now, UpdatedAt: now}, nil)

	mockSvc.EXPECT().
		GenerateAIHint(gomock.Any(), gomock.Any()).
		Return(&hint)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

	mockSvc.EXPECT().
		GenerateTranslation(gomock.Any()).
		DoAndReturn(func(req *models.GenerateTranslationRequest) (*models.GenerateTranslationResponse, error) {
			translation := "γεια σας"
			if req.Content == "hello" && req.FromLang == "en" && req.ToLang == "el" {
				return &models.GenerateTranslationResponse{
					Translation: translation,
					Cached:      false,
					CacheKey:    "mock-key-123",
				}, nil
			}
			return &models.GenerateTranslationResponse{
				Translation: "translated: " + req.Content,
				Cached:      false,
				CacheKey:    "mock-key-123",
			}, nil
		})

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockFlashcardServiceInterface(ctrl)
	h := handlers.NewFlashcardHandler(mockSvc)

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
