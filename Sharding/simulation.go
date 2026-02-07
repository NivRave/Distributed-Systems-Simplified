package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Simulation orchestrates the Sharding system.
type Simulation struct {
	Router    *RouterComponent
	NumShards int
	// Config
	TotalKeys  int
	KeyPrefix  string
	WriteDelay time.Duration
}

// NewSimulation creates a simulation with default settings.
func NewSimulation(numShards int) *Simulation {
	return &Simulation{
		Router:     NewRouterComponent(numShards),
		NumShards:  numShards,
		TotalKeys:  50,
		KeyPrefix:  "user",
		WriteDelay: 50 * time.Millisecond,
	}
}

// Run executes the full simulation lifecycle.
func (s *Simulation) Run(ctx context.Context) {
	fmt.Println("=== Initializing Sharding Topology ===")
	fmt.Printf("Topology ready: %d Shards\n", s.NumShards)

	// 1. Start Readers (Background)
	// They will try to read keys even before they exist, showing "MISS"
	go s.runReaders(ctx)

	// 2. Start Writers (Foreground Workload)
	fmt.Println("=== Starting Workload ===")
	s.runWriters(ctx)

	fmt.Println("=== Workload Complete. Showing final stats... ===")
	time.Sleep(1 * time.Second)

	// 3. Final Distribution
	s.Router.PrintDistribution()
}

func (s *Simulation) runWriters(ctx context.Context) {
	for i := 0; i < s.TotalKeys; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// Deterministic keys for demo purposes so readers can find them
			key := fmt.Sprintf("%s:%d", s.KeyPrefix, i)
			val := fmt.Sprintf("data-%d", i)

			fmt.Printf("\033[32m[Write]\033[0m Key: %s -> ", key)

			// Get routing info for visualization
			_, shardID := s.Router.GetShard(key)

			// Perform Write
			s.Router.AddKey(key, val)

			fmt.Printf("Shard %d\n", shardID)
			time.Sleep(s.WriteDelay)
		}
	}
}

func (s *Simulation) runReaders(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Pick a random key from the range [0, TotalKeys]
			// Some might be written already, some not.
			targetID := rand.Intn(s.TotalKeys)
			key := fmt.Sprintf("%s:%d", s.KeyPrefix, targetID)

			val, found := s.Router.GetKey(key)
			_, shardID := s.Router.GetShard(key) // Just to know where it came from

			if found {
				fmt.Printf("       \033[36m[Read Hit]\033[0m  Key: %s (from Shard %d) => %s\n", key, shardID, val)
			} else {
				// Don't spam misses too much, or maybe do to show emptiness
				// fmt.Printf("       \033[33m[Read Miss]\033[0m Key: %s (Shard %d empty)\n", key, shardID)
			}

			time.Sleep(s.randDuration(20, 100))
		}
	}
}

// Helper for random duration
func (s *Simulation) randDuration(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min)+min) * time.Millisecond
}
