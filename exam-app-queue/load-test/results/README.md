# Load Test Results Directory

This directory will contain k6 test results.

After running tests, you'll find:
- `direct_summary.json` - Results from direct insert approach
- `queue_summary.json` - Results from queue-based approach

## Running Tests

```bash
# Test direct approach
k6 run test_direct.js

# Test queue approach (ensure workers are running!)
k6 run test_queue.js
```

## Comparing Results

Look for these key metrics:

### Response Times
- avg, p95, p99
- Direct should be 200-800ms
- Queue should be < 20ms

### Success Rate  
- Direct: 60-80%
- Queue: 99%+

### Throughput
- Direct: ~500 req/s max
- Queue: 3000+ req/s
