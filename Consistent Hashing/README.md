# Consistent Hashing

This project demonstrates the concept of **Consistent Hashing**, a special kind of hashing such that when a hash table is resized, only $k/n$ keys need to be remapped on average.

## What is Consistent Hashing?

Consistent Hashing solves the **rebalancing problem** in distributed systems. When you have a cluster of 10 nodes and you add an 11th node, you don't want to shuffle 100% of your data. You only want to move ~1/11th of the data to the new node.

*   **The Ring**: We map both Nodes and Keys to a 32-bit integer circle.
*   **Assignment**: A key belongs to the first node found moving clockwise on the ring.

## Why use it?
1.  **Minimize Data Movement**: Adding/Removing nodes causes minimal disruption.
2.  **Horizontal Scalability**: Allows dynamic scaling without "stopping the world" to rebalance.
3.  **High Availability**: If a node dies, only its keys are reassigned to its neighbor.

### The Innovation: Virtual Nodes
With only 3 physical nodes (A, B, C), they might be clumped together on the ring, causing Node A to hold 90% of the data.
*   **Solution**: We hash each physical node multiple times (e.g., "Node A#1", "Node A#2"...).
*   **Result**: Nodes are statistically scattered evenly around the ring.

---

## How This Project Simulates the Theory

This code models the Ring structure and simulates a scaling event.

| Theory Concept | Simulation Implementation |
| :--- | :--- |
| **Hash Space (Ring)** | `uint32` integer space (0 to 4 Billion). |
| **Virtual Node** | We hash strings like `Node-A#1`, `Node-A#2` and store them in a sorted slice. |
| **Node Lookup** | `sort.Search` (Binary Search) to find the nearest node >= KeyHash. |
| **Scaling Event** | We calculate the distribution with 3 nodes, then add a 4th, and count exactly how many keys moved. |

### Key Files
*   `ring.go`: Implements the core logic (Adding nodes, Removing nodes, Finding the owner of a key).
*   `node.go`: Represents a physical storage unit (just holds an ID in this demo).
*   `simulation.go`: The scenario runner. It generates 2000 keys, maps them, adds a node, and calculates migration stats.

---

## How to Run

Ensure you have Go installed.

```bash
cd "Consistent Hashing"
go run .
```

### Sample Output Explanation

**1. Initial Distribution (3 Nodes)**
```text
Node-A  : [████████                 ]  678 keys (33.9%)
Node-B  : [████████                 ]  651 keys (32.5%)
Node-C  : [████████                 ]  671 keys (33.6%)
```
> **Balance**: Thanks to Virtual Nodes, data is split evenly (~33% each).

**2. Scaling Up (Adding Node-D)**
```text
[Step 2] Adding Node-D to the cluster...
  -> Key user-session-12 moved from Node-C to Node-D
  -> Key user-session-19 moved from Node-B to Node-D
  ... (total 492 keys moved)
```
> **Minimal Movement**: Only keys that "fall" into Node-D's new ranges are moved. Keys on Node-A that are far away on the ring stay untouched.

**3. Final Analysis**
```text
[Analysis]
Total Moved: 492 keys (24.60%)
Ideal Movement: 25.00% (1/4 of keys)
Result: Very efficient rebalancing!
```
> **Efficiency**: In a naive system, 75% of keys would move. Here, only ~25% moved, which is mathematically perfect for moving from 3 to 4 nodes.
