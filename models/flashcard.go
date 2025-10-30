package models

import "time"

type Flashcard struct {
	ID        int       `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFlashcardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type UpdateFlashcardRequest struct {
	Question *string `json:"question,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

// RandomFlashcardResponse represents the payload returned by the random flashcard endpoint.
// It contains the flashcard and an optional AI-generated hint or translation.
type RandomFlashcardResponse struct {
	Flashcard *Flashcard `json:"flashcard"`
	AIHint    *string    `json:"ai_hint,omitempty"`
}
