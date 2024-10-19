package cache

import (
	"sync"
)

// Cache is a simple in-memory key-value store. using a map and a mutex for concurrency.
type Cache struct {
	data map[string]interface{} // map stores key-value pairs
	mu   sync.RWMutex           // mutex for concurrency control
}

// NewCache creates a new cache instance.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Set adds a key-value pair to the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()         // Lock the mutex to ensure access to the map.
	defer c.mu.Unlock() // ensure the mutex is unlocked after the func returns
	c.data[key] = value // store value by key
}

// Get retrieves a value from the cache by key.
