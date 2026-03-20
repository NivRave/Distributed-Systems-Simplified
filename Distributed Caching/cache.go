package main

import (
    "container/list"
    "sync"
    "time"
)

// CacheItem stores the KV pair alongside its expiration time.
type CacheItem struct {
    Key       string
    Value     string
    ExpiresAt time.Time
}

// LRUCache implements a thread-safe cache with capacity bounds and TTL.
type LRUCache struct {
    capacity  int
    items     map[string]*list.Element
    evictList *list.List
    mu        sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity:  capacity,
        items:     make(map[string]*list.Element),
        evictList: list.New(),
    }
}

// Get returns the value and a boolean indicating Hit/Miss.
// If the item exists but its TTL has expired, it is purged and returns a Miss.
func (c *LRUCache) Get(key string) (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if el, ok := c.items[key]; ok {
        item := el.Value.(*CacheItem)
        
        // Check TTL Expiration
        if time.Now().After(item.ExpiresAt) {
            c.evictList.Remove(el)
            delete(c.items, key)
            return "", false // Cache Miss (Expired)
        }
        
        // Update LRU because it was recently accessed
        c.evictList.MoveToFront(el)
        return item.Value, true // Cache Hit
    }
    
    return "", false // Cache Miss (Not Found)
}

// Set adds or updates an item, refreshing its TTL and pushing it to the front of LRU.
func (c *LRUCache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 1. Update if exists
    if el, ok := c.items[key]; ok {
        c.evictList.MoveToFront(el) // Freshly used
        item := el.Value.(*CacheItem)
        item.Value = value
        item.ExpiresAt = time.Now().Add(ttl)
        return
    }

    // 2. Add new item
    item := &CacheItem{
        Key:       key,
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
    el := c.evictList.PushFront(item)
    c.items[key] = el

    // 3. Evict oldest if we exceed capacity
    if c.evictList.Len() > c.capacity {
        c.removeOldest()
    }
}

// removeOldest pops the tail element from the linked list.
func (c *LRUCache) removeOldest() {
    el := c.evictList.Back()
    if el != nil {
        c.evictList.Remove(el)
        item := el.Value.(*CacheItem)
        delete(c.items, item.Key) // Free memory from map
    }
}

// Keys returns an ordered list of keys (Most recently used first).
func (c *LRUCache) Keys() []string {
    c.mu.Lock()
    defer c.mu.Unlock()
    var keys []string
    for el := c.evictList.Front(); el != nil; el = el.Next() {
        keys = append(keys, el.Value.(*CacheItem).Key)
    }
    return keys
}
