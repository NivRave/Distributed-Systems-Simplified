package main

import (
	"fmt"
)

// HealthChecker is responsible for identifying system failures cleanly without affecting client load.
type HealthChecker struct {
	allServers []*Server
	lb         *LoadBalancer
	knownDead  map[string]bool // Used to prevent spamming the console logs.
}

func NewHealthChecker(servers []*Server, lb *LoadBalancer) *HealthChecker {
	return &HealthChecker{
		allServers: servers,
		lb:         lb,
		knownDead:  make(map[string]bool),
	}
}

// RunCheck acts as a sweeping interval task. 
// Note: In real applications, this loops on a `time.Ticker`, but for cleaner visual simulation output,
// we invoke the sweep directly in a deterministic linear flow.
func (hc *HealthChecker) RunCheck() {
	var healthyPool []*Server
	poolChanged := false

	for _, srv := range hc.allServers {
		// Native ping / health probe
		err := srv.Ping()
		isDead := (err != nil)

		// State Change Tracker: Just died
		if isDead && !hc.knownDead[srv.ID] {
			fmt.Printf("   \033[31m[HealthCheck]\033[0m Probe Failed! Node %s is DEAD. Evicting from LB pool.\n", srv.ID)
			hc.knownDead[srv.ID] = true
			poolChanged = true
		} else if !isDead && hc.knownDead[srv.ID] {
			// State Change Tracker: Just recovered
			fmt.Printf("   \033[32m[HealthCheck]\033[0m Probe Restored! Node %s is HEALTHY. Re-attaching to LB pool.\n", srv.ID)
			hc.knownDead[srv.ID] = false
			poolChanged = true
		}

		if !isDead {
			healthyPool = append(healthyPool, srv)
		}
	}

	// Publish changes to central router logic so the load balancer doesn't route dead traffic anymore.
	if poolChanged {
		hc.lb.UpdateActivePool(healthyPool)
		fmt.Printf("   [HealthCheck] LB Active Pool sync updated. Target Servers available: %d/%d\n", len(healthyPool), len(hc.allServers))
	} else if len(healthyPool) == len(hc.allServers) {
		// Log just for initial setup sanity
		fmt.Printf("   [HealthCheck] Sweep Complete. All %d servers are completely healthy.\n", len(hc.allServers))
	}
}
