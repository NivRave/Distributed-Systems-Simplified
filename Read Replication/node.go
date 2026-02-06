package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

// ReplicationEvent represents a data change to be replicated.
type ReplicationEvent struct {
	Key       string
	Value     string
	Timestamp int64
}

// Node represents a distributed system node (Leader or Follower).
type Node struct {
	ID       int
	IsLeader bool
	Data     map[string]string
	mu       sync.RWMutex
	replicas []chan ReplicationEvent // Send-only channels for replication
	logger   *slog.Logger
}

// NewNode creates a new Node instance.
func NewNode(id int, isLeader bool) *Node {
	return &Node{
		ID:       id,
		IsLeader: isLeader,
		Data:     make(map[string]string),
		replicas: make([]chan ReplicationEvent, 0),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})).With("node_id", id, "role", roleName(isLeader)),
	}
}

func roleName(isLeader bool) string {
	if isLeader {
		return "Leader"
	}
	return "Follower"
}

// AddReplica registers a follower's channel to receiving updates.
func (n *Node) AddReplica(ch chan ReplicationEvent) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.replicas = append(n.replicas, ch)
}

// Write (Leader only) updates local state and broadcasts to followers.
func (n *Node) Write(key, value string) error {
	if !n.IsLeader {
		return fmt.Errorf("node %d is not a leader", n.ID)
	}

	n.mu.Lock()
	n.Data[key] = value
	n.mu.Unlock()

	event := ReplicationEvent{
		Key:       key,
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	}

	n.logger.Info("Persisted data", "key", key, "value", value)

	// Broadcast to replicas asynchronously to avoid unnecessary blocking
	// In a real system, you might wait for a quorum here.
	for _, ch := range n.replicas {
		go func(c chan ReplicationEvent) {
			select {
			case c <- event:
			case <-time.After(500 * time.Millisecond):
				n.logger.Warn("Replication timed out for a follower")
			}
		}(ch)
	}
	return nil
}

// Read returns the value for a given key.
func (n *Node) Read(key string) (string, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	val, ok := n.Data[key]
	return val, ok
}

// StartFollower starts the listening loop for a follower node.
// It blocks until the context is cancelled.
func (n *Node) StartFollower(ctx context.Context, updates <-chan ReplicationEvent) {
	if n.IsLeader {
		n.logger.Error("Attempted to start follower loop on a leader node")
		return
	}

	n.logger.Info("Starting replication listener")

	for {
		select {
		case <-ctx.Done():
			n.logger.Info("Stopping replication listener")
			return
		case event := <-updates:
			// Simulate network latency
			time.Sleep(time.Duration(50+event.Timestamp%100) * time.Millisecond)

			n.mu.Lock()
			n.Data[event.Key] = event.Value
			n.mu.Unlock()

			n.logger.Info("Replicated data", "key", event.Key, "value", event.Value, "lag_ms", (time.Now().UnixNano()-event.Timestamp)/1000000)
		}
	}
}
