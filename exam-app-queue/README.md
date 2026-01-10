# Exam App Queue - Performance Comparison

A Go application demonstrating the performance difference between **direct database inserts** vs **queue-based inserts** when handling high-concurrency exam submissions.

## 🎯 Purpose

This project simulates a real-world scenario: **10,000 students submitting exams simultaneously** (generating ~3,000 req/s). It compares two architectural approaches to show when and why message queues are essential.

## 📊 Two Approaches

### Approach 1: Direct Insert (No Queue)
```
HTTP Request → API Server → PostgreSQL → Response
```
- **Port**: 8080
- **Behavior**: Waits for database before responding
- **Expected**: Slow response times, connection pool exhaustion under load

### Approach 2: Queue-based Insert
```
HTTP Request → API Server → NATS Queue → Fast Response
                             ↓
                          Workers → PostgreSQL
```
- **Port**: 8081
- **Behavior**: Publishes to queue and responds immediately
- **Expected**: Fast response times, controlled database load

## 🏗️ Architecture

### Resource-Limited Environment

Docker containers are configured with limited resources to simulate real production servers (not M4 Pro power):

**PostgreSQL**:
- CPU: 1 core
- RAM: 512MB
- Max connections: 100

**NATS**:
- CPU: 0.5 core
- RAM: 256MB

## 📁 Project Structure

```
exam-app-queue/
├── docker-compose.yaml          # PostgreSQL + NATS with resource limits
├── go.mod
├── README.md
│
├── api/
│   ├── direct/main.go          # Approach 1: Direct DB insert
│   └── queue/main.go           # Approach 2: NATS queue
│
├── worker/main.go              # Background workers (queue consumers)
│
├── store/
│   ├── models.go               # Data structures
│   └── db.go                   # Database operations
│
├── migrate/
│   ├── migrate.go
│   └── migrations/
│       ├── 001_create_exam_submissions.up.sql
│       └── 001_create_exam_submissions.down.sql
│
├── load-test/
│   ├── test_direct.js          # k6 load test for direct approach
│   ├── test_queue.js           # k6 load test for queue approach
│   └── results/                # Test results
│
└── env/
    ├── env.yaml               # Your config
    └── sample.env.yaml        # Template
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23.6+
- Docker & Docker Compose
- [k6](https://k6.io/docs/get-started/installation/) for load testing

### Installation

1. **Start Infrastructure**

```bash
cd exam-app-queue
docker-compose up -d
```

Check containers:
```bash
docker-compose ps
```

2. **Configure Application**

```bash
cp env/sample.env.yaml env/env.yaml
```

Edit `env/env.yaml` if needed (defaults should work).

3. **Run Database Migrations**

```bash
cd migrate
go run migrate.go up
```

4. **Install Dependencies**

```bash
cd ..
go mod download
```

## 🎮 Running the Applications

### Approach 1: Direct Insert API

```bash
cd api/direct
go run main.go
```

Server starts on `:8080`

Test it:
```bash
curl -X POST http://localhost:8080/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "student_001",
    "exam_id": "math_101",
    "answers": {"q1": "A", "q2": "B", "q3": "C"}
  }'
```

### Approach 2: Queue-based API + Workers

**Terminal 1 - Start API**:
```bash
cd api/queue
go run main.go
```

Server starts on `:8081`

**Terminal 2 - Start Workers**:
```bash
cd worker
go run main.go
```

10 workers start consuming from NATS queue

Test it:
```bash
curl -X POST http://localhost:8081/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "student_001",
    "exam_id": "math_101",
    "answers": {"q1": "A", "q2": "B", "q3": "C"}
  }'
