package cache

import (
	"sync"
	"time"
)

// CacheItem struct to store the value and expiry time of a cache item
type CacheItem struct {
	value  interface{}
	expiry time.Time // expiry time of the cache item
}

// Cache is a simple in-memory key-value store. using a map and a mutex for concurrency.
type Cache struct {
	data map[string]CacheItem // map stores key-value pairs as CacheItems
	mu   sync.RWMutex         // mutex for concurrency control
}

// NewCache creates a new cache instance.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

// Set adds or stores a key-value pair to the cache.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheItem{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiry) {
		// Item has expired, delete it and return as if it doesn't exist
		c.mu.RUnlock()
		c.Delete(key)
		return nil, false
	}

	return item.value, true
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
