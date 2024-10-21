package cache

import (
	"container/list"
	"sync"
)

// Evictionpolicy is an interface that defines the methods for an eviction policy
type Evictionpolicy interface {
	RecordAceess(key string) //record the access of the key
	RemoveEviction() string  // return the key to evict
}

// LRU is a struct that implements the Evictionpolicy interface
type LRUCache struct {
	Queue *list.List               // queue to store the keys in the order of their access
	items map[string]*list.Element // map to store the keys and their positions in the queue
	mu    sync.Mutex
}

//LRUEntry holds key for easy tracking in the list.
type LRUEntry struct {
	key string
}

// NewLRUCache creates a new LRUCache
func NewLRUCache() *LRUCache {
	return &LRUCache{
		Queue: list.New(),
		items: make(map[string]*list.Element),
	}
}
