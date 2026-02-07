# Distributed Systems Simplified

A comprehensive deep dive into distributed systems concepts, implemented in Go. This repository breaks down complex theories into clean, professional, and simulate-able code projects.

## About the Author
**Niv Rave**
*   LinkedIn: [https://www.linkedin.com/in/niv-rave/](https://www.linkedin.com/in/niv-rave/)
*   Email: nivikr@gmail.com

---

## Project Roadmap

Each topic is implemented as an independent project in its own folder, complete with a simulation and explanation.

### Phase 1: Data Distribution & Scaling
*   **[Read Replication](./Read%20Replication/)** ✅ (Completed)
    *   *Concept*: Scaling read-heavy workloads using a Leader-Follower model.
    *   *Simulation*: Leader DB writing updates that asynchronously propagate to follower nodes with replication lag.

*   **[Sharding](./Sharding/)** ✅ (Completed)
    *   *Concept*: Scaling write-heavy workloads and massive datasets by partitioning.
    *   *Simulation*: Consistent hashing router distributing distinct user keys across multiple shard nodes.

*   **Consistent Hashing** (Upcoming)
    *   *Concept*: Solving the "rebalancing" problem when adding/removing nodes in a sharded cluster.

### Phase 2: Reliability & Theory
*   **CAP Theorem** (Upcoming)
    *   *Concept*: The fundamental trade-off between Consistency, Availability, and Partition Tolerance.
*   **Distributed Transactions (2PC)** (Upcoming)
    *   *Concept*: How to ensure an "all or nothing" operation across multiple physical nodes.

### Phase 3: Coordination & Time
*   **Network Time Protocol (NTP) & Drift** (Upcoming)
    *   *Concept*: Why "Global Time" is an illusion in distributed systems.
*   **Vector Clocks** (Upcoming)
    *   *Concept*: Tracking causality—knowing which event happened before another without a master clock.
*   **Distributed Consensus (Paxos/Raft)** (Upcoming)
    *   *Concept*: How nodes agree on a single value (the "holy grail" of distributed logic).

### Phase 4: Big Data & Batch Processing
*   **MapReduce** (Upcoming)
    *   *Concept*: Building a Master-Worker architecture to process terabytes of data in parallel.
*   **Hadoop vs. Spark vs. Storm** (Upcoming)
    *   *Concept*: Exploring the evolution of the Big Data ecosystem.
*   **Lambda Architecture** (Upcoming)
    *   *Concept*: Combining batch and real-time speed layers for the best of both worlds.

### Phase 5: Messaging & Orchestration
*   **Messaging Introduction (Pub/Sub)** (Upcoming)
    *   *Concept*: Decoupling services using message queues.
*   **Kafka** (Upcoming)
    *   *Concept*: Building a persistent, append-only distributed log.
*   **Zookeeper** (Upcoming)
    *   *Concept*: Managing configuration and service discovery in a clustered environment.