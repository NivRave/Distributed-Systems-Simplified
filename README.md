# Distributed Systems Simplified

A comprehensive deep dive into distributed systems concepts, implemented in Go. This repository breaks down complex theories into clean, professional, and simulate-able code projects.

## About the Author
**Niv Rave**
*   LinkedIn: [https://www.linkedin.com/in/niv-rave/](https://www.linkedin.com/in/niv-rave/)
*   Email: nivikr@gmail.com

---

## Project Roadmap

Each topic is implemented as an independent project in its own folder, complete with a simulation and explanation.

### **Phase 1: Data Distribution & Scaling (The Storage Foundation)**

* **[Read Replication](./Read%20Replication/)** ✅
    * *Concept*: Scaling read-heavy workloads using a Leader-Follower model.
    * *Simulation*: Asynchronous data propagation with simulated replication lag.
* **[Sharding](./Sharding/)** ✅
    * *Concept*: Horizontal partitioning to scale write throughput and manage massive datasets.
    * *Simulation*: A router distributing unique keys across independent physical shards.
* **[Consistent Hashing](./Consistent%20Hashing/)** ✅
    * *Concept*: Minimizing data movement when scaling a cluster (adding/removing nodes).
    * *Simulation*: Visualization of a Ring structure, Virtual Nodes, and efficient rebalancing stats.
* **[CAP Theorem](./CAP%20Theorem/)** ✅
    * *Concept*: The fundamental trade-off: in a Partition (P), choose Consistency (CP) or Availability (AP).
    * *Simulation*: A partitioned network where AP nodes accept writes while CP nodes reject them.

### **Phase 2: Performance & Traffic Management (The Gateway)**

* **[Distributed Caching](./Distributed%20Caching/)** ✅
    * *Concept*: Reducing database pressure and latency by storing "hot" data in-memory.
    * *Simulation*: Implementing LRU eviction and TTL expiration, comparing Cache Hit vs. Cache Miss latency.
* **[Load Balancing & Health Checks](./Load%20Balancing/)** ✅
    * *Concept*: Intelligent traffic routing and automatic failover for high availability.
    * *Simulation*: A dynamic balancer that detects node crashes and reroutes traffic in real-time.
* **Rate Limiting (Token Bucket)** 🆕
    * *Concept*: Protecting system resources from abuse and ensuring fair usage (Noisy Neighbor problem).
    * *Simulation*: Using Go channels and tickers to simulate a request "refill" bucket.

### **Phase 3: Service Resiliency (The Bodyguard)**

* **The Circuit Breaker** 🆕
    * *Concept*: Preventing cascading failures by "tripping the fuse" when a downstream service is struggling.
    * *Simulation*: A state machine (Closed/Open/Half-Open) protecting a fragile mock service.
* **API Gateway & Middleware** 🆕
    * *Concept*: Centralizing cross-cutting concerns like Authentication, Logging, and Request Routing.
    * *Simulation*: A reverse-proxy that validates requests before they reach the internal service mesh.
* **Idempotency & Retries** 🆕
    * *Concept*: Ensuring distributed operations (like payments) are safe to retry without side effects.
    * *Simulation*: Using Unique Request IDs to skip duplicate processing in a multi-service flow.

### **Phase 4: Asynchronous Coordination (The Event-Driven World)**

* **Message Queues (Pub/Sub)** 🆕
    * *Concept*: Decoupling services for non-blocking, high-performance background processing.
    * *Simulation*: A Producer/Consumer model using Go channels to simulate an async job worker.
* **The Saga Pattern** 🆕
    * *Concept*: Managing long-running distributed transactions with compensating (Undo) logic.
    * *Simulation*: Orchestrating a "Trip Booking" (Flight + Hotel) that rolls back if one part fails.
* **Distributed Tracing (Context Propagation)** 🆕
    * *Concept*: Visualizing the "Life of a Request" as it travels through multiple microservices.
    * *Simulation*: Injecting and extracting Trace IDs using Go’s `context` package across service boundaries.
