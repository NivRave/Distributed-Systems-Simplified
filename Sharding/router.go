package main

import (
	"fmt"
	"hash/crc32"
	"log/slog"
	"os"
)

// RouterComponent is responsible for routing client requests to the correct shard.
// It implements the Sharding Strategy (Hash-based partitioning).
type RouterComponent struct {
	shards map[int]*ShardNode
	count  int
	logger *slog.Logger
}

// NewRouterComponent creates a new router with N shards.
func NewRouterComponent(numShards int) *RouterComponent {
	shards := make(map[int]*ShardNode)
	for i := 0; i < numShards; i++ {
		shards[i] = NewShardNode(i)
	}

	return &RouterComponent{
		shards: shards,
		count:  numShards,
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})).With("component", "ROUTER"),
	}
}

// GetShard calculates the target shard ID for a given key.
// Strategy: Hash(key) % NumberOfShards
func (r *RouterComponent) GetShard(key string) (*ShardNode, int) {
	checksum := crc32.ChecksumIEEE([]byte(key))
	shardID := int(checksum) % r.count
	r.logger.Info("Routing request", "key", key, "hash", checksum, "shard_id", shardID)
	return r.shards[shardID], shardID
}

// AddKey routes a write request to the appropriate shard.
func (r *RouterComponent) AddKey(key, value string) {
	shard, _ := r.GetShard(key)
	shard.Write(key, value)
}

// GetKey routes a read request to the appropriate shard.
func (r *RouterComponent) GetKey(key string) (string, bool) {
	shard, _ := r.GetShard(key)
	return shard.Read(key)
}

// PrintDistribution logs how many keys ended up on each shard.
// This is critical to visualize data skew vs. uniform distribution.
func (r *RouterComponent) PrintDistribution() {
	fmt.Println("\n=== Shard Distribution Stats ===")
	for i := 0; i < r.count; i++ {
		count := r.shards[i].Stats()
		bar := ""
		for j := 0; j < count; j++ {
			bar += "█"
		}
		fmt.Printf("Shard %d: [%-20s] %d keys\n", i, bar, count)
	}
	fmt.Println("================================")
}
