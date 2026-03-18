# bptree

A read-optimized B+ Tree implementation in Go, designed as the foundation for a high-performance cache and storage layer.

## Overview

`bptree` implements a fully functional B+ Tree optimized for **read-heavy workloads**. The long-term goal is to evolve this into a storage layer suitable for caching systems, database indexing, and key-value engines.

B+ Trees are the index structure of choice in systems like MySQL InnoDB and PostgreSQL because they keep all data in leaf nodes (enabling efficient sequential scans), maintain a balanced height for consistent `O(log n)` lookups, and expose a natural API for range queries via linked leaves.

---

## Features

- [x] Insert, search
- [ ] Delete
- [ ] Balanced tree with automatic node splitting
- [ ] Linked leaf nodes for range scans
- [ ] Bulk loading
- [ ] Range scan API
- [ ] Persistence layer (disk-backed nodes)
- [ ] Concurrency support (RW locks, latch coupling)
- [ ] Cache eviction policies (LRU, LFU)
- [ ] Benchmarks (read vs write throughput)

---

## Getting Started

```bash
git clone https://github.com/your-username/bptree.git
cd bptree
go build ./...
go test ./...
```

### Basic usage

```go
tree := bptree.New(order)

tree.Insert(42, "hello")
tree.Insert(17, "world")

val, ok := tree.Search(42)
// val = "hello", ok = true

results := tree.RangeScan(10, 50)
// returns all key-value pairs with keys in [10, 50]

tree.Delete(17)
```

---

## Design

### Why B+ Tree over alternatives?

| Structure | Point lookup | Range scan | Write-heavy | Notes |
|---|---|---|---|---|
| B+ Tree | `O(log n)` | Excellent | Good | All data in leaves; cache-friendly |
| Hash map | `O(1)` | None | Excellent | No ordering |
| LSM Tree | `O(log n)` | Good | Excellent | Write-optimized; read amplification |
| Skip list | `O(log n)` | Good | Good | Simpler but higher memory overhead |

For read-heavy workloads with range query requirements, B+ Tree is the natural fit.

### Planned architecture

```
┌─────────────────────────────┐
│        Cache Layer          │  ← eviction policies, TTL
├─────────────────────────────┤
│        B+ Tree Index        │  ← sorted keys, range scans
├─────────────────────────────┤
│      Storage Backend        │  ← in-memory or disk-backed
└─────────────────────────────┘
```

---

## Roadmap

**Phase 1 — Core data structure** *(current)*
- B+ Tree with insert, search, delete
- Node splitting and merging
- Linked leaves for range traversal

**Phase 2 — Performance**
- Cache-locality optimizations
- Node traversal profiling
- Benchmarks against `map` and competing tree structures

**Phase 3 — Cache system**
- Key-value cache API on top of the tree
- In-memory + disk hybrid storage
- LRU and LFU eviction strategies
- Concurrency via RW locks and latch coupling
