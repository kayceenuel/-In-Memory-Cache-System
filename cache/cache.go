package cache

import (
	"fmt"
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
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheItem{ // store the cache item in the map
		value:  value,               // value of the cache item
		expiry: time.Now().Add(ttl), //expiry time of the cache item
	}
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()         // read lock
	defer c.mu.RUnlock() // read unlock

	item, exists := c.data[key] // retrieve item by key
	if !exists {
		return nil, false
	}

	cacheItem, ok := item.(CacheItem)
	if !ok {
		// If the item is not of type CacheItem, something went wron
		// We'll delete it and return as if it doesn't exist
		c.mu.RUnlock()
		c.Delete(key)
		c.mu.RLock()
		return nil, false
	}

	if time.Now().After(cacheItem.expiry) { // check if the item has expired
		// Item has expired
		c.mu.RUnlock()
		c.Delete(key)
		c.mu.RLock()
		return nil, false
	}

	return cacheItem.value, true
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key) // delete key-value pair from map
}

func main() {
	cache := NewCache()

	//storing data in the cache with a ttl of 1 hour
	cache.Set("name", "John Doe", 1*time.Hour)
	cache.Set("age", 30, 1*time.Hour)
	cache.Set("city", "New York", 1*time.Hour)

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
