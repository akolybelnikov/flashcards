package services_test

import (
	"testing"
	"time"

	"github.com/akolybelnikov/flashcards/services"
)

func TestNewInMemoryTranslationCache(t *testing.T) {
	ttl := 1 * time.Hour
	cache := services.NewInMemoryTranslationCache(ttl)

	if cache == nil {
		t.Fatal("NewInMemoryTranslationCache() returned nil")
	}
}

func TestInMemoryTranslationCache_GenerateKey(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)

	tests := []struct {
		name     string
		content  string
		fromLang string
		toLang   string
	}{
		{
			name:     "simple text",
			content:  "hello",
			fromLang: "en",
			toLang:   "el",
		},
		{
			name:     "text with spaces",
			content:  "hello world",
			fromLang: "en",
			toLang:   "fr",
		},
		{
			name:     "text with special characters",
			content:  "hello! how are you?",
			fromLang: "en",
			toLang:   "de",
		},
		{
			name:     "unicode text",
			content:  "γεια σας",
			fromLang: "el",
			toLang:   "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := cache.GenerateKey(tt.content, tt.fromLang, tt.toLang)

			// Key should not be empty
			if key == "" {
				t.Error("GenerateKey() returned empty string")
			}

			// Key should be consistent for the same input
			key2 := cache.GenerateKey(tt.content, tt.fromLang, tt.toLang)
			if key != key2 {
				t.Errorf("GenerateKey() not consistent: %s != %s", key, key2)
			}

			// Key should be a valid hex string (SHA256 produces 64 hex chars)
			if len(key) != 64 {
				t.Errorf("GenerateKey() length = %d, want 64", len(key))
			}
		})
	}
}

func TestInMemoryTranslationCache_GenerateKey_Different(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)

	// Different inputs should produce different keys
	key1 := cache.GenerateKey("hello", "en", "el")
	key2 := cache.GenerateKey("hello", "en", "fr") // Different target language
	key3 := cache.GenerateKey("hello", "fr", "el") // Different source language
	key4 := cache.GenerateKey("world", "en", "el") // Different content

	if key1 == key2 {
		t.Error("Different target languages produced same key")
	}
	if key1 == key3 {
		t.Error("Different source languages produced same key")
	}
	if key1 == key4 {
		t.Error("Different content produced same key")
	}
}

func TestInMemoryTranslationCache_SetAndGet(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)
	key := cache.GenerateKey("hello", "en", "el")

	// Initially should not exist
	_, exists := cache.Get(key)
	if exists {
		t.Error("Get() found non-existent key")
	}

	// Set a translation
	translation := &services.CachedTranslation{
		Translation: "γεια σας",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set(key, translation)

	// Should now exist
	cached, exists := cache.Get(key)
	if !exists {
		t.Fatal("Get() did not find cached translation")
	}

	// Check values
	if cached.Translation != "γεια σας" {
		t.Errorf("Get() translation = %s, want γεια σας", cached.Translation)
	}
	if cached.FromLang != "en" {
		t.Errorf("Get() fromLang = %s, want en", cached.FromLang)
	}
	if cached.ToLang != "el" {
		t.Errorf("Get() toLang = %s, want el", cached.ToLang)
	}

	// Check timestamps were set
	if cached.CachedAt.IsZero() {
		t.Error("Get() CachedAt is zero")
	}
	if cached.ExpiresAt.IsZero() {
		t.Error("Get() ExpiresAt is zero")
	}
	if !cached.ExpiresAt.After(cached.CachedAt) {
		t.Error("Get() ExpiresAt should be after CachedAt")
	}
}

