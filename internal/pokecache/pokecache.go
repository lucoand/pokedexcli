package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntries map[string]cacheEntry
	mtx          sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	entry, exists := c.CacheEntries[key]
	if !exists {
		return []byte{}, false
	}

	val := entry.val
	return val, true
}

func (c *Cache) Add(key string, val []byte) {
	if key == "" {
		return
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	entry := cacheEntry{}
	entry.createdAt = time.Now()
	entry.val = val
	c.CacheEntries[key] = entry
	return
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for true {
		tick := <-ticker.C
		for key, entry := range c.CacheEntries {
			if tick.Sub(entry.createdAt) > interval {
				delete(c.CacheEntries, key)
			}
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{}
	cache.CacheEntries = make(map[string]cacheEntry)
	go cache.reapLoop(interval)
	return &cache
}
