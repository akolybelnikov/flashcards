package services

//go:generate mockgen -destination=mocks/mock_llm_client.go -package=mocks github.com/akolybelnikov/flashcards/services LLMClient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// LLMClient defines the interface for language model operations
type LLMClient interface {
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)
}

// OpenAIClient implements LLMClient using OpenAI
type OpenAIClient struct {
	llm *openai.LLM
}

// NewOpenAIClient creates a new OpenAI client with the provided API key
func NewOpenAIClient(apiKey string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	llm, err := openai.New(openai.WithToken(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return &OpenAIClient{llm: llm}, nil
}

// Translate translates text from source language to target language
// Language codes should follow ISO 639-1 standard (e.g., "en", "el", "fr")
func (c *OpenAIClient) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if c.llm == nil {
		return "", errors.New("LLM client not initialized")
	}

	// Validate language codes
	if err := validateLanguageCode(sourceLang); err != nil {
		return "", fmt.Errorf("invalid source language: %w", err)
	}
	if err := validateLanguageCode(targetLang); err != nil {
		return "", fmt.Errorf("invalid target language: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// OpenAI understands ISO 639-1 language codes directly
	prompt := fmt.Sprintf(
		"Translate the following text from %s to %s. Provide ONLY the translation, no explanations or additional text.\n\nText: %s",
		sourceLang,
		targetLang,
		text,
	)

	response, err := llms.GenerateFromSinglePrompt(ctx, c.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}

	return response, nil
}

// validateLanguageCode checks if a language code follows ISO 639-1 format (2 lowercase letters)
func validateLanguageCode(code string) error {
	if len(code) != 2 {
		return fmt.Errorf("language code must be 2 characters, got %d", len(code))
	}
	for _, r := range code {
		if r < 'a' || r > 'z' {
			return fmt.Errorf("language code must contain only lowercase letters")
		}
	}
	return nil
}
