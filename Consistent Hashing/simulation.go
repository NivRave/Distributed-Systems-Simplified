package main

import (
	"fmt"
	"sort"
)

// Simulation runs consistent hashing scenarios.
type Simulation struct {
	Ring      *ConsistencyRing
	TotalKeys int
}

// NewSimulation creates the simulation runner.
func NewSimulation() *Simulation {
	return &Simulation{
		Ring:      NewConsistencyRing(50), // 50 virtual nodes for better distribution
		TotalKeys: 2000,
	}
}

// Run executes the simulation lifecycle.
func (s *Simulation) Run() {
	fmt.Println("=== Consistent Hashing Simulation ===")
	fmt.Printf("Total Keys: %d\n", s.TotalKeys)

	// 1. Initial Cluster Setup
	nodes := []string{"Node-A", "Node-B", "Node-C"}
	for _, n := range nodes {
		s.Ring.AddNode(n)
	}
	fmt.Printf("\n[Step 1] Initial Cluster: %v\n", nodes)

	// 2. Distribute Keys
	// We map keys to their assigned nodes to track movement later.
	keyMap := make(map[string]string) // Key -> NodeID
	keys := make([]string, s.TotalKeys)

	initialDist := make(map[string]int)

	for i := 0; i < s.TotalKeys; i++ {
		key := fmt.Sprintf("user-session-%d", i)
		keys[i] = key

		node := s.Ring.GetNode(key)
		keyMap[key] = node
		initialDist[node]++
	}

	s.printDistribution(initialDist)

	// 3. Scale Up: Add a new node
	newNode := "Node-D"
	fmt.Printf("\n[Step 2] Adding %s to the cluster...\n", newNode)
	s.Ring.AddNode(newNode)

	// 4. Analyze Rebalancing
	// Consistent Hashing promise: Only k/N keys should move.
	// We expect ~1/4 of keys (500) to move to Node-D.
	movedCount := 0
	newDist := make(map[string]int)

	for _, key := range keys {
		prevNode := keyMap[key]
		currNode := s.Ring.GetNode(key)
		newDist[currNode]++

		if prevNode != currNode {
			movedCount++
			if movedCount <= 5 {
				// Show a few examples
				fmt.Printf("  -> Key %s moved from %s to %s\n", key, prevNode, currNode)
			}
		}
	}
	fmt.Printf("  ... (total %d keys moved)\n", movedCount)

	s.printDistribution(newDist)

	// 5. Statistics
	percentMoved := float64(movedCount) / float64(s.TotalKeys) * 100
	fmt.Printf("\n[Analysis]\n")
	fmt.Printf("Total Moved: %d keys (%.2f%%)\n", movedCount, percentMoved)
	fmt.Printf("Ideal Movement: 25.00%% (1/4 of keys)\n")
	fmt.Printf("Result: Very efficient rebalancing! (Modulo hashing would move ~75%%)\n")
}

func (s *Simulation) printDistribution(dist map[string]int) {
	fmt.Println("\n--- Node Distribution ---")

	// Sort keys for consistent output
	var keys []string
	for k := range dist {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, node := range keys {
		count := dist[node]
		percent := float64(count) / float64(s.TotalKeys) * 100

		// Simple ASCII bar
		barLen := int(percent / 2)
		bar := ""
		for i := 0; i < barLen; i++ {
			bar += "█"
		}

		fmt.Printf("%-8s: [%-25s] %4d keys (%.1f%%)\n", node, bar, count, percent)
	}
}
