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

	t.Run("Expiration", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key2", "value2", 100*time.Millisecond)

		time.Sleep(200 * time.Millisecond)

		_, exists := cache.Get("key2")
		if exists {
			t.Errorf("Expected key 'key2' to be expired")
		}
	})

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

	t.Run("Delete", func(t *testing.T) {
		cache := NewCache()
		cache.Set("key4", "value4", 1*time.Second)
		cache.Delete("key4")

		_, exists := cache.Get("key4")
		if exists {
			t.Errorf("Expected key 'key4' to be deleted")
		}
	})

	t.Run("Non-existent Key", func(t *testing.T) {
		cache := NewCache()
		_, exists := cache.Get("non_existent")
		if exists {
			t.Errorf("Expected non-existent key to return false")
		}
	})

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