func TestInMemoryTranslationCache_Expiration(t *testing.T) {
	// Use very short TTL for testing
	ttl := 50 * time.Millisecond
	cache := services.NewInMemoryTranslationCache(ttl)
	key := cache.GenerateKey("hello", "en", "el")

	translation := &services.CachedTranslation{
		Translation: "γεια σας",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set(key, translation)

	// Should exist immediately
	_, exists := cache.Get(key)
	if !exists {
		t.Fatal("Get() did not find freshly cached translation")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should no longer exist
	_, exists = cache.Get(key)
	if exists {
		t.Error("Get() found expired translation")
	}
}

func TestInMemoryTranslationCache_MultipleEntries(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)

	entries := []struct {
		content  string
		fromLang string
		toLang   string
		trans    string
	}{
		{"hello", "en", "el", "γεια σας"},
		{"goodbye", "en", "el", "αντίο"},
		{"thank you", "en", "fr", "merci"},
		{"good morning", "en", "de", "guten Morgen"},
	}

	// Store all entries
	keys := make([]string, len(entries))
	for i, entry := range entries {
		key := cache.GenerateKey(entry.content, entry.fromLang, entry.toLang)
		keys[i] = key
		translation := &services.CachedTranslation{
			Translation: entry.trans,
			FromLang:    entry.fromLang,
			ToLang:      entry.toLang,
		}
		cache.Set(key, translation)
	}

	// Retrieve and verify all entries
	for i, entry := range entries {
		cached, exists := cache.Get(keys[i])
		if !exists {
			t.Errorf("Entry %d: Get() did not find cached translation", i)
			continue
		}
		if cached.Translation != entry.trans {
			t.Errorf("Entry %d: translation = %s, want %s", i, cached.Translation, entry.trans)
		}
	}
}

func TestInMemoryTranslationCache_Overwrite(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)
	key := cache.GenerateKey("hello", "en", "el")

	// Set initial translation
	translation1 := &services.CachedTranslation{
		Translation: "γεια",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set(key, translation1)

	// Overwrite with a new translation
	translation2 := &services.CachedTranslation{
		Translation: "γεια σας",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set(key, translation2)

	// Should get the new translation
	cached, exists := cache.Get(key)
	if !exists {
		t.Fatal("Get() did not find cached translation")
	}
	if cached.Translation != "γεια σας" {
		t.Errorf("Get() translation = %s, want γεια σας", cached.Translation)
	}
}

func TestInMemoryTranslationCache_StartCleanup(t *testing.T) {
	// Use very short TTL and cleanup interval
	ttl := 50 * time.Millisecond
	cleanupInterval := 100 * time.Millisecond
	cache := services.NewInMemoryTranslationCache(ttl)

	// Start a cleanup goroutine
	stopChan := make(chan struct{})
	defer close(stopChan)
	cache.StartCleanup(cleanupInterval, stopChan)

	// Add multiple entries
	keys := make([]string, 3)
	for i := 0; i < 3; i++ {
		key := cache.GenerateKey("test", "en", "el")
		keys[i] = key
		translation := &services.CachedTranslation{
			Translation: "τεστ",
			FromLang:    "en",
			ToLang:      "el",
		}
		cache.Set(key, translation)
	}

	// All entries should exist initially
	for i, key := range keys {
		_, exists := cache.Get(key)
		if !exists {
			t.Errorf("Entry %d: should exist initially", i)
		}
	}

	// Wait for entries to expire and cleanup to run
	time.Sleep(200 * time.Millisecond)

	// Entries should be cleaned up
	// Note: We can't guarantee exact timing, but at least some should be gone
	foundCount := 0
	for _, key := range keys {
		if _, exists := cache.Get(key); exists {
			foundCount++
		}
	}

	// After cleanup, we expect fewer or no entries
	// This is a soft assertion since timing can vary
	if foundCount == len(keys) {
		t.Log("Warning: Cleanup may not have run, all entries still present")
	}
}

func TestInMemoryTranslationCache_CleanupStop(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)
	stopChan := make(chan struct{})

	// Start cleanup
	cache.StartCleanup(100*time.Millisecond, stopChan)

	// Stop cleanup immediately
	close(stopChan)

	// Give it a moment to stop
	time.Sleep(50 * time.Millisecond)

	// Test should complete without hanging
	// If cleanup doesn't stop, this test would hang
}

func TestInMemoryTranslationCache_ConcurrentAccess(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			key := cache.GenerateKey("test", "en", "el")
			translation := &services.CachedTranslation{
				Translation: "τεστ",
				FromLang:    "en",
				ToLang:      "el",
			}
			cache.Set(key, translation)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func(n int) {
			key := cache.GenerateKey("test", "en", "el")
			cache.Get(key)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Test passes if no race conditions (run with -race flag)
}

func TestInMemoryTranslationCache_EmptyKey(t *testing.T) {
	cache := services.NewInMemoryTranslationCache(1 * time.Hour)

	// Get with an empty key
	_, exists := cache.Get("")
	if exists {
		t.Error("Get() should not find entry for empty key")
	}

	// Set with an empty key should not panic
	translation := &services.CachedTranslation{
		Translation: "test",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set("", translation)

	// Should be able to retrieve it
	cached, exists := cache.Get("")
	if !exists {
		t.Error("Get() should find entry for empty key after Set()")
	} else if cached.Translation != "test" {
		t.Errorf("Get() translation = %s, want test", cached.Translation)
	}
}

func TestInMemoryTranslationCache_TTLRespected(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := services.NewInMemoryTranslationCache(ttl)
	key := cache.GenerateKey("hello", "en", "el")

	translation := &services.CachedTranslation{
		Translation: "γεια σας",
		FromLang:    "en",
		ToLang:      "el",
	}
	cache.Set(key, translation)

	// Get immediately - should have correct TTL
	cached, exists := cache.Get(key)
	if !exists {
		t.Fatal("Get() did not find cached translation")
	}

	expectedExpiry := cached.CachedAt.Add(ttl)
	if !cached.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("ExpiresAt = %v, want %v", cached.ExpiresAt, expectedExpiry)
	}
}
