# Load Test Results: Direct vs Queue Approach

## Test Configuration
- **Duration**: 5 minutes
- **Peak Load**: 3000 concurrent users
- **Total Requests**: ~413,000

## 📊 Resource Usage Comparison

### PostgreSQL (exam-db)

| Metric | Direct Approach | Queue Approach | Difference |
|--------|----------------|----------------|------------|
| **Avg CPU** | 26.97% | 12.09% | -14.88% |
| **Max CPU** | 51.05% | 42.69% | -8.36% |
| **Avg Memory** | 197 MB | 198 MB | +1 MB |
| **Max Memory** | 294 MB | 223 MB | -71 MB |
| **Avg Mem %** | 38.41% | 38.67% | +0.26% |

### NATS (exam-nats) - Queue Only

| Metric | Value |
|--------|-------|
| **Avg CPU** | 43.57% |
| **Max CPU** | 92.52% |
| **Avg Memory** | 30 MB |
| **Max Memory** | 42 MB |

## 🎯 Key Insights

### Direct Approach
- ✅ Lower latency (direct write to database)
- ⚠️ Higher database load during peaks
- ⚠️ No buffering - all requests hit DB immediately

### Queue Approach
- ✅ Decoupled - API responds faster
- ✅ Smoother database load (batched writes)
- ✅ Better resilience (messages queued if DB slow)
- ⚠️ Additional component (NATS) to manage
- ⚠️ Small overhead from message passing

## 📈 Visualization

![Resource Usage Comparison](comparison_chart.png)

---
*Generated automatically from load test results*
