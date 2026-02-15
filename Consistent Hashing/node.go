package main

import (
	"fmt"
	"sync"
)

// Node represents a storage server in the cluster.
type Node struct {
	ID   string
	Data map[string]string
	mu   sync.RWMutex
}

func NewNode(id string) *Node {
	return &Node{
		ID:   id,
		Data: make(map[string]string),
	}
}

func (n *Node) Put(key, value string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Data[key] = value
}

func (n *Node) Get(key string) (string, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	val, ok := n.Data[key]
	return val, ok
}

// KeyCount returns number of items stored
func (n *Node) KeyCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.Data)
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %s: %d keys", n.ID, n.KeyCount())
}
