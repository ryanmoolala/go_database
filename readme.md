# bptree
A read-optimized B+ Tree implementation in Go, designed as the foundation for a high-performance cache and storage layer.

## Overview

`bptree` implements a fully functional B+ Tree optimized for **read-heavy workloads**. The long-term goal is to evolve this into a storage layer suitable for caching systems, database indexing, and key-value engines.

B+ Trees are the index structure of choice in systems like MySQL InnoDB and PostgreSQL because they keep all data in leaf nodes (enabling efficient sequential scans), maintain a balanced height for consistent `O(log n)` lookups, and expose a natural API for range queries via linked leaves.

## Features
- [x] Insert, search
- [x] Delete
- [x] Balanced tree with automatic node splitting, underflow & overflow
- [ ] Duplicate keys
- [x] Linked leaf nodes for range scans
- [x] Bulk loading
- [x] Range scan API
- [ ] Persistence layer (disk-backed nodes)
- [ ] Concurrency support (RW locks, latch coupling)
- [ ] Cache eviction policies (LRU, LFU)
- [ ] Benchmarks (read vs write throughput)


## Basic usage
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


## Planned architecture
```
┌─────────────────────────────┐
│        Cache Layer          │  
├─────────────────────────────┤
│        B+ Tree Index        │  
├─────────────────────────────┤
│      Storage Backend        │ 
└─────────────────────────────┘
```


## Roadmap

**Phase 1 — Core data structure** *(current)*
- B+ Tree with insert, search, delete
- Node splitting and merging
- Linked leaves for range traversal
