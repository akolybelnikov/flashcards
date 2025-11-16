package services_test

import (
	"testing"

	dbmocks "github.com/akolybelnikov/flashcards/db/mocks"
	"github.com/akolybelnikov/flashcards/models"
	"github.com/akolybelnikov/flashcards/services"
	"github.com/akolybelnikov/flashcards/services/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateFlashcardValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockLLM := mocks.NewMockLLMClient(ctrl)
	mockCache := mocks.NewMockTranslationCache(ctrl)

	questionLang := "en"
	answerLang := "el"

	// Set up expectation for Create
	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(&models.Flashcard{
			ID:           1,
			Question:     "hello",
			Answer:       "γεια σας",
			QuestionLang: &questionLang,
			AnswerLang:   &answerLang,
		}, nil)

	svc := services.NewFlashcardService(mockRepo, mockLLM, mockCache)

	// Both fields present - no translation needed
	fc, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "hello",
		Answer:   "γεια σας",
		FromLang: "en",
		ToLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Question != "hello" || fc.Answer != "γεια σας" {
		t.Fatalf("flashcard fields don't match input")
	}
}

func TestCreateFlashcardBothFieldsRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)

	// Service without LLM client
	svc := services.NewFlashcardService(mockRepo, nil, nil)

	// Both empty should fail
	_, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "",
		Answer:   "",
	})
	if err == nil {
		t.Fatalf("expected error when both fields are empty")
	}

	// Empty question with no LLM should fail
	_, err = svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "",
		Answer:   "γεια σας",
		FromLang: "el",
		ToLang:   "en",
	})
	if err == nil {
		t.Fatalf("expected error when LLM is nil and translation is needed")
	}
	if err.Error() != "AI translation not available: API key not configured" {
		t.Fatalf("expected AI translation error, got: %v", err)
	}

	// Empty answer with no LLM should fail
	_, err = svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "hello",
		Answer:   "",
		FromLang: "en",
		ToLang:   "el",
	})
	if err == nil {
		t.Fatalf("expected error when LLM is nil and translation is needed")
	}
	if err.Error() != "AI translation not available: API key not configured" {
		t.Fatalf("expected AI translation error, got: %v", err)
	}
}

func TestGenerateTranslation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockLLM := mocks.NewMockLLMClient(ctrl)
	mockCache := mocks.NewMockTranslationCache(ctrl)

	// Set up cache expectations - first call: cache miss
	cacheKey := "test-cache-key"
	mockCache.EXPECT().
		GenerateKey("hello", "en", "el").
		Return(cacheKey).
		Times(2) // Called twice

	mockCache.EXPECT().
		Get(cacheKey).
		Return(nil, false) // Cache miss

	// LLM will be called
	mockLLM.EXPECT().
		Translate(gomock.Any(), "hello", "en", "el").
		Return("γεια σας", nil)

	// Cache will store result
	mockCache.EXPECT().
		Set(cacheKey, gomock.Any())

	svc := services.NewFlashcardService(mockRepo, mockLLM, mockCache)

	// First call - cache miss
	resp, err := svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "hello",
		FromLang: "en",
		ToLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Translation != "γεια σας" {
		t.Fatalf("expected translation 'γεια σας', got '%s'", resp.Translation)
	}
	if resp.Cached {
		t.Fatalf("expected cached to be false on first call")
	}
	if resp.CacheKey == "" {
		t.Fatalf("expected cache key to be set")
	}

	// Second call - cache hit
	mockCache.EXPECT().
		Get(cacheKey).
		Return(&services.CachedTranslation{
			Translation: "γεια σας",
		}, true) // Cache hit

	resp2, err := svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "hello",
		FromLang: "en",
		ToLang:   "el",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.Translation != "γεια σας" {
		t.Fatalf("expected same translation from cache")
	}
	if !resp2.Cached {
		t.Fatalf("expected cached to be true on second call")
	}
	if resp2.CacheKey != resp.CacheKey {
		t.Fatalf("expected same cache key")
	}
}

func TestGenerateTranslationValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockLLM := mocks.NewMockLLMClient(ctrl)
	mockCache := mocks.NewMockTranslationCache(ctrl)

	svc := services.NewFlashcardService(mockRepo, mockLLM, mockCache)

	// Empty content
	_, err := svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "",
		FromLang: "en",
		ToLang:   "el",
	})
	if err == nil {
		t.Fatalf("expected error when content is empty")
	}

	// Empty from_lang
	_, err = svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "hello",
		FromLang: "",
		ToLang:   "el",
	})
	if err == nil {
		t.Fatalf("expected error when from_lang is empty")
	}

	// Empty to_lang
	_, err = svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "hello",
		FromLang: "en",
		ToLang:   "",
	})
	if err == nil {
		t.Fatalf("expected error when to_lang is empty")
	}
}

func TestGenerateTranslationWithoutLLMClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockCache := mocks.NewMockTranslationCache(ctrl)

	svc := services.NewFlashcardService(mockRepo, nil, mockCache)

	// Should fail when LLM client is nil
	_, err := svc.GenerateTranslation(&models.GenerateTranslationRequest{
		Content:  "hello",
		FromLang: "en",
		ToLang:   "el",
	})
	if err == nil {
		t.Fatalf("expected error when LLM client is nil")
	}
}

func TestUpdateFlashcardValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)

	svc := services.NewFlashcardService(mockRepo, nil, nil)

	// Both fields nil
	_, err := svc.UpdateFlashcard(1, &models.UpdateFlashcardRequest{})
	if err == nil {
		t.Fatalf("expected error when update request has no fields")
	}
}

func TestGetRandomFlashcardReturnsFlashcard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)

	// Set up expectation for GetRandom
	mockRepo.EXPECT().
		GetRandom().
		Return(&models.Flashcard{
			ID:       1,
			Question: "hello",
			Answer:   "γεια σας",
		}, nil)

	svc := services.NewFlashcardService(mockRepo, nil, nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)

	// Service without LLM client
	svc := services.NewFlashcardService(mockRepo, nil, nil)
	fc := &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σας"}

	hint := svc.GenerateAIHint(fc, "el")
	if hint != nil {
		t.Fatalf("expected nil hint when LLM client is nil, got '%s'", *hint)
	}
}

func TestGenerateAIHintWithLLMClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockLLM := mocks.NewMockLLMClient(ctrl)

	// Set up expectation for Translate call
	mockLLM.EXPECT().
		Translate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("A common Greek greeting", nil)

	// Service with mock LLM client
	svc := services.NewFlashcardService(mockRepo, mockLLM, nil)
	fc := &models.Flashcard{ID: 1, Question: "hello", Answer: "γεια σας"}

	hint := svc.GenerateAIHint(fc, "el")
	if hint == nil {
		t.Fatalf("expected hint when LLM client is available")
	}
	if *hint == "" {
		t.Fatalf("expected non-empty hint")
	}
}
