package cache

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem stores the value and expiry time of a cache item
type CacheItem struct {
	value  interface{}
	expiry time.Time
}

// Cache is a simple in-memory key-value store with LRU eviction.
type Cache struct {
	data     map[string]*list.Element // map to store key-value pairs
	eviction *list.List               // doubly linked list to track LRU order
	mu       sync.RWMutex             // mutex for concurrency control
	maxSize  int                      // max size of cache
}

// Entry represents an entry in the eviction list
type Entry struct {
	key   string
	value CacheItem
}

// NewCache creates a new cache instance with a given max size.
func NewCache(maxSize int) *Cache {
	return &Cache{
		data:     make(map[string]*list.Element),
		eviction: list.New(),
		maxSize:  maxSize,
	}
}

// Set adds or stores a key-value pair to the cache with a TTL (time to live).
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If the key already exists and hasn't expired, update it
	if element, exists := c.data[key]; exists {
		entry := element.Value.(*Entry)
		if time.Now().After(entry.value.expiry) {
			c.eviction.Remove(element)
			delete(c.data, key)
		} else {
			c.eviction.MoveToFront(element)
			entry.value = CacheItem{value: value, expiry: time.Now().Add(ttl)}
			return
		}
	}

	// Add new item to cache
	item := CacheItem{value: value, expiry: time.Now().Add(ttl)}
	entry := &Entry{key: key, value: item}
	element := c.eviction.PushFront(entry)
	c.data[key] = element

	// Evict least recently used if cache exceeds size
	if c.eviction.Len() > c.maxSize {
		c.evict()
	}
}

// Get retrieves a value from the cache by key, removing it if expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, exists := c.data[key]
	if !exists {
		return nil, false
	}

	entry := element.Value.(*Entry)
	if time.Now().After(entry.value.expiry) {
		// If the item has expired, remove it immediately
		c.mu.RUnlock()
		c.Delete(key)
		c.mu.RLock()
		return nil, false
	}

	// Move to front (recently used)
	c.eviction.MoveToFront(element)
	return entry.value.value, true
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, exists := c.data[key]
	if !exists {
		return
	}

	c.eviction.Remove(element)
	delete(c.data, key)
}

// evict removes the least recently used (LRU) item from the cache.
func (c *Cache) evict() {
	for {
		element := c.eviction.Back()
		if element == nil {
			break
		}

		entry := element.Value.(*Entry)
		if time.Now().After(entry.value.expiry) {
			c.eviction.Remove(element)
			delete(c.data, entry.key)
			continue
		}

		// Remove LRU item
		c.eviction.Remove(element)
		delete(c.data, entry.key)
		break
	}
}
