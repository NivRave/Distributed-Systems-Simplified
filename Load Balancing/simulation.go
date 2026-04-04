package main

import (
	"fmt"
)

type Simulation struct {
	LB            *LoadBalancer
	HealthChecker *HealthChecker
	Servers       []*Server
}

func NewSimulation() *Simulation {
	// 1. Initialize Backend Servers
	srv1 := NewServer("Server-A", "10.0.0.1")
	srv2 := NewServer("Server-B", "10.0.0.2")
	srv3 := NewServer("Server-C", "10.0.0.3")

	servers := []*Server{srv1, srv2, srv3}
	
	// 2. Wrap them behind a Load Balancer
	lb := NewLoadBalancer(servers)
	
	// 3. Attach a Health Checker daemon
	hc := NewHealthChecker(servers, lb)

	return &Simulation{
		LB:            lb,
		HealthChecker: hc,
		Servers:       servers,
	}
}

// SendClientTraffic simulates external users making HTTP requests.
func (sim *Simulation) SendClientTraffic(count int) {
	for i := 1; i <= count; i++ {
		req := fmt.Sprintf("Client_Req_%d", i)
		resp, err := sim.LB.RouteRequest(req)
		
		if err != nil {
			fmt.Printf("   <- [Client] Error: %v\n", err)
		} else {
			fmt.Printf("   <- [Client] RX: %s\n", resp)
		}
	}
}

func (sim *Simulation) Run() {
	fmt.Println("=== Load Balancing & Automatic Failover Simulation ===")
	
	// Initialize the active pool before serving traffic
	sim.HealthChecker.RunCheck()

	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 1: Happy Path (Round-Robin Routing)")
	fmt.Println("Goal: Show perfectly even traffic distribution across 3 healthy nodes.")
	fmt.Println("-----------------------------------------------------")
	
	fmt.Println("\n[Actions] Firing 6 concurrent requests through Load Balancer:")
	sim.SendClientTraffic(6)

	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 2: Node Failure (Automatic Failover)")
	fmt.Println("Goal: Show the system routing around a crashed server safely.")
	fmt.Println("-----------------------------------------------------")
	
	// Artificially crash Server B
	sim.Servers[1].Crash() 
	
	fmt.Println("\n[Background] Health Checker heavily sweeps the cluster...")
	sim.HealthChecker.RunCheck() // Detects failure and surgically removes B
	
	fmt.Println("\n[Actions] Firing 4 new client requests:")
	sim.SendClientTraffic(4)

	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Scenario 3: Node Recovery")
	fmt.Println("Goal: Reintegrating a server safely once it passes medical health checks again.")
	fmt.Println("-----------------------------------------------------")
	
	// Reboot Server B
	sim.Servers[1].Recover()
	
	fmt.Println("\n[Background] Health Checker sweeps the cluster...")
	sim.HealthChecker.RunCheck() // Detects recovery and re-adds B
	
	fmt.Println("\n[Actions] Firing 3 final requests to confirm Normal Operations:")
	sim.SendClientTraffic(3)
}
