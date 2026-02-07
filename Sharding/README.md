# Sharding (Horizontal Partitioning) Simulation

This project demonstrates the concept of **Database Sharding** (specifically Hash-Based Sharding) using Go.

## What is Sharding?

Sharding is a method of splitting and storing a single logical dataset across multiple databases (physical shards). This allows systems to scale horizontally beyond the limits of a single server.

*   **Logic**: Instead of a monolithic DB, data is partitioned.
*   **Routing**: A `Router` or Proxy sits in front of the shards and decides: "Key 'user:123' belongs to Shard 2".

## Sharding Strategies

Partitioning data effectively requires understanding your access patterns. Here are three common strategies:

### 1. Range-Based Sharding
Divides data based on ranges of key values.
*   **Example**: Transactions dated `2023-01` to `2023-06` go to Shard A; `2023-07` to `2023-12` go to Shard B.
*   **Pros**: Excellent for range queries (e.g., "Give me all records from March").
*   **Cons**: Susceptible to **Hotspots**. If everyone writes data with "today's date", one shard takes 100% of the load while others sit idle.

### 2. Directory / Geometric Sharding
Routes data based on a specific attribute, often utilizing a lookup table or geographic rules.
*   **Example**: Users with `region=US` go to a New York shard; `region=EU` go to a Frankfurt shard.
*   **Pros**: Low latency for local users; complies with data sovereignty laws.
*   **Cons**: High risk of uneven distribution if one region is significantly larger than others.

### 3. Hash-Based Sharding (Selected Strategy)
Uses a consistent hash function to randomize the shard assignment: `ShardID = Hash(Key) % TotalShards`.
*   **Example**: `CRC32("user:55")` results in 1,234, which maps to Shard 2.
*   **Pros**: **Uniform Distribution**. Random keys are spread evenly, preventing hotspots.
*   **Cons**: Range queries are impossible without querying *all* shards.

> **Why this project uses Hash-Based Sharding**:
> We assume a random distribution of User IDs. Hashing ensures that even if we insert `user:1`, `user:2`, `user:3` sequentially, they will likely land on different shards, utilizing the full capacity of the cluster.

---

## How This Project Simulates the Theory

This code models a sharded cluster in a single process using multiple Go structs.

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **Physical Shard** | `ShardNode` struct. It acts as an independent database with its own `map` and Mutex lock. |
| **Shard Router** | `RouterComponent` struct. It holds references to all shards and implements the routing logic. |
| **Partitioning Logic** | `crc32.ChecksumIEEE(key) % numShards`. A deterministic way to map strings to integers `[0..N]`. |
| **Data Skew** | The final bar chart output. Since we use random input keys, you will see some shards hold slightly more data than others (variance). |

### Key Files
*   `node.go`: Defines the storage unit (`ShardNode`). It focuses on storage (`Write`/`Read`) and thread-safety.
*   **`router.go`**: The brain of the operation. Contains the `GetShard` function which implements the hashing strategy.
*   `main.go`: The simulation runner. It generates random traffic (`user:1`...`user:1000`) and prints the visualization.

---

## How to Run

Ensure you have Go installed.

```bash
cd Sharding
go run .
```

### Sample Output Explanation

**1. Concurrent Reads & Writes**
```text
[Write] Key: user:10 -> time=... msg="Routing request" ... shard_id=2
       [Read Hit]  Key: user:3 (from Shard 1) => data-3
```
> **Real-time Traffic**: You will see green **[Write]** logs as data is ingested, and cyan **[Read Hit]** logs as the background readers successfully find the data on the correct shards.

**2. Distribution Visualization**
```text
=== Shard Distribution Stats ===
Shard 0: [█████████           ] 12 keys
Shard 1: [████████████        ] 13 keys
Shard 2: [█████████████       ] 12 keys
Shard 3: [████████████████    ] 13 keys
================================
```
> **Uniformity**: Notice that while not *perfectly* equal, the load is distributed across all 4 nodes without any manual intervention. This proves the hashing strategy works for random input.
