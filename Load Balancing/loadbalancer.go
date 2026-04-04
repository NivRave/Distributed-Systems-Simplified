package main

import (
	"fmt"
	"sync"
)

// LoadBalancer manages routing traffic optimally across multiple scalable servers.
type LoadBalancer struct {
	activeServers []*Server // Only holds currently healthy components
	currentIndex  int       // Counter for round-robin math
	mu            sync.RWMutex
}

func NewLoadBalancer(initialServers []*Server) *LoadBalancer {
	return &LoadBalancer{
		activeServers: initialServers,
		currentIndex:  0,
	}
}

// UpdateActivePool is triggered by the Health Checker daemon to adjust available capacity natively.
func (lb *LoadBalancer) UpdateActivePool(healthyList []*Server) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	lb.activeServers = healthyList
	// Reset the round-robin index safely so we don't 'out-of-bounds' short circuit
	lb.currentIndex = 0 
}

// RouteRequest implements the logic to forward traffic using the Round-Robin algorithm.
func (lb *LoadBalancer) RouteRequest(req string) (string, error) {
	lb.mu.Lock()
	poolSize := len(lb.activeServers)
	
	if poolSize == 0 {
		lb.mu.Unlock() // avoid hanging lock on return early
		return "", fmt.Errorf("[503 Service Unavailable] Zero healthy backend servers to field the request")
	}

	// Mathematical Round Robin calculation (e.g. 0 % 3 = Server 0, 1 % 3 = Server 1)
	selectedIndex := lb.currentIndex % poolSize
	targetServer := lb.activeServers[selectedIndex]
	
	lb.currentIndex++ // Iterate for the next request in line
	lb.mu.Unlock()

	// Network jump
	fmt.Printf("   [LB Router]   -> Forwarding '%s' strictly to %s\n", req, targetServer.ID)
	
	return targetServer.ServeRequest(req)
}
