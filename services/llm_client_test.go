package services_test

import (
	"context"
	"testing"

	"github.com/akolybelnikov/flashcards/services"
)

func TestNewOpenAIClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty API key returns error",
			apiKey:  "",
			wantErr: true,
			errMsg:  "OpenAI API key is required",
		},
		{
			name:    "valid API key creates client",
			apiKey:  "sk-test-key-123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := services.NewOpenAIClient(tt.apiKey)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOpenAIClient() expected error, got nil")
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("NewOpenAIClient() error = %v, want %v", err.Error(), tt.errMsg)
				}
				if client != nil {
					t.Errorf("NewOpenAIClient() expected nil client on error, got %v", client)
				}
			} else {
				if err != nil {
					t.Errorf("NewOpenAIClient() unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("NewOpenAIClient() expected client, got nil")
				}
			}
		})
	}
}

func TestOpenAIClient_Translate_ValidationErrors(t *testing.T) {
	// Note: We can't fully test Translate without mocking the OpenAI API
	// But we can test the validation logic and error cases

	tests := []struct {
		name       string
		sourceLang string
		targetLang string
		text       string
		wantErr    bool
		errSubstr  string
	}{
		{
			name:       "invalid source language - too short",
			sourceLang: "e",
			targetLang: "el",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid source language",
		},
		{
			name:       "invalid source language - too long",
			sourceLang: "eng",
			targetLang: "el",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid source language",
		},
		{
			name:       "invalid source language - uppercase",
			sourceLang: "EN",
			targetLang: "el",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid source language",
		},
		{
			name:       "invalid source language - mixed case",
			sourceLang: "En",
			targetLang: "el",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid source language",
		},
		{
			name:       "invalid source language - numbers",
			sourceLang: "e1",
			targetLang: "el",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid source language",
		},
		{
			name:       "invalid target language - too short",
			sourceLang: "en",
			targetLang: "e",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid target language",
		},
		{
			name:       "invalid target language - too long",
			sourceLang: "en",
			targetLang: "ell",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid target language",
		},
		{
			name:       "invalid target language - uppercase",
			sourceLang: "en",
			targetLang: "EL",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid target language",
		},
		{
			name:       "invalid target language - special characters",
			sourceLang: "en",
			targetLang: "e-",
			text:       "hello",
			wantErr:    true,
			errSubstr:  "invalid target language",
		},
	}

	// Create a client with a dummy API key for validation testing
	// The actual OpenAI call won't succeed, but we can test validation
	client, err := services.NewOpenAIClient("sk-test-validation-key")
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := client.Translate(ctx, tt.text, tt.sourceLang, tt.targetLang)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Translate() expected error containing %q, got nil", tt.errSubstr)
					return
				}
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("Translate() error = %v, want error containing %q", err, tt.errSubstr)
				}
			}
		})
	}
}

func TestOpenAIClient_Translate_ValidLanguageCodes(t *testing.T) {
	// Test that valid ISO 639-1 language codes pass validation
	validCodes := []struct {
		name       string
		sourceLang string
		targetLang string
	}{
		{"English to Greek", "en", "el"},
		{"French to German", "fr", "de"},
		{"Spanish to Italian", "es", "it"},
		{"Japanese to Korean", "ja", "ko"},
		{"Chinese to Arabic", "zh", "ar"},
		{"Portuguese to Russian", "pt", "ru"},
		{"Dutch to Swedish", "nl", "sv"},
		{"Polish to Czech", "pl", "cs"},
	}

	client, err := services.NewOpenAIClient("sk-test-validation-key")
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	for _, tt := range validCodes {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// We expect the validation to pass, but the API call will fail with test key
			// The important thing is we don't get a validation error
			_, err := client.Translate(ctx, "test", tt.sourceLang, tt.targetLang)

			// We should get an API error, not a validation error
			if err != nil && (contains(err.Error(), "invalid source language") ||
				contains(err.Error(), "invalid target language")) {
				t.Errorf("Translate() got validation error for valid codes: %v", err)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
