package main

import (
	"fmt"
	"time"
)

type Simulation struct{}

func (s *Simulation) Run() {
	fmt.Println("=== CAP Theorem Simulation ===")

	s.RunScenarioAP()
	time.Sleep(1 * time.Second) // Pause for effect
	s.RunScenarioCP()
}

func (s *Simulation) RunScenarioAP() {
	fmt.Println("\n--------------------------------------------------")
	fmt.Println("Scenario 1: Availability (AP) - The Social Network")
	fmt.Println("Goal: Keep the service UP, even if data is stale.")
	fmt.Println("--------------------------------------------------")

	// Setup
	net := NewNetwork()
	nodeA := NewNode("Node-A", ModeAP, net)
	nodeB := NewNode("Node-B", ModeAP, net)

	// Link up (Peer Discovery)
	nodeA.Peer = nodeB
	nodeB.Peer = nodeA
	net.Connect(nodeA, nodeB)

	// Normal Operation
	fmt.Println("\n[Step 1] Normal Operation (Network Healthy)")
	nodeA.Write("Status: Online")
	fmt.Printf("   Node-B sees: '%s' (Consistent)\n", nodeB.Read())

	// Partition
	fmt.Println("\n[Step 2] Network Partition Occurs (Cable Cut)")
	net.Disconnect(nodeA, nodeB)

	// Write during Partition
	fmt.Println("\n[Step 3] Client updates status on Node-A")
	nodeA.Write("Status: Busy")

	// Check Result
	fmt.Println("\n[Step 4] Reading from both nodes")
	valA := nodeA.Read()
	valB := nodeB.Read()
	fmt.Printf("   Node-A (Local):  '%s'\n", valA)
	fmt.Printf("   Node-B (Remote): '%s'\n", valB)

	if valA != valB {
		fmt.Println("\nResult: System is AVAILABLE (Write succeeded), but INCONSISTENT.")
	}
}

func (s *Simulation) RunScenarioCP() {
	fmt.Println("\n--------------------------------------------------")
	fmt.Println("Scenario 2: Consistency (CP) - The Bank")
	fmt.Println("Goal: Prevent conflicting data, even if we must go down.")
	fmt.Println("--------------------------------------------------")

	// Setup
	net := NewNetwork()
	nodeA := NewNode("Node-A", ModeCP, net)
	nodeB := NewNode("Node-B", ModeCP, net)

	// Link up
	nodeA.Peer = nodeB
	nodeB.Peer = nodeA
	net.Connect(nodeA, nodeB)

	// Partition
	fmt.Println("\n[Step 1] Network Partition Occurs")
	net.Disconnect(nodeA, nodeB)

	// Write during Partition
	fmt.Println("\n[Step 2] Client attempts transaction on Node-A")
	err := nodeA.Write("Balance: $500")

	if err != nil {
		fmt.Printf("   Client received error: %v\n", err)
	}

	// Check Result
	fmt.Println("\n[Step 3] Verifying Data Integrity")
	valA := nodeA.Read()
	valB := nodeB.Read()
	fmt.Printf("   Node-A State: '%s'\n", valA)
	fmt.Printf("   Node-B State: '%s'\n", valB)

	if valA == valB {
		fmt.Println("\nResult: System is CONSISTENT (Data matches), but UNAVAILABLE for writes.")
	}
}
