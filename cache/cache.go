package cache

import (
	"fmt"
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

// Set adds or stores a key-value pair to the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()         // Lock the mutex to ensure access to the map.
	defer c.mu.Unlock() // ensure the mutex is unlocked after the func returns
	c.data[key] = value // store value by key
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key] // retrieve value by key
	return value, exists         // return value and boolean indicating if the key exists
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key) // delete key-value pair from map
}

func main() {
	cache := NewCache()

	//storing data in the cache
	cache.Set("name", "John Doe")
	cache.Set("age", 30)
	cache.Set("city", "New York")

	//Getting data from the cache
	if value, exists := cache.Get("name"); exists {
		fmt.Println("Name:", value)
	}

	//Deleting data from the cache
	cache.Delete("age")

	//checking if the key exists
	if _, exists := cache.Get("name"); !exists {
		fmt.Println("Key 'name' no longer exists")
	}
}
