In-Memory Cache with LRU Eviction
This project is an in-memory key-value store with support for TTL (Time-to-Live) expiration and LRU (Least Recently Used) eviction. It provides a simple way to store, retrieve, and manage cached data with automatic removal of the least recently used items when the cache size limit is reached.

Features
Set and Get Operations: Store and retrieve values by key.
TTL Expiration: Automatically remove entries after a specified time-to-live.
LRU Eviction: Evict the least recently used items when the cache exceeds its maximum size.
Concurrency Safe: Supports concurrent access for both reads and writes.
Usage
Create a Cache Instance

go
Copy code
cache := NewCache(5) // Max size of 5 items
Set a Key-Value Pair

go
Copy code
cache.Set("key", "value", 10*time.Second) // TTL of 10 seconds
Get a Value by Key

go
Copy code
value, exists := cache.Get("key")
if exists {
    fmt.Println("Value:", value)
}
Delete a Key

go
Copy code
cache.Delete("key")
Testing
Run the following command to execute the unit tests:

bash
Copy code
go test -v
The tests cover:

Basic set/get/delete operations
Key expiration with TTL
LRU eviction when the cache is full
Concurrency handling