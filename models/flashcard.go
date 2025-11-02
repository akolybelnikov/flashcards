package models

import "time"

type Flashcard struct {
	ID                   int       `json:"id"`
	Question             string    `json:"question"`
	Answer               string    `json:"answer"`
	QuestionLang         *string   `json:"question_lang,omitempty"`
	AnswerLang           *string   `json:"answer_lang,omitempty"`
	AITranslatedQuestion bool      `json:"ai_translated_question"`
	AITranslatedAnswer   bool      `json:"ai_translated_answer"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateFlashcardRequest struct {
	Question             string  `json:"question"`
	Answer               string  `json:"answer"`
	QuestionLang         *string `json:"question_lang,omitempty"`
	AnswerLang           *string `json:"answer_lang,omitempty"`
	AITranslatedQuestion *bool   `json:"ai_translated_question,omitempty"`
	AITranslatedAnswer   *bool   `json:"ai_translated_answer,omitempty"`
}

type UpdateFlashcardRequest struct {
	Question *string `json:"question,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

// GenerateTranslationRequest represents a request to generate an AI translation
type GenerateTranslationRequest struct {
	Content  string `json:"content"`
	FromLang string `json:"from_lang"`
	ToLang   string `json:"to_lang"`
}

// GenerateTranslationResponse represents the response from translation generation
type GenerateTranslationResponse struct {
	Translation string `json:"translation"`
	Cached      bool   `json:"cached"`
	CacheKey    string `json:"cache_key"`
}

// RandomFlashcardResponse represents the payload returned by the random flashcard endpoint.
// It contains the flashcard and an optional AI-generated hint or translation.
type RandomFlashcardResponse struct {
	Flashcard *Flashcard `json:"flashcard"`
	AIHint    *string    `json:"ai_hint,omitempty"`
}
