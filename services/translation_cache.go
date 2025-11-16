package services

//go:generate mockgen -destination=mocks/mock_translation_cache.go -package=mocks github.com/akolybelnikov/flashcards/services TranslationCache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	gc "github.com/patrickmn/go-cache"
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

// GoCacheTranslationCache is a TranslationCache implementation backed by github.com/patrickmn/go-cache
type GoCacheTranslationCache struct {
	cache *gc.Cache
	ttl   time.Duration
}

// NewGoCacheTranslationCache creates a new GoCacheTranslationCache.
// ttl: default expiration for entries. cleanupInterval: how often the internal janitor runs (passed to go-cache New).
func NewGoCacheTranslationCache(ttl, cleanupInterval time.Duration) *GoCacheTranslationCache {
	c := gc.New(ttl, cleanupInterval)
	return &GoCacheTranslationCache{
		cache: c,
		ttl:   ttl,
	}
}

// NewInMemoryTranslationCache maintains backward compatibility with the previous in-memory
// constructor used throughout the codebase and tests. It returns a TranslationCache with
// sensible defaults. The caller provides a TTL for entries; a cleanup interval defaults to 10 minutes.
func NewInMemoryTranslationCache(ttl time.Duration) TranslationCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return NewGoCacheTranslationCache(ttl, 10*time.Minute)
}

// Get retrieves a cached translation by key
// Returns nil, false if not found or expired
func (c *GoCacheTranslationCache) Get(key string) (*CachedTranslation, bool) {
	item, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	ct, ok := item.(*CachedTranslation)
	if !ok {
		// type mismatch; remove the corrupted entry
		c.cache.Delete(key)
		return nil, false
	}
	// go-cache guarantees expired items are not returned, so this is valid
	return ct, true
}

// Set stores a translation in the cache using configured TTL
func (c *GoCacheTranslationCache) Set(key string, translation *CachedTranslation) {
	translation.CachedAt = time.Now()
	translation.ExpiresAt = translation.CachedAt.Add(c.ttl)
	c.cache.Set(key, translation, c.ttl)
}

// GenerateKey creates a deterministic cache key from content and language codes
func (c *GoCacheTranslationCache) GenerateKey(content, fromLang, toLang string) string {
	data := fmt.Sprintf("%s:%s:%s", content, fromLang, toLang)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// StartCleanup starts a background goroutine that periodically forces deletion of expired entries.
// It will stop when stopChan is closed.
func (c *GoCacheTranslationCache) StartCleanup(interval time.Duration, stopChan <-chan struct{}) {
	// If an interval <= 0, rely on go-cache's internal janitor and do nothing.
	if interval <= 0 {
		return
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.cache.DeleteExpired()
			case <-stopChan:
				return
			}
		}
	}()
}
