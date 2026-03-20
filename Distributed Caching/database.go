package main

import (
	"fmt"
	"time"
)

// Database simulates a slow, persistent storage backend.
type Database struct {
	data map[string]string
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]string),
	}
}

// Read simulates a slow query by intentionally blocking for 50ms.
func (db *Database) Read(key string) string {
	time.Sleep(50 * time.Millisecond) // Artificial latency
	
	if val, ok := db.data[key]; ok {
		return val
	}
	
	// Default value if not found
	return fmt.Sprintf("DataFor(%s)", key)
}

// Write simulates a regular write to persistent storage.
func (db *Database) Write(key, value string) {
	db.data[key] = value
}
