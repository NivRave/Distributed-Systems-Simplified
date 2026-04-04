package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Server represents a single backend application instance inside a data center.
type Server struct {
	ID      string
	Address string
	IsAlive bool
	mu      sync.RWMutex
}

func NewServer(id, address string) *Server {
	return &Server{
		ID:      id,
		Address: address,
		IsAlive: true,
	}
}

// ServeRequest simulates accepting and processing client API requests.
func (s *Server) ServeRequest(req string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.IsAlive {
		return "", errors.New("connection refused (Server is down)")
	}

	// Tiny artificial delay for processing
	time.Sleep(2 * time.Millisecond) 
	return fmt.Sprintf("200 OK Response from %s (%s)", s.ID, s.Address), nil
}

// Ping simulates an internal '/healthz' HTTP endpoint called by the Load Balancer setup.
func (s *Server) Ping() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.IsAlive {
		return errors.New("timeout reached")
	}
	return nil
}

// Crash artificially simulates an unhandled panic or hardware failure.
func (s *Server) Crash() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAlive = false
	fmt.Printf("\n[SIMULATION] 💥 Hardware Failure -> %s has CRASHED!\n", s.ID)
}

// Recover simulates rebooting the instance successfully.
func (s *Server) Recover() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAlive = true
	fmt.Printf("\n[SIMULATION] 🛠️ Systems Admin intervened -> %s is BACK ONLINE!\n", s.ID)
}
