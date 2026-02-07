package main

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

// ShardNode represents a single shard in the cluster.
// It is responsible for storing a subset of the total data.
type ShardNode struct {
	ID     int
	Data   map[string]string
	mu     sync.RWMutex
	logger *slog.Logger
}

func NewShardNode(id int) *ShardNode {
	return &ShardNode{
		ID:   id,
		Data: make(map[string]string),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})).With("shard_id", id),
	}
}

// Write simulates a database write operation.
// It includes a small artificial delay to mimic disk I/O.
func (sn *ShardNode) Write(key, value string) {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	// Simulate disk I/O latency
	time.Sleep(10 * time.Millisecond)

	sn.Data[key] = value
	sn.logger.Info("Data written to shard", "key", key, "value", value)
}

// Read simulates a database read operation.
func (sn *ShardNode) Read(key string) (string, bool) {
	sn.mu.RLock()
	defer sn.mu.RUnlock()

	val, ok := sn.Data[key]
	if ok {
		sn.logger.Info("Data read from shard", "key", key)
	} else {
		sn.logger.Warn("Key not found in shard", "key", key)
	}
	return val, ok
}

// Stats returns the number of keys stored in this shard.
func (sn *ShardNode) Stats() int {
	sn.mu.RLock()
	defer sn.mu.RUnlock()
	return len(sn.Data)
}
