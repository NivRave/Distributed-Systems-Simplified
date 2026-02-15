package main

import (
	"fmt"
	"hash/crc32"
	"sort"
)

// ConsistencyRing manages the consistent hashing ring.
// It maps keys to nodes using virtual nodes for better distribution.
type ConsistencyRing struct {
	virtualNodes int               // Number of virtual replicas per physical node
	sortedHashes []uint32          // Sorted list of all virtual node hashes
	ring         map[uint32]string // Map from virtual node hash -> physical node ID
}

// NewConsistencyRing creates a new ring.
// virtualNodes: recommended 10-100 for better distribution.
func NewConsistencyRing(virtualNodes int) *ConsistencyRing {
	return &ConsistencyRing{
		virtualNodes: virtualNodes,
		sortedHashes: []uint32{},
		ring:         make(map[uint32]string),
	}
}

// AddNode adds a physical node to the ring (creating virtual replicas).
func (c *ConsistencyRing) AddNode(nodeID string) {
	for i := 0; i < c.virtualNodes; i++ {
		// Create virtual node key: e.g. "NodeA#1"
		virtualKey := fmt.Sprintf("%s#%d", nodeID, i)
		hash := c.hash(virtualKey)

		c.sortedHashes = append(c.sortedHashes, hash)
		c.ring[hash] = nodeID
	}
	// Keep the ring sorted for binary search
	sort.Slice(c.sortedHashes, func(i, j int) bool {
		return c.sortedHashes[i] < c.sortedHashes[j]
	})
}

// RemoveNode removes a physical node from the ring.
func (c *ConsistencyRing) RemoveNode(nodeID string) {
	newSortedHashes := []uint32{}
	// Rebuild slice without the removed node's hashes
	for _, h := range c.sortedHashes {
		if c.ring[h] == nodeID {
			delete(c.ring, h)
		} else {
			newSortedHashes = append(newSortedHashes, h)
		}
	}
	c.sortedHashes = newSortedHashes
}

// GetNode returns the physical node ID responsible for a given key.
func (c *ConsistencyRing) GetNode(key string) string {
	if len(c.sortedHashes) == 0 {
		return ""
	}

	hash := c.hash(key)

	// Binary search for the first virtual node where hash(node) >= hash(key)
	idx := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] >= hash
	})

	// Wrap around to the start of the ring if we went past the end
	if idx == len(c.sortedHashes) {
		idx = 0
	}

	return c.ring[c.sortedHashes[idx]]
}

func (c *ConsistencyRing) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
