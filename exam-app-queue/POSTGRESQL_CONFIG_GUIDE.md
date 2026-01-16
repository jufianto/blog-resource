# PostgreSQL Configuration Guide

A comprehensive guide to understanding PostgreSQL configuration for resource-constrained environments (512MB RAM).

---

## Table of Contents

1. [Critical Memory Configurations](#critical-memory-configurations)
2. [Performance Tuning for SSD](#performance-tuning-for-ssd)
3. [Write-Ahead Log (WAL) Settings](#write-ahead-log-wal-settings)
4. [Configuration Summary](#configuration-summary)

---

## Critical Memory Configurations

These settings directly control PostgreSQL's memory usage and must be carefully tuned for servers with limited RAM.

### 1. `shared_buffers` - PostgreSQL Cache

**Default:** 128MB  
**Our setting:** 128MB  
**Recommendation:** 25% of total RAM

**What it does:**
- Shared memory cache for table and index data
- All connections share this pool
- Postgres checks here before reading from disk

**Impact:**

```
Query: SELECT * FROM exam_submissions WHERE user_id = 123;

First execution:
  → Reads from disk → stores in shared_buffers → 100ms

Second execution:
  → Reads from shared_buffers → 1ms (100x faster!)
```

**Sizing guidelines:**

| RAM | shared_buffers | Reasoning |
|-----|----------------|-----------|
| 512MB | 128MB (25%) | Our setup - balanced |
| 2GB | 512MB (25%) | Standard server |
| 8GB | 2GB (25%) | Production server |
| 32GB | 8GB (max) | Beyond 8GB, OS cache is more effective |

**Too low (32MB):**
- Excessive disk I/O
- Slow query performance
- Cache misses on hot data

**Too high (400MB on 512MB server):**
- Starves OS and other processes
- Risk of OOM
- Diminishing returns

---

### 2. `work_mem` - Per-Operation Memory

**Default:** 4MB  
**Our setting:** 2MB (reduced for high concurrency)  
**Warning:** This is per operation, not per connection!

**What it does:**
- Memory allocated for each sort/hash operation in queries
- Used by: ORDER BY, JOIN, GROUP BY, DISTINCT, window functions

**Operations that consume work_mem:**

```sql
-- Query with 5 operations
SELECT 
  user_id, 
  AVG(score) as avg_score
FROM exam_submissions
WHERE created_at > '2024-01-01'
GROUP BY user_id           -- Operation 1: uses work_mem
HAVING AVG(score) > 70     -- Operation 2: uses work_mem
ORDER BY avg_score DESC    -- Operation 3: uses work_mem
LIMIT 10;

-- Memory used: 3 operations × 2MB = 6MB
```

**Impact of too low (1MB):**
```
Sort operation > 1MB
→ Spills to disk (/tmp)
→ 10-100x slower
→ "temporary file" warnings in logs
```

**Impact of too high (64MB):**
```
Complex query: 5 operations × 64MB = 320MB
10 concurrent queries = 3.2GB memory usage
Server RAM = 512MB
Result: OOM crash!
```

---

### 3. `max_connections` × `work_mem` = Memory Risk

**Default max_connections:** 100  
**Our setting:** 100

**The Danger Formula:**
```
Max possible memory = max_connections × operations × work_mem
                    = 100 × 5 × 2MB
                    = 1000MB (1GB)
```

**Real-world scenarios:**

| Scenario | Connections | Operations/query | Memory Used | Status |
|----------|-------------|------------------|-------------|--------|
| Simple queries | 50 | 0 (no sort/join) | 0MB | ✅ Safe |
| Medium load | 30 | 1 (ORDER BY) | 60MB | ✅ Safe |
| Heavy analytics | 50 | 5 (complex) | 500MB | ⚠️ Risky |
| Peak disaster | 100 | 5 (complex) | 1000MB | ❌ CRASH |

**Example disaster scenario:**
```sql
-- All 100 users run this at peak exam time:
SELECT 
  u.name,
  COUNT(DISTINCT e.id),          -- Op 1
  AVG(e.score)                   -- Op 2
FROM exam_submissions e
JOIN users u ON e.user_id = u.id -- Op 3 (hash join)
GROUP BY u.name                   -- Op 4
ORDER BY avg_score DESC;          -- Op 5

Memory: 100 connections × 5 ops × 2MB = 1GB
Server RAM: 512MB
→ PostgreSQL killed by Linux OOM killer
```

**Why work_mem=2MB is defensive:**
- Assumes not all 100 connections active
- Most queries are simple (few operations)
- Complex queries don't all run simultaneously
- Better slow queries than crashed server

---

### 4. `effective_cache_size` - Query Planner Hint

**Default:** 4GB (varies by system)  
**Our setting:** 256MB (50% of 512MB RAM)  
**Important:** Does NOT allocate memory!

**What it does:**
- Tells query planner how much OS cache is available
- Influences index vs sequential scan decisions
- Pure hint - doesn't actually use RAM

**Impact of wrong value:**

**Scenario: Set to 4GB on 512MB server**
```
Planner thinks: "4GB cache available, use complex index plans"
Reality: Only 256MB available
Result: Slower queries due to cache misses
```

**Scenario: Set correctly to 256MB**
```
Planner thinks: "Limited cache, make conservative decisions"
Reality: Matches actual cache
Result: Optimal query plans
```

**If you don't set it:**
- Postgres uses compiled default (~4GB)
- Query planner makes wrong assumptions
- Performance suffers silently

---

## Performance Tuning for SSD

These settings optimize Postgres for modern SSD storage instead of legacy HDDs.

### 5. `random_page_cost` - Disk Access Cost Estimate

**Default:** 4.0 (assumes HDD)  
**Our setting:** 1.1 (optimized for SSD)

**What it does:**
- Cost factor for random disk reads
- Query planner uses this to choose between sequential scan vs index scan

**How index scans work:**

**Index structure (B-Tree):**
```
              [Root: 1-1000000]
              /              \
        [1-500000]        [500001-1000000]
        /        \           /           \
   [1-250k]  [250k-500k] [500k-750k]  [750k-1M]
     ...         ...         ...          ...
   [user_id=123 → page 4567]  ← Points to data page
```

**Finding user_id = 123:**
```
1. Read root page       → "123 in left branch"
2. Read branch page     → "123 in left-left branch"  
3. Read leaf page       → "user_id=123 on page 4567"
4. Jump to page 4567    ← RANDOM READ
5. Fetch row data       ← RANDOM READ

Total: ~5 random reads (NOT linear search!)
```

**Why it's called "random":**

**Sequential Scan (seq_page_cost=1.0):**
```
Disk: [1][2][3][4][5][6][7][8]...
Read:  →→→→→→→→→→→
Smooth, predictable
```

**Index Scan (random_page_cost):**
```
Disk: [1][2][3][4][5][6][7][8]...[5000]
Index points to: page 127, 2034, 4567
Read order: 127 → 2034 → 4567
        →→→ JUMP ←←← JUMP →→→

Scattered/random locations
```

**HDD impact:**
- Disk head physically moves
- Random seeks = 10-100x slower than sequential
- Planner avoids indexes for large result sets

**SSD impact:**
- No moving parts
- Random reads ≈ sequential reads
- Planner can use indexes more aggressively

---

**Query planner calculation example:**

**Table:** 1,000,000 rows (100,000 pages)

```sql
SELECT * FROM exam_submissions WHERE user_id = 123;
-- Returns 10 rows
```

**Option A: Sequential Scan**
```
Cost = 100,000 pages × 1.0 = 100,000
```

**Option B: Index Scan with random_page_cost=4.0 (HDD)**
```
Index reads: 5 pages × 4.0 = 20
Heap reads:  10 pages × 4.0 = 40
Total: 60 < 100,000 → Choose index ✓
```

**Option C: Index Scan with random_page_cost=1.1 (SSD)**
```
Index reads: 5 pages × 1.1 = 5.5
Heap reads:  10 pages × 1.1 = 11
Total: 16.5 < 100,000 → Confidently choose index ✓
```

---

**Impact on query plans:**

**Returning 30 rows (0.003%):**

```sql
EXPLAIN SELECT * FROM exam_submissions 
WHERE user_id IN (123, 456, 789);
```

**With random_page_cost=4.0:**
```
Seq Scan on exam_submissions  (cost=0.00..25000.00)
  Filter: (user_id = ANY ('{123,456,789}'))
  
→ Reads ALL 1M rows, filters in memory
→ Slow
```

**With random_page_cost=1.1:**
```
Index Scan using idx_user_id  (cost=0.42..68.50)
  Index Cond: (user_id = ANY ('{123,456,789}'))
  
→ Reads only 30 rows via index
→ 1000x faster!
```

---

**When index is NOT used (even with SSD setting):**

```sql
SELECT * FROM exam_submissions WHERE score > 50;
-- Returns 500,000 rows (50% of table)
```

**Sequential Scan:**
```
Cost = 100,000 pages × 1.0 = 100,000
```

**Index Scan:**
```
Index + heap reads: 500,000 pages × 1.1 = 550,000
Still more expensive than sequential!
→ Planner correctly chooses sequential scan
```

**Rule of thumb:**
- Index used: returning <5-10% of rows
- Sequential scan: returning >10% of rows
- SSD setting makes this threshold more favorable to indexes

---

### 6. `effective_io_concurrency` - Parallel I/O

**Default:** 1 (single HDD)  
**Our setting:** 200 (NVMe SSD)

**What it does:**
- Number of simultaneous I/O operations disk can handle
- Used for bitmap heap scans and parallel queries

**Impact:**

**With effective_io_concurrency=1 (HDD):**
```sql
SELECT * FROM exam_submissions 
WHERE score > 80 OR created_at > '2024-01-01';

Execution:
  Read block 1 → wait
  Read block 2 → wait
  Read block 3 → wait
  ...
  
Duration: 5 seconds
```

**With effective_io_concurrency=200 (SSD):**
```sql
Same query

Execution:
  Issue 200 read requests in parallel
  SSD handles them simultaneously
  
Duration: 1 second (5x faster)
```

**Best for:**
- Bitmap heap scans (OR conditions, IN clauses)
- Parallel sequential scans
- Queries reading scattered data

---

## Write-Ahead Log (WAL) Settings

WAL ensures data durability but impacts write performance.

### 7. `wal_buffers`

**Default:** Typically 16MB or -1 (auto-sized)  
**Our setting:** 16MB

**What it does:**
- Memory buffer for WAL writes
- Reduces disk I/O for transaction logs

**Impact:**
- Small writes batched in memory
- Flushed to disk at commit
- Larger = fewer disk writes

---

### 8. `min_wal_size` & `max_wal_size`

**Defaults:** 80MB / 1GB  
**Our setting:** 1GB / 2GB

**What it does:**
- Controls checkpoint frequency
- Larger = less frequent checkpoints = better write performance
- Trade-off: longer recovery time after crash

---

### 9. `checkpoint_completion_target`

**Default:** 0.9  
**Our setting:** 0.9

**What it does:**
- Spreads checkpoint I/O over time
- 0.9 = use 90% of checkpoint interval
- Prevents I/O spikes

---

### 10. `maintenance_work_mem`

**Default:** 64MB  
**Our setting:** 64MB

**What it does:**
- Memory for VACUUM, CREATE INDEX, ALTER TABLE
- Per operation, not shared

---

### 11. `default_statistics_target`

**Default:** 100  
**Our setting:** 100

**What it does:**
- Amount of statistics collected per column
- Higher = better query plans, slower ANALYZE

---

## Configuration Summary

### Our Full Config (docker-compose.yaml)

```yaml
postgres:
  image: postgres:16
  mem_limit: 512m
  cpus: 0.5
  command: >
    postgres
    -c max_connections=100              # Default: 100
    -c shared_buffers=128MB             # Default: 128MB (25% of RAM)
    -c effective_cache_size=256MB       # Default: 4GB → Set to 50% of RAM
    -c maintenance_work_mem=64MB        # Default: 64MB
    -c checkpoint_completion_target=0.9 # Default: 0.9
    -c wal_buffers=16MB                 # Default: 16MB
    -c default_statistics_target=100    # Default: 100
    -c random_page_cost=1.1             # Default: 4.0 → SSD optimized
    -c effective_io_concurrency=200     # Default: 1 → SSD optimized
    -c work_mem=2MB                     # Default: 4MB → Reduced for concurrency
    -c min_wal_size=1GB                 # Default: 80MB
    -c max_wal_size=2GB                 # Default: 1GB
```

---

### Priority Matrix

| Priority | Setting | Why Changed | Impact |
|----------|---------|-------------|--------|
| 🔴 Critical | `work_mem=2MB` | Prevent OOM with 100 connections | High concurrency safety |
| 🔴 Critical | `effective_cache_size=256MB` | Match actual RAM | Correct query plans |
| 🟡 Important | `random_page_cost=1.1` | SSD storage | 10-1000x faster queries |
| 🟡 Important | `effective_io_concurrency=200` | SSD parallel I/O | 2-5x faster scans |
| 🟢 Performance | `min/max_wal_size` | Reduce checkpoint frequency | Better write performance |
| 🟢 Performance | `checkpoint_completion_target=0.9` | Smooth I/O | Avoid I/O spikes |

---

### Quick Reference: When to Adjust

**Got more RAM (upgrade to 2GB)?**
```
shared_buffers: 128MB → 512MB
effective_cache_size: 256MB → 1GB
work_mem: 2MB → 4MB (if reducing max_connections)
```

**Need higher concurrency (200 connections)?**
```
work_mem: 2MB → 1MB (or reduce max_connections)
shared_buffers: Keep at 128MB
Consider connection pooling (PgBouncer)
```

**Using HDD instead of SSD?**
```
random_page_cost: 1.1 → 4.0
effective_io_concurrency: 200 → 1
```

**Write-heavy workload?**
```
wal_buffers: 16MB → 32MB
max_wal_size: 2GB → 4GB
```

---

## Testing Your Configuration

**Check current values:**
```sql
SHOW shared_buffers;
SHOW work_mem;
SHOW max_connections;
SHOW effective_cache_size;
SHOW random_page_cost;
```

**Monitor memory usage:**
```sql
SELECT 
  setting,
  unit,
  boot_val,
  reset_val
FROM pg_settings
WHERE name IN (
  'shared_buffers',
  'work_mem',
  'effective_cache_size'
);
```

**Check for disk spills (work_mem too low):**
```sql
-- Enable logging
SET log_temp_files = 0;

-- Look for "temporary file" in logs
-- Increase work_mem if frequent
```

**Verify query plans use indexes:**
```sql
EXPLAIN (ANALYZE, BUFFERS) 
SELECT * FROM exam_submissions WHERE user_id = 123;

-- Look for "Index Scan" vs "Seq Scan"
-- Check "Buffers: shared hit" (cache hits)
```

---

## Common Mistakes

❌ **Setting work_mem too high**
```
work_mem=64MB with max_connections=100
→ Potential 6.4GB usage on 512MB server
→ OOM crash
```

❌ **Not adjusting effective_cache_size**
```
Leave at default 4GB on 512MB server
→ Query planner thinks unlimited cache
→ Suboptimal query plans
```

❌ **Using HDD settings on SSD**
```
random_page_cost=4.0 on SSD
→ Planner avoids indexes unnecessarily
→ Slow queries
```

❌ **Forgetting the multiplication effect**
```
"2MB is small, surely safe"
→ 100 connections × 5 operations = 1GB
→ Server crash
```

---

## Further Reading

- [PostgreSQL Tuning Guide](https://wiki.postgresql.org/wiki/Tuning_Your_PostgreSQL_Server)
- [PgTune - Configuration Calculator](https://pgtune.leopard.in.ua/)
- [Explain Plan Visualization](https://explain.dalibo.com/)

---

**Related:** See [README.md](./README.md) for the full exam app queue architecture and load testing guide.
