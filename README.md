# In-Memory Cache with LRU Eviction

This project is an in-memory key-value store with support for **TTL (Time-to-Live)** expiration and **LRU (Least Recently Used)** eviction. It provides a simple way to store, retrieve, and manage cached data with automatic removal of the least recently used items when the cache size limit is reached.

## Features

- **Set and Get Operations**: Store and retrieve values by key.
- **TTL Expiration**: Automatically remove entries after a specified time-to-live.s
- **LRU Eviction**: Evict the least recently used items when the cache exceeds its maximum size.
- **Concurrency Safe**: Supports concurrent access for both reads and writes.
  
## Usage

### 1. Create a Cache Instance
```go
cache := NewCache(5) // Max size of 5 items
```
### 2. Set a Key-Value Pair
```go
cache.Set("key", "value", 10*time.Second) // TTL of 10 seconds
```

### 3. Get a Value by Key
```go
value, exists := cache.Get("key")
if exists {
    fmt.Println("Value:", value)
}

```
### 4. Delete a Key-Value Pair
```go
cache.Delete("key")
```

### Testing
To run the tests, use the following command:
```bash
go test -v
```
    The tests cover:

- Basic set/get/delete operations
- Key expiration with TTL
- LRU eviction when the cache is full
- Concurrency handling