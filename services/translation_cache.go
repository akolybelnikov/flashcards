package services

//go:generate mockgen -destination=mocks/mock_translation_cache.go -package=mocks github.com/akolybelnikov/flashcards/services TranslationCache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CachedTranslation represents a cached translation with metadata
type CachedTranslation struct {
	Translation string
	FromLang    string
	ToLang      string
	CachedAt    time.Time
	ExpiresAt   time.Time
}

// TranslationCache defines the interface for translation caching
type TranslationCache interface {
	Get(key string) (*CachedTranslation, bool)
	Set(key string, translation *CachedTranslation)
	GenerateKey(content, fromLang, toLang string) string
	StartCleanup(interval time.Duration, stopChan <-chan struct{})
}

// InMemoryTranslationCache implements TranslationCache using an in-memory map
type InMemoryTranslationCache struct {
	cache map[string]*CachedTranslation
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewInMemoryTranslationCache creates a new in-memory translation cache
func NewInMemoryTranslationCache(ttl time.Duration) *InMemoryTranslationCache {
	return &InMemoryTranslationCache{
		cache: make(map[string]*CachedTranslation),
		ttl:   ttl,
	}
}

// Get retrieves a cached translation by key
// Returns nil, false if not found or expired
func (c *InMemoryTranslationCache) Get(key string) (*CachedTranslation, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(cached.ExpiresAt) {
		return nil, false
	}

	return cached, true
}

// Set stores a translation in the cache
func (c *InMemoryTranslationCache) Set(key string, translation *CachedTranslation) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Set expiration time
	translation.CachedAt = time.Now()
	translation.ExpiresAt = translation.CachedAt.Add(c.ttl)

	c.cache[key] = translation
}

// GenerateKey creates a deterministic cache key from content and language codes
func (c *InMemoryTranslationCache) GenerateKey(content, fromLang, toLang string) string {
	// Create a consistent hash of the input parameters
	data := fmt.Sprintf("%s:%s:%s", content, fromLang, toLang)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// StartCleanup starts a background goroutine that periodically removes expired entries
func (c *InMemoryTranslationCache) StartCleanup(interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-stopChan:
				return
			}
		}
	}()
}

// cleanup removes expired entries from the cache
func (c *InMemoryTranslationCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	// Find expired keys
	for key, cached := range c.cache {
		if now.After(cached.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired keys
	for _, key := range expiredKeys {
		delete(c.cache, key)
	}
}