```

## 📈 Load Testing

### Install k6

macOS:
```bash
brew install k6
```

Other platforms: https://k6.io/docs/get-started/installation/

### Test Direct Approach

```bash
cd load-test
k6 run test_direct.js
```

This simulates:
- Ramp up to 3,000 concurrent users
- Sustained load at 3,000 req/s
- Measures response times, success rate, throughput

### Test Queue Approach

**Ensure workers are running!**

```bash
k6 run test_queue.js
```

### Compare Results

Results are saved in `load-test/results/`:
- `direct_summary.json`
- `queue_summary.json`

## 📊 Expected Results

### Direct Insert (Approach 1)

| Metric | Expected Value |
|--------|----------------|
| Avg Response Time | 300-800ms |
| P95 Response Time | 1000-2000ms |
| P99 Response Time | 2000-5000ms |
| Success Rate | 60-80% |
| Max Throughput | ~500 req/s |
| DB Connection Pool | Exhausted (100/100) |

**Symptoms**:
- Slow responses
- Timeout errors
- Connection pool exhaustion
- Database CPU at 100%

### Queue-based Insert (Approach 2)

| Metric | Expected Value |
|--------|----------------|
| Avg Response Time | 5-20ms |
| P95 Response Time | 30-60ms |
| P99 Response Time | 80-150ms |
| Success Rate | 99%+ |
| Max Throughput | 3000+ req/s |
| DB Connection Pool | Healthy (10-20/100) |

**Benefits**:
- Fast API responses
- No timeout errors
- Controlled DB load
- Queue absorbs traffic spikes

## 🔍 Monitoring

### NATS Monitoring

Visit: http://localhost:8222

- `/varz` - Server info
- `/connz` - Connection info
- `/jsz` - JetStream info

### Database Queries

```bash
docker exec -it exam-db psql -U postgres -d exam_app
```

```sql
-- Check submission count
SELECT status, COUNT(*) FROM exam_submissions GROUP BY status;

-- Recent submissions
SELECT * FROM exam_submissions ORDER BY submitted_at DESC LIMIT 10;

-- Processing lag
SELECT 
  AVG(EXTRACT(EPOCH FROM (processed_at - submitted_at))) as avg_lag_seconds
FROM exam_submissions 
WHERE processed_at IS NOT NULL;
```

## 🧪 Database Schema

```sql
CREATE TABLE exam_submissions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    exam_id VARCHAR(50) NOT NULL,
    answers JSONB NOT NULL,
    score DECIMAL(5, 2),
    submitted_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Status values**:
- `pending` - In queue (approach 2)
- `processed` - Saved to database
- `failed` - Error occurred

## 🎓 Key Learnings

### When to Use Direct Insert
✅ Low traffic (< 100 req/s)  
✅ Simple CRUD operations  
✅ Single server setup  
✅ Immediate consistency required  

### When to Use Queue
✅ High concurrent writes (> 500 req/s)  
✅ Traffic spikes expected  
✅ User-facing applications  
✅ Need fast response times  
✅ Want retry capability  
✅ Need to scale independently  

## 🛠️ Configuration

### Worker Configuration

Edit `env/env.yaml`:

```yaml
worker:
  count: 10          # Number of workers
  batch_size: 100    # Records per batch insert
```

**Tuning**:
- More workers = higher throughput (but more DB connections)
- Larger batch = fewer DB calls (but higher memory)
- Sweet spot: 10-20 workers, 50-100 batch size

### Database Pool

In API/Worker code:

```go
configPool.MaxConns = 50   // Max connections
configPool.MinConns = 10   // Idle connections
```

**Rule**: MaxConns ≥ number of workers + API requests

## 🐛 Troubleshooting

### Connection refused to PostgreSQL

```bash
docker-compose ps
docker-compose up -d postgres
```

### NATS connection failed

```bash
docker-compose logs nats
docker-compose restart nats
```

### Workers not processing

1. Check NATS is running
2. Verify queue API published messages
3. Check worker logs for errors

### Load test fails

1. Ensure correct API is running (8080 or 8081)
2. Start workers for queue approach
3. Check Docker containers have enough resources

## 📝 Clean Up

```bash
# Stop containers
docker-compose down

# Remove data volumes
docker-compose down -v

# Reset database
cd migrate
go run migrate.go down
```

## 🎯 Blog Post Topics

This project demonstrates:

1. **Connection pooling importance** (pgxpool)
2. **Queue vs direct insert trade-offs**
3. **Handling traffic spikes**
4. **Backpressure management**
5. **Worker pool patterns**
6. **Batch insert optimization**
7. **Load testing with k6**
8. **Resource-constrained environments**

## 📚 Technologies Used

- **Go** - Application code
- **PostgreSQL** - Database
- **NATS JetStream** - Message queue
- **pgx/pgxpool** - PostgreSQL driver & connection pooling
- **k6** - Load testing
- **Docker** - Container orchestration
- **Viper** - Configuration management

## 📄 License

This is a learning/blog resource project. Feel free to use for educational purposes.

## 🤝 Contributing

This is a demonstration project for blog content. Suggestions welcome!

---

**Happy Learning! 🚀**
