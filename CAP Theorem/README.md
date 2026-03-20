# CAP Theorem: Consistency vs Availability

This project simulates the fundamental dilemma in distributed systems: **When the network breaks, do you stop the world (Consistency) or keep guessing (Availability)?**

## What is the CAP Theorem?

CAP states that a distributed data store can only provide **two** of the following three guarantees:

1.  **Consistency (C)**: Every read receives the most recent write or an error. (All nodes see the same data at the same time).
2.  **Availability (A)**: Every request receives a (non-error) response, without the guarantee that it contains the most recent write.
3.  **Partition Tolerance (P)**: The system continues to operate despite an arbitrary number of messages being dropped or delayed by the network.

### The "Gotcha": P is Mandatory
In a distributed system, you cannot "choose" to have no network failures. Cables get cut, routers crash. **Partition Tolerance (P) is a reality, not an option.**

Therefore, you must choose between:
*   **CP (Consistency + Partition Tolerance)**: If the network breaks, refuse updates to prevent data corruption. (e.g., Banking, Stock Markets).
*   **AP (Availability + Partition Tolerance)**: If the network breaks, accept updates even if they might conflict later. (e.g., Social Media feeds, Shopping Carts).

---

## How This Project Simulates the Theory

This code simulates a 2-node cluster and forces a **Network Partition** to demonstrate the divergence in behavior.

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **Cluster** | Two `Node` structs (Node-A, Node-B) that replicate data to each other. |
| **Network Partition (P)** | `network.Disconnect(NodeA, NodeB)` sets a boolean flag preventing message delivery. |
| **AP Strategy** | `Node.Write()` detects failure but logs a `[Warning]` and writes to local memory anyway. |
| **CP Strategy** | `Node.Write()` detects failure and returns `error("503 Service Unavailable")`, rejecting the write. |

### Key Files
*   `network.go`: Acts as the physical layer (cables). It has a switch to `Disconnect` nodes.
*   `node.go`: The Server logic. It has a `Mode` (AP or CP) switch that dictates its decision logic during a partition.
*   `simulation.go`: The orchestration engine. It runs two scenarios back-to-back:
    1.  **The Social Network (AP)**: Keeps accepting posts even if followers can't see them yet.
    2.  **The Bank (CP)**: Refuses transactions if it can't guarantee safety ("Split Brain" prevention).

---

## How to Run

Ensure you have Go installed.

```bash
cd "CAP Theorem"
go run .
```

### Sample Output Explanation

**1. Scenario: AP (Social Network)**
Goal: Keep the service UP, even if data is stale.

```text
[Network] [Partition] Connection severed between Node-A and Node-B...

[Step 3] Client writes to Node-A...
   -> [Warning] [AP Decision] Replication Failed, but proceeding (Availability First!).
   -> [Success] [AP Decision] Replication OK.

[Step 4] Reading from both nodes
   Node-A (Local):  'Status: Busy'
   Node-B (Remote): 'Status: Online'

Result: System is AVAILABLE (Write succeeded), but INCONSISTENT.
```
> **Observation**: Node-A accepted the new status "Busy". However, because the cable is cut, Node-B still thinks the status is "Online". The system is **Available** (no errors), but **Inconsistent** (users see different things).

**2. Scenario: CP (Bank)**
Goal: Prevent conflicting data, even if we must go down.

```text
[Network] [Partition] Connection severed between Node-A and Node-B...

[Step 2] Client attempts transaction on Node-A...
   -> [Error] [CP Decision] Replication Failed. Aborting Write to serve Consistency.
   Client received error: 503 Service Unavailable: Replication Quorum Failed

Result: System is CONSISTENT (Data matches), but UNAVAILABLE for writes.
```
> **Observation**: Node-A realized it couldn't talk to Node-B. To prevent a "Split Brain" (where both sides accept money transfers unknowingly), it **shut down** the write operation. The system is **Consistent** (no corrupted data), but **Unavailable**.
