package main

import (
	"context"
	"testing"
	"time"
)

func TestNodeWriteAndRead(t *testing.T) {
	// Create a leader
	leader := NewNode(1, true)

	key := "test_key"
	value := "test_val"

	err := leader.Write(key, value)
	if err != nil {
		t.Fatalf("Failed to write to leader: %v", err)
	}

	readVal, ok := leader.Read(key)
	if !ok {
		t.Fatalf("Leader should have data immediately")
	}
	if readVal != value {
		t.Errorf("Expected %s, got %s", value, readVal)
	}
}

func TestReplication(t *testing.T) {
	leader := NewNode(1, true)
	follower := NewNode(2, false)

	// Create a channel to connect them
	replChan := make(chan ReplicationEvent, 10)
	leader.AddReplica(replChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start follower
	go follower.StartFollower(ctx, replChan)

	// Write to leader
	key := "repl_key"
	value := "repl_val"

	err := leader.Write(key, value)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Wait for eventually consistency (allow some time for channel and sleep simulation)
	time.Sleep(200 * time.Millisecond)

	readVal, ok := follower.Read(key)
	if !ok {
		// Wait a bit longer if needed (flaky test prevention)
		time.Sleep(500 * time.Millisecond)
		readVal, ok = follower.Read(key)
		if !ok {
			t.Fatalf("Follower did not receive data")
		}
	}

	if readVal != value {
		t.Errorf("Follower has wrong value. Expected %s, got %s", value, readVal)
	}
}

func TestWriteOnFollowerShouldFail(t *testing.T) {
	follower := NewNode(2, false)
	err := follower.Write("k", "v")
	if err == nil {
		t.Error("Writing to a follower should return an error")
	}
}
