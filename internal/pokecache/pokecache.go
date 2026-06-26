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
	go newCache.reapLoop()

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
	container := make([]byte, 0)
	container = append(container, data.val)
	return container, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.cache)
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
