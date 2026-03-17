package cache

import (
	"sync"
	"time"
)

type tokenEntry struct {
	accountID string
	expiresAt time.Time
}

// TokenCache maps accessToken → accountId with TTL
type TokenCache struct {
	mu  sync.RWMutex
	m   map[string]*tokenEntry
	ttl time.Duration
}

func NewTokenCache(ttl time.Duration) *TokenCache {
	return &TokenCache{ttl: ttl, m: make(map[string]*tokenEntry)}
}

func (c *TokenCache) Get(token string) (string, bool) {
	c.mu.RLock()
	e, ok := c.m[token]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return "", false
	}
	return e.accountID, true
}

func (c *TokenCache) Set(token, accountID string) {
	c.SetWithTTL(token, accountID, c.ttl)
}

func (c *TokenCache) SetWithTTL(token, accountID string, ttl time.Duration) {
	c.mu.Lock()
	c.m[token] = &tokenEntry{accountID: accountID, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *TokenCache) Delete(token string) {
	c.mu.Lock()
	delete(c.m, token)
	c.mu.Unlock()
}
