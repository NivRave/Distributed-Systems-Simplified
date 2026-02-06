package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Simulation orchestrates the Leader-Follower system and generates load.
type Simulation struct {
	Leader    *Node
	Followers []*Node
	// Configuration
	NumFollowers  int
	WorkloadCount int
	MaxDelay      int // Max simulated network delay in ms
}

// NewSimulation creates a simulation with default settings.
func NewSimulation(numFollowers int) *Simulation {
	return &Simulation{
		NumFollowers:  numFollowers,
		WorkloadCount: 5,
		MaxDelay:      800,
	}
}

// Run executes the full simulation lifecycle.
func (s *Simulation) Run(ctx context.Context) {
	fmt.Println("=== Initializing Simulation Topology ===")

	// 1. Initialize Leader
	s.Leader = NewNode(0, true)

	// 2. Initialize Followers
	s.Followers = make([]*Node, s.NumFollowers)
	for i := 0; i < s.NumFollowers; i++ {
		s.Followers[i] = NewNode(i+1, false)

		// Create and register replication channel
		replChan := make(chan ReplicationEvent, 10)
		s.Leader.AddReplica(replChan)

		// Start follower background listener
		go s.Followers[i].StartFollower(ctx, replChan)
	}

	fmt.Printf("Topology ready: 1 Leader + %d Followers\n", s.NumFollowers)

	// 3. Start Background Workload & Readers
	// We use a WaitGroup (implicitly or explicitly) or just channels to coordinate.
	// Since the original code relied on time and context, we'll keep it simple but structured.

	// Start Readers
	go s.runReaders(ctx)

	// Run Writers (Blocking or separate)
	// We'll run writers in a goroutine so we can wait for them specifically if needed,
	// or just run them here. To match previous behavior + gracefulness:

	fmt.Println("=== Starting Workload and Readers ===")
	s.runWriters(ctx)

	fmt.Println("=== Workload Complete. Waiting for final replication... ===")
	time.Sleep(1 * time.Second)
}

func (s *Simulation) runWriters(ctx context.Context) {
	keys := []string{"user:100", "product:550", "order:999"}

	for i := 0; i < s.WorkloadCount; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// Generate random write
			key := keys[rand.Intn(len(keys))]
			val := fmt.Sprintf("v%d-%d", i, rand.Intn(1000))

			// Write to Leader
			err := s.Leader.Write(key, val)
			if err != nil {
				fmt.Printf("Error writing to leader: %v\n", err)
			}

			// Simulate user think time / delay
			time.Sleep(time.Duration(rand.Intn(s.MaxDelay)) * time.Millisecond)
		}
	}
}

func (s *Simulation) runReaders(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(s.Followers) == 0 {
				return
			}

			// Randomly pick a follower
			f := s.Followers[rand.Intn(len(s.Followers))]

			// Random key (focusing on one mostly for consistency checks)
			key := "user:100"

			val, ok := f.Read(key)
			if ok {
				// We use a distinct format to differentiate from system logs
				fmt.Fprintf(os.Stdout, "\033[36m[Client Read]\033[0m Node %d served key=%s value=%s\n", f.ID, key, val)
			} else {
				fmt.Fprintf(os.Stdout, "\033[33m[Client Read]\033[0m Node %d miss for key=%s (not propagated yet)\n", f.ID, key)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
