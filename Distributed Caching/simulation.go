package main

import (
	"fmt"
	"time"
)

type Simulation struct {
	Cache *LRUCache
	DB    *Database
}

func NewSimulation() *Simulation {
	return &Simulation{
		Cache: NewLRUCache(3), // Small capacity to show LRU eviction easily
		DB:    NewDatabase(),
	}
}

// FetchData implements the Cache-Aside pattern.
func (s *Simulation) FetchData(key string, ttl time.Duration) string {
	start := time.Now()
	
	// 1. Ask Cache
	val, hit := s.Cache.Get(key)
	if hit {
		latency := time.Since(start)
		fmt.Printf("   \033[36m[CACHE HIT]\033[0m  Key: %-15s | Latency: %v\n", key, latency)
		return val
	}
	
	// 2. Cache Miss -> Fallback to slow DB
	val = s.DB.Read(key)
	
	// 3. Save to Cache for next time
	s.Cache.Set(key, val, ttl)
	
	latency := time.Since(start)
	fmt.Printf("   \033[33m[CACHE MISS]\033[0m Key: %-14s | Latency: %v (Fetched from DB)\n", key, latency)
	
	return val
}

func (s *Simulation) Run() {
	fmt.Println("=== Distributed Caching Simulation ===")
	
	s.RunScenario1_HotKey()
	s.RunScenario2_LRU()
	s.RunScenario3_TTL()
}

func (s *Simulation) RunScenario1_HotKey() {
	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 1: The Concept of a 'Hot Key'")
	fmt.Println("Goal: Show dramatic latency reduction on repeated reads")
	fmt.Println("-----------------------------------------------------")
	
	key := "trending_news:1"
	ttl := 1 * time.Minute // Doesn't expire during this scenario
	
	fmt.Println("\n[Step 1] First time reading (DB Query required).")
	s.FetchData(key, ttl)
	
	fmt.Println("\n[Step 2] Next 4 reads: Item is in cache. Loads instantly.")
	for i := 0; i < 4; i++ {
		s.FetchData(key, ttl)
	}
}

func (s *Simulation) RunScenario2_LRU() {
	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 2: LRU Eviction (Memory Management)")
	fmt.Println("Goal: Show that cache purges old items when full (Capacity = 3)")
	fmt.Println("-----------------------------------------------------")
	
	// Reset fresh cache
	s.Cache = NewLRUCache(3)
	ttl := 1 * time.Minute
	
	fmt.Println("\n[Step 1] Insert 3 items into Cache (Fills Capacity).")
	s.FetchData("item:A", ttl)
	s.FetchData("item:B", ttl)
	s.FetchData("item:C", ttl)
	
	fmt.Println("\n[Cache State] Ordered Most-Recently-Used -> Least-Recently-Used:")
	fmt.Println("   Keys:", s.Cache.Keys())
	
	fmt.Println("\n[Step 2] Fetch a new item:D. This will evict the oldest (item:A).")
	s.FetchData("item:D", ttl)
	fmt.Println("   Keys:", s.Cache.Keys(), "<- Notice item:A is gone")
	
	fmt.Println("\n[Step 3] Try to read item:A again.")
	s.FetchData("item:A", ttl) // Will Miss!
}

func (s *Simulation) RunScenario3_TTL() {
	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 3: TTL (Time-To-Live) Expiration")
	fmt.Println("Goal: Ensure data freshness by expiring old data")
	fmt.Println("-----------------------------------------------------")
	
	s.Cache = NewLRUCache(3)
	key := "stock_price:TSLA"
	ttl := 1 * time.Second // Short TTL
	
	fmt.Println("\n[Step 1] Fetch price for the first time.")
	s.FetchData(key, ttl)
	
	fmt.Println("\n[Step 2] Read immediately again (Data is fresh).")
	s.FetchData(key, ttl)
	
	fmt.Println("\n[Step 3] Wait 1.1 seconds for TTL to expire...")
	time.Sleep(1100 * time.Millisecond)
	
	fmt.Println("\n[Step 4] Read again it after TTL. Cache should purge it and fetch anew.")
	s.FetchData(key, ttl) // Will miss!
}
