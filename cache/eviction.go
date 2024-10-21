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

// LRUEntry holds key for easy tracking in the list.
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

// RecordsAccess records the access of the keys in the queue and moves the key to the front of the queue...
func (lru *LRUCache) RecordAceess(key string) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element, ok := lru.items[key]; ok {
		lru.Queue.MoveToFront(element)
	} else {
		element := lru.Queue.PushFront(LRUEntry{key: key})
		lru.items[key] = element
	}
}
