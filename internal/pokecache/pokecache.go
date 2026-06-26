package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func NewCache(duration time.Duration) *Cache {
	newCache := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: duration,
	}
	go newCache.reapLoop(duration)

	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return data.val, true
}

func (c *Cache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key := range c.cache {
			if time.Since(c.cache[key].createdAt) > duration {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
