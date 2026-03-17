package cache

import (
	"sync"
	"time"

	"demo-auth-center/entity"
)

type cacheEntry struct {
	user      *entity.User
	expiresAt time.Time
}

type UserCache struct {
	mu  sync.RWMutex
	ttl time.Duration
	m   map[string]*cacheEntry // key: accountId
}

func NewUserCache(ttl time.Duration) *UserCache {
	return &UserCache{
		ttl: ttl,
		m:   make(map[string]*cacheEntry),
	}
}

func (c *UserCache) Get(accountID string) (*entity.User, bool) {
	c.mu.RLock()
	entry, ok := c.m[accountID]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.user, true
}

func (c *UserCache) Set(user *entity.User) {
	c.SetWithTTL(user, c.ttl)
}

func (c *UserCache) SetWithTTL(user *entity.User, ttl time.Duration) {
	c.mu.Lock()
	c.m[user.AccountID] = &cacheEntry{
		user:      user,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}

func (c *UserCache) Delete(accountID string) {
	c.mu.Lock()
	delete(c.m, accountID)
	c.mu.Unlock()
}
