package main

import (
	"fmt"
	"sync"
)

// ConsistencyMode determines how a node reacts to partitions.
type ConsistencyMode string

const (
	ModeCP ConsistencyMode = "CP (Strong Consistency)"
	ModeAP ConsistencyMode = "AP (High Availability)"
)

// Node represents a server in our distributed bank/social network.
type Node struct {
	Name    string
	Data    string          // Represents sensitive data (e.g., Balance or Status)
	Mode    ConsistencyMode // How this node behaves under stress
	Network *Network
	Peer    *Node // The single peer for this demo (simplification)
	mu      sync.RWMutex
}

func NewNode(name string, mode ConsistencyMode, net *Network) *Node {
	return &Node{
		Name:    name,
		Data:    "InitialState",
		Mode:    mode,
		Network: net, // Pointer to the shared network switch
	}
}

// Write attempts to update the data.
// This is the CRITICAL logic for CAP Theorem.
func (n *Node) Write(newValue string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("[%s] Received Write Request ('%s')...\n", n.Name, newValue)

	// Attempt Replication First (Synchronous for CP logic simulation)
	// In reality, AP is often async, but here we simulate the *result* of the network call.
	err := n.Network.SendMessage(n, n.Peer, newValue)

	if n.Mode == ModeCP {
		// CP Logic: If we can't talk to the peer, we MUST reject the write.
		// Otherwise we risk Split-Brain (two different truths).
		if err != nil {
			fmt.Printf("   -> [Error] [CP Decision] Replication Failed. Aborting Write to serve Consistency.\n")
			return fmt.Errorf("503 Service Unavailable: Replication Quorum Failed")
		}
		fmt.Printf("   -> [Success] [CP Decision] Replication OK. Committing.\n")
	} else {
		// AP Logic: We don't care if the peer is dead. We accept the write anyway.
		if err != nil {
			fmt.Printf("   -> [Warning] [AP Decision] Replication Failed, but proceeding (Availability First!).\n")
		} else {
			fmt.Printf("   -> [Success] [AP Decision] Replication OK.\n")
		}
	}

	// Commit locally (For CP, we only reach here if err == nil)
	n.Data = newValue
	return nil
}

// Receive handles incoming replication messages from the network.
func (n *Node) Receive(msg string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Data = msg
	// fmt.Printf("   -> [%s] Updated via replication.\n", n.Name)
	return nil
}

// Read returns the local state.
func (n *Node) Read() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Data
}
