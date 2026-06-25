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

func NewCache(duration time.Duration) Cache {
	newCache := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: duration,
	}
	return *newCache
}

func (c *Cache) Add(key string, val []byte) int {
	return 0
}

func (c *Cache) Get() ([]byte, bool) {
	return make([]byte, 0), false
}
