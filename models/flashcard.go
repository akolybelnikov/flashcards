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
