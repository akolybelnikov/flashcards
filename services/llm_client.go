package services

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
func (c *OpenAIClient) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if c.llm == nil {
		return "", errors.New("LLM client not initialized")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Map language codes to full names for better prompt clarity
	langMap := map[string]string{
		"en": "English",
		"el": "Greek",
	}

	sourceLanguage := langMap[sourceLang]
	if sourceLanguage == "" {
		sourceLanguage = sourceLang
	}

	targetLanguage := langMap[targetLang]
	if targetLanguage == "" {
		targetLanguage = targetLang
	}

	prompt := fmt.Sprintf(
		"Translate the following text from %s to %s. Provide ONLY the translation, no explanations or additional text.\n\nText: %s",
		sourceLanguage,
		targetLanguage,
		text,
	)

	response, err := llms.GenerateFromSinglePrompt(ctx, c.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}

	return response, nil
}
