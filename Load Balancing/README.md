# Load Balancing & Health Checks

This project simulates how a **Load Balancer** elegantly distributes traffic across multiple backend servers to prevent overload, and how it actively monitors their health to provide **Automatic Failover** when hardware crashes.

## What is Load Balancing?

Imagine a popular website that gets millions of requests per second. A single server would inevitably crash under the load. Instead, we use multiple backend servers hidden behind a single entry point: the Load Balancer.

### Core Concepts Simulated Here

1.  **Round-Robin Routing**: 
    The easiest way to distribute traffic fairly. The Load Balancer simply goes sequentially down the list: Request 1 goes to Server A, Request 2 to Server B, Request 3 to Server C, Request 4 wraps back to Server A.

2.  **Health Checks (Active Probing)**:
    If a server crashes, the Load Balancer needs to know immediately. Otherwise, it will send user traffic to a dead server (resulting in 503 errors). A background daemon constantly "pings" the servers (e.g., hitting a `/healthz` HTTP endpoint).

3.  **Automatic Failover**:
    When the Health Checker notices a server is dead, it surgically removes it from the Load Balancer's active pool. The Round-Robin algorithm logically routes around the dead node. When the node recovers eventually, it is integrated seamlessly back into the pool.

---

## How This Project Simulates the Theory

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **Backend Server** | `Server` struct. Contains a `Ping()` health endpoint and a mock `ServeRequest()` request handler. |
| **Hardware Crash** | The `Crash()` function artificially flips a boolean `IsAlive = false`, simulating an unrecoverable panic. |
| **Round Robin Algorithm** | `loadbalancer.go` utilizes modulo arithmetic logic to loop indices: `idx = lb.currentIndex % len(activeServers)`. |
| **Active Monitoring** | A dedicated watcher (`HealthChecker`) loops through all known servers. If one fails, it actively mutates the LB's routing array. |

### Key Files
*   `server.go`: The backend compute node holding mock API business logic.
*   `loadbalancer.go`: The traffic router. Exposes only currently healthy servers to external clients.
*   `healthcheck.go`: The daemon watcher. Detects anomalies and updates the Load Balancer dynamically.
*   `simulation.go`: The orchestration script demonstrating the failover event timeline accurately.

---

## How to Run

Ensure you have Go installed natively.

```bash
cd "Load Balancing"
go run .
```

### Sample Output Explanation

**1. Scenario: Happy Path**
```text
[Actions] Firing 6 concurrent requests through Load Balancer:
   [LB Router]   -> Forwarding 'Client_Req_1' strictly to Server-A
   <- [Client] RX: 200 OK Response from Server-A (10.0.0.1)
   [LB Router]   -> Forwarding 'Client_Req_2' strictly to Server-B
   ...
   [LB Router]   -> Forwarding 'Client_Req_4' strictly to Server-A
```
> **Observation**: The requests alternate perfectly via `A -> B -> C -> A -> B -> C`. The systemic load is mathematically balanced globally.

**2. Scenario: Automatic Failover**
```text
[SIMULATION] 💥 Hardware Failure -> Server-B has CRASHED!
[Background] Health Checker heavily sweeps the cluster...
   [HealthCheck] Probe Failed! Node Server-B is DEAD. Evicting from LB pool.

[Actions] Firing 4 new client requests:
   [LB Router]   -> Forwarding 'Client_Req_1' strictly to Server-A
   [LB Router]   -> Forwarding 'Client_Req_2' strictly to Server-C
   [LB Router]   -> Forwarding 'Client_Req_3' strictly to Server-A
```
> **Observation**: Server-B is completely dead and unreachable. Because the active Health Checker removed it, the Load Balancer organically begins routing exclusively `A -> C -> A -> C`. Zero physical client requests were dropped by the actual web router.

**3. Scenario: Node Recovery**
```text
[SIMULATION] 🛠️ Systems Admin intervened -> Server-B is BACK ONLINE!
[Background] Health Checker sweeps the cluster...
   [HealthCheck] Probe Restored! Node Server-B is HEALTHY. Re-attaching to LB pool.
```
> **Observation**: The application instance recovers. The background probe discovers it autonomously and pushes it live instantly, allowing it to receive fractions of the overall traffic efficiently.
