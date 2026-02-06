package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. Setup Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down simulation...")
		cancel()
	}()

	// 2. Run Simulation
	sim := NewSimulation(3) // 3 Followers
	sim.Run(ctx)

	fmt.Println("Simulation finished.")
}
