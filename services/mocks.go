package services

import "context"

// MockLLMClient is a mock implementation of LLMClient for testing
type MockLLMClient struct {
	TranslateFunc func(ctx context.Context, text, sourceLang, targetLang string) (string, error)
}

func (m *MockLLMClient) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(ctx, text, sourceLang, targetLang)
	}

	// Default mock behavior - simple translations
	if sourceLang == "en" && targetLang == "el" {
		switch text {
		case "hello":
			return "γεια σας", nil
		case "goodbye":
			return "αντίο", nil
		default:
			return "μετάφραση", nil // "translation" in Greek
		}
	}

	if sourceLang == "el" && targetLang == "en" {
		switch text {
		case "γεια σας":
			return "hello", nil
		case "αντίο":
			return "goodbye", nil
		default:
			return "translation", nil
		}
	}

	return text + " (translated)", nil
}
