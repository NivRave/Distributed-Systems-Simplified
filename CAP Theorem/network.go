package main

import (
	"errors"
	"fmt"
	"sync"
)

// Network acts as a switch between nodes (simulating cables).
type Network struct {
	connections map[string]bool // "NodeA->NodeB" = true/false
	mu          sync.RWMutex
}

func NewNetwork() *Network {
	return &Network{connections: make(map[string]bool)}
}

// Connect establishes a bidirectional link.
func (n *Network) Connect(nodeA, nodeB *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	key1 := fmt.Sprintf("%s->%s", nodeA.Name, nodeB.Name)
	key2 := fmt.Sprintf("%s->%s", nodeB.Name, nodeA.Name)
	n.connections[key1] = true
	n.connections[key2] = true
}

// Disconnect simulates a network partition (cable cut).
func (n *Network) Disconnect(nodeA, nodeB *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	key1 := fmt.Sprintf("%s->%s", nodeA.Name, nodeB.Name)
	key2 := fmt.Sprintf("%s->%s", nodeB.Name, nodeA.Name)
	n.connections[key1] = false
	n.connections[key2] = false

	fmt.Printf("\n[Network] [Partition] Connection severed between %s and %s...\n", nodeA.Name, nodeB.Name)
}

// SendMessage simulates RPC. Returns error if partitioned.
func (n *Network) SendMessage(sender, receiver *Node, msg string) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	key := fmt.Sprintf("%s->%s", sender.Name, receiver.Name)
	if !n.connections[key] {
		return errors.New("network_error: destination unreachable")
	}

	// Direct function call simulates RPC
	return receiver.Receive(msg)
}
