# Read Replication Pattern

This project is a simulation of the **Read Replication** (or Leader-Follower) architectural pattern, commonly used in distributed databases like PostgreSQL (hot standby), MySQL, and MongoDB.

## What is Read Replication?

Read Replication is a database architecture that separates **Write** operations from **Read** operations to scale performance.

*   **Leader (Primary/Master)**: The single source of truth. Handles all write operations (INSERT, UPDATE, DELETE) and replicates changes to followers.
*   **Followers (Replicas/Slaves)**: Read-only copies of the database. They apply updates from the Leader asynchronously and serve read requests.

### Why use it?
1.  **Read Scalability**: Most web applications are "read-heavy" (e.g., 99 reads for every 1 write). You can spread this load across 10 or 100 followers.
2.  **High Availability**: If the Leader fails, one of the Followers can be promoted to become the new Leader.
3.  **Geographic Proximity**: You can place read replicas in different regions closer to users to reduce latency.

### The Trade-off: Eventual Consistency
Because replication is often asynchronous (to keep writes fast), a Follower might be a few milliseconds behind the Leader. A user might write data and immediately try to read it, but find the data missing or "stale". This is known as **Replication Lag**.

---

## How This Project Simulates the Theory

This code models a distributed system in a single process using Go's concurrency primitives.

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **Leader Node** | `Node` struct with `IsLeader=true`. Only this node accepts `Write()` calls. |
| **Follower Nodes** | `Node` structs with `IsLeader=false`. They listen on a channel for updates. |
| **Network Replication** | `replicas` channel (buffered). Sends `ReplicationEvent` structs. |
| **Replication Lag** | Artificial `time.Sleep` added in the follower loop to mimic network delays. |
| **Eventual Consistency** | The console logs show `[Client Read]` events. You will see followers serve data *after* the leader writes it, often with a delay. |

### Key Files
*   `node.go`: Defines the Node behavior (Leader writes, Follower listens).
*   `simulation.go`: Sets up the topology (1 Leader, 3 Followers) and generates artificial traffic.

---

## How to Run

Ensure you have Go installed.

```bash
go run .
```

### Sample Output Explanation

```text
time=... msg="Persisted data" role=Leader key=user:100 value=v0-123 ...
```
> **Write**: The Leader saves data to its local map and broadcasts the event.

```text
[Client Read] Node 1 miss for key=user:100 (not propagated yet)
```
> **Stale Read**: A client tried to read from Node 1 immediately, but the data hadn't arrived yet (Consistency Trade-off).

```text
time=... msg="Replicated data" node_id=1 ... lag_ms=50
[Client Read] Node 1 served key=user:100 value=v0-123
```
> **Eventual Consistency**: After 50ms, the data arrives at Node 1, and subsequent reads are successful.
