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
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key] // retrieve value by key
	return value, ok         // return value and boolean indicating if the key exists
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key) // delete key-value pair from map
}
