package cache

import (
	"testing"
	"time"
)

// TestCache tests the cache functionality.
func TestCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key1", "value1", 1*time.Second)

		value, exists := cache.Get("key1")
		if !exists {
			t.Errorf("Expected key 'key1' to exist")
		}
		if value != "value1" {
			t.Errorf("Expected value 'value1', got '%v'", value)
		}
	})

	// Test expiration.
	t.Run("Expiration", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key2", "value2", 100*time.Millisecond)

		// Wait for the key to expire
		time.Sleep(150 * time.Millisecond)

		// Now attempt to retrieve the key after the TTL
		_, exists := cache.Get("key2")
		if exists {
			t.Errorf("Expected key 'key2' to be expired")
		}
	})

	// Test overwrite.
	t.Run("Overwrite", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key3", "value3", 1*time.Second)
		cache.Set("key3", "new_value3", 1*time.Second)

		value, exists := cache.Get("key3")
		if !exists {
			t.Errorf("Expected key 'key3' to exist")
		}
		if value != "new_value3" {
			t.Errorf("Expected value 'new_value3', got '%v'", value)
		}
	})

	// Test delete.
	t.Run("Delete", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key4", "value4", 1*time.Second)
		cache.Delete("key4")

		_, exists := cache.Get("key4")
		if exists {
			t.Errorf("Expected key 'key4' to be deleted")
		}
	})

	// Test non-existent key.
	t.Run("Non-existent Key", func(t *testing.T) {
		cache := NewCache()
		_, exists := cache.Get("non_existent")
		if exists {
			t.Errorf("Expected non-existent key to return false")
		}
	})

	// Test multiple keys
	t.Run("Multiple Keys", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key5", "value5", 1*time.Second)
		cache.Set("key6", "value6", 2*time.Second)
		cache.Set("key7", "value7", 3*time.Second)

		testCases := []struct {
			key      string
			expected string
		}{
			{"key5", "value5"},
			{"key6", "value6"},
			{"key7", "value7"},
		}

		for _, tc := range testCases {
			value, exists := cache.Get(tc.key)
			if !exists {
				t.Errorf("Expected key '%s' to exist", tc.key)
			}
			if value != tc.expected {
				t.Errorf("Expected value '%s', got '%v'", tc.expected, value)
			}
		}
	})
}

// TestLRU: test the LRU eviction policy
func TestLRU(t *testing.T) {
	cache := NewCache()             // Create a new cache
	evictionPolicy := NewLRUCache() // Create a New LRU Cache

	// Add items to cache and simulate access
	cache.Set("key1", "value1", 1*time.Hour)
	evictionPolicy.RecordAceess("key1")

	// Add another item
	cache.Set("key2", "value2", 1*time.Hour)
	evictionPolicy.RecordAceess("key2")

	cache.Set("key3", "value3", 1*time.Hour)
	evictionPolicy.RecordAceess("key3")

	// Exceed cache size
	evictedKey := evictionPolicy.RemoveEviction()

	// Check if the least recently used key is evicted
	if evictedKey != "key1" {
		t.Errorf("Expected key 'key1' (least recently used) to be evicted")
	}

	// Check if "key1" no longer exists in cache
	if _, exists := cache.Get("key1"); exists {
		t.Errorf("Expected key 'key1' to be evicted from cache")
	}
}
