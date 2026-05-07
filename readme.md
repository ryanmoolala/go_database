# bptree
A read-optimized B+ Tree implementation in Go, designed as the foundation for a high-performance cache and storage layer.

## Overview

`bptree` implements a fully functional B+ Tree optimized for **read-heavy workloads**. The long-term goal is to evolve this into a storage layer suitable for caching systems, database indexing, and key-value engines.

B+ Trees are the index structure of choice in systems like MySQL InnoDB and PostgreSQL because they keep all data in leaf nodes (enabling efficient sequential scans), maintain a balanced height for consistent `O(log n)` lookups, and expose a natural API for range queries via linked leaves.

## Features
- [x] Insert, search
- [x] Delete
- [x] Search with conditions
- [x] Balanced tree with automatic node splitting, underflow & overflow
- [ ] Duplicate keys
- [x] Linked leaf nodes for range scans
- [x] Bulk loading
- [x] Range scan API
- [x] Interactive CLI with commands
- [ ] Persistence layer (disk-backed nodes)
- [ ] Concurrency support (RW locks, latch coupling)
- [ ] Cache eviction policies (LRU, LFU)
- [ ] Benchmarks (read vs write throughput)


## Basic usage
```go
CREATE DUMMYTABLE 3

INSERT DUMMYTABLE VALUES ("TEST", "VALUE")
INSERT DUMMYTABLE VALUES ("TEST2", "VALUE2")

DELETE FROM DUMMYTABLE "TEST"

PRINT DUMMYTABLE

SELECT * FROM DUMMYTABLE WHERE KEY >= "TEST"
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

