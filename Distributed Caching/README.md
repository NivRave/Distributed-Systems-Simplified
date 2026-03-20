# Distributed Caching (LRU & TTL)

This project simulates how **In-Memory Caching** drastically reduces database load and response latency by storing frequently accessed ("hot") data.

## What is Distributed Caching?

When a database is under heavy read pressure, querying the disk repeatedly for the same data is inefficient. A cache sits between the application and the database, storing data in extremely fast RAM.

### Core Concepts Simulated Here

1.  **Cache-Aside Pattern**: 
    The application checks the cache first.
    *   **Cache Hit**: Return data immediately (microseconds).
    *   **Cache Miss**: Query slow database (milliseconds), save the result in the cache, then return.

2.  **LRU (Least Recently Used) Eviction**:
    RAM is expensive and limited. When the cache reaches capacity, it must delete old data to make room for new data. The LRU policy evicts the item that hasn't been accessed for the longest time, assuming it's no longer "hot".

3.  **TTL (Time-To-Live) Expiration**:
    Data changes over time. To prevent the cache from serving stale (outdated) data forever, every item is given a TTL. Once that time passes, the cache deletes the item, forcing the next request to fetch fresh data from the database.

---

## How This Project Simulates the Theory

This project implements a custom `LRUCache` from scratch using standard Go libraries.

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **LRU Cache** | Designed using a `container/list` (Doubly Linked List) combined with a `map`. This ensures `O(1)` time complexity for reads, writes, and evictions. |
| **Slow Database** | A `Database` struct that artificially pauses via `time.Sleep(50ms)` on every read. |
| **TTL Expiration** | A `CacheItem` struct stores an `ExpiresAt` timestamp. The `Get` method checks if `time.Now() > ExpiresAt`. If so, it purges the item. |

### Key Files
*   `cache.go`: The core logic for the LRU cache (Set, Get, Evict, Map+LinkedList integration).
*   `database.go`: The mock database with artificial latency.
*   `simulation.go`: The orchestration engine running three distinct performance scenarios.

---

## How to Run

Ensure you have Go installed.

```bash
cd "Distributed Caching"
go run .
```

### Sample Output Explanation

**1. Scenario: The "Hot Key"**
```text
[Step 1] First time reading (DB Query required).
   [CACHE MISS] Key: trending_news:1 | Latency: 51.2ms (Fetched from DB)

[Step 2] Next 4 reads: Item is in cache. Loads instantly.
   [CACHE HIT]  Key: trending_news:1 | Latency: 0.1ms
   [CACHE HIT]  Key: trending_news:1 | Latency: 0.1ms
```
> **Observation**: The first read took ~50ms because it had to query the DB. Subsequent reads took less than a millisecond, a 500x speedup.

**2. Scenario: LRU Eviction**
```text
[Step 1] Insert 3 items into Cache (Fills Capacity).
   [Cache State] Ordered Most-Recently-Used -> Least-Recently-Used:
   Keys: [item:C item:B item:A]

[Step 2] Fetch a new item:D. This will evict the oldest (item:A).
   Keys: [item:D item:C item:B] <- Notice item:A is gone

[Step 3] Try to read item:A again.
   [CACHE MISS] Key: item:A        | Latency: 50.8ms (Fetched from DB)
```
> **Observation**: Because the cache capacity was 3, adding `item:D` forced the system to drop `item:A`. When we requested it again, it resulted in a Cache Miss.

**3. Scenario: TTL Expiration**
```text
[Step 1] Fetch price for the first time.
   [CACHE MISS] Key: stock_price   | Latency: 50.5ms (Fetched from DB)

[Step 2] Read immediately again (Data is fresh).
   [CACHE HIT]  Key: stock_price   | Latency: 0.1ms

[Step 3] Wait 1.1 seconds for TTL to expire...

[Step 4] Read again it after TTL. Cache should purge it and fetch anew.
   [CACHE MISS] Key: stock_price   | Latency: 51.1ms (Fetched from DB)
```
> **Observation**: Even though the item was technically in memory, the cache saw its `ExpiresAt` timestamp had passed, deleted it, and forced a fresh read from the DB to ensure data accuracy.
