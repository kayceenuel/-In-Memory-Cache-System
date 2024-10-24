package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestCache tests the basic cache functionality (Set/Get/Delete).
func TestCache(t *testing.T) {
	// Test Set and Get functionality.
	t.Run("Set and Get", func(t *testing.T) {
		cache := NewCache(3) // Set a reasonable max size

		// Add a key-value pair
		cache.Set("key1", "value1", 1*time.Second)
		value, exists := cache.Get("key1")
		if !exists {
			t.Errorf("Expected key 'key1' to exist")
		}
		if value != "value1" {
			t.Errorf("Expected value 'value1', got '%v'", value)
		}
	})

	// Test TTL expiration.
	t.Run("Expiration", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("key2", "value2", 100*time.Millisecond)

		// Wait for the key to expire
		time.Sleep(150 * time.Millisecond)

		// Now attempt to retrieve the key after the TTL
		_, exists := cache.Get("key2")
		if exists {
			t.Errorf("Expected key 'key2' to be expired")
		}
	})

	// Test manual deletion of keys.
	t.Run("Delete", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("key3", "value3", 1*time.Second)

		// Delete the key
		cache.Delete("key3")

		_, exists := cache.Get("key3")
		if exists {
			t.Errorf("Expected key 'key3' to be deleted")
		}
	})

	// Test LRU eviction behavior.
	t.Run("LRU Eviction", func(t *testing.T) {
		cache := NewCache(2) // Set max size to 2 for eviction testing

		// Add items to the cache
		cache.Set("key4", "value4", 1*time.Hour)
		cache.Set("key5", "value5", 1*time.Hour)

		// Access key4 to make it recently used
		cache.Get("key4")

		// Add another item to trigger eviction
		cache.Set("key6", "value6", 1*time.Hour)

		// Since "key5" is the least recently used, it should be evicted
		if _, exists := cache.Get("key5"); exists {
			t.Errorf("Expected key 'key5' to be evicted")
		}

		// Check if "key4" and "key6" still exist
		for _, key := range []string{"key4", "key6"} {
			if _, exists := cache.Get(key); !exists {
				t.Errorf("Expected key '%s' to still exist", key)
			}
		}
	})

	// Test multiple keys with long TTL.
	t.Run("Multiple Keys", func(t *testing.T) {
		cache := NewCache(3) // Max size set to 3

		// Set multiple keys
		cache.Set("keyA", "valueA", 5*time.Second)
		cache.Set("keyB", "valueB", 5*time.Second)
		cache.Set("keyC", "valueC", 5*time.Second)

		// Define expected key-value pairs for testing
		testCases := []struct {
			key      string
			expected string
		}{
			{"keyA", "valueA"},
			{"keyB", "valueB"},
			{"keyC", "valueC"},
		}

		// Test retrieval of each key
		for _, tc := range testCases {
			value, exists := cache.Get(tc.key)
			if !exists {
				t.Errorf("Expected key '%s' to exist", tc.key)
				continue
			}

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

	// Test concurrency by performing simultaneous read/write operations.
	t.Run("Concurrency", func(t *testing.T) {
		cache := NewCache(5)
		var wg sync.WaitGroup

		// Write values to the cache concurrently
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := "key" + strconv.Itoa(i) // Convert int to string of digits
				cache.Set(key, i, 10*time.Second)
			}(i)
		}

		// Read values concurrently
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := "key" + strconv.Itoa(i) // Convert int to string of digits
				cache.Get(key)
			}(i)
		}

		wg.Wait()
	})
}
