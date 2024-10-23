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

		// Set multiple keys with a long TTL to avoid expiration during the test
		cache.Set("key5", "value5", 10*time.Second) // Increased TTL to 10 seconds
		cache.Set("key6", "value6", 10*time.Second)
		cache.Set("key7", "value7", 10*time.Second)

		// Define expected key-value pairs for testing
		testCases := []struct {
			key      string
			expected string
		}{
			{"key5", "value5"},
			{"key6", "value6"},
			{"key7", "value7"},
		}

		// Test retrieval of each key
		for _, tc := range testCases {
			value, exists := cache.Get(tc.key)

			// Log additional information to help debug why "key5" isn't found
			t.Logf("Checking key '%s': exists=%v, value=%v", tc.key, exists, value)

			if !exists {
				t.Errorf("Expected key '%s' to exist", tc.key)
				continue
			}

			// Convert value to string for comparison
			valueStr, ok := value.(string)
			if !ok {
				t.Errorf("Expected value for key '%s' to be a string, got %T", tc.key, value)
				continue
			}

			if valueStr != tc.expected {
				t.Errorf("Expected value '%s', got '%v'", tc.expected, valueStr)
			}
		}
	})

	// Test LRU
	t.Run("LRU Eviction", func(t *testing.T) {
		cache := NewCache() // Default maxSize is 2 for simplicity

		// Add items to cache and simulate access
		cache.Set("key1", "value1", 1*time.Hour)
		cache.Set("key2", "value2", 1*time.Hour)
		cache.Set("key3", "value3", 1*time.Hour) // This should trigger eviction

		// Check if the least recently used key "key1" is evicted
		if _, exists := cache.Get("key1"); exists {
			t.Errorf("Expected key 'key1' to be evicted")
		}

		// Check if other keys still exist
		for _, key := range []string{"key2", "key3"} {
			if _, exists := cache.Get(key); !exists {
				t.Errorf("Expected key '%s' to still exist in cache", key)
			}
		}

		// Add another key to exceed the cache size and evict another key
		cache.Set("key4", "value4", 1*time.Hour)
		if _, exists := cache.Get("key2"); exists {
			t.Errorf("Expected key 'key2' to be evicted")
		}
	})
}
