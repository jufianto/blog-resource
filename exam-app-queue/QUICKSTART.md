# Quick Start Guide

## 1. Start Infrastructure (1 minute)

```bash
cd exam-app-queue
docker-compose up -d
```

Wait ~10 seconds for PostgreSQL to be ready.

## 2. Run Migrations (30 seconds)

```bash
cd migrate
go run migrate.go up
cd ..
```

## 3. Test Both Approaches

### Option A: Direct Insert (Simple Test)

**Terminal 1:**
```bash
cd api/direct
go run main.go
```

**Terminal 2:**
```bash
curl -X POST http://localhost:8080/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "student_001",
    "exam_id": "math_101",
    "answers": {"q1": "A", "q2": "B", "q3": "C"}
  }'
```

You should see response in ~10-50ms (light load).

### Option B: Queue-based (Two Terminals Needed)

**Terminal 1 - Start API:**
```bash
cd api/queue
go run main.go
```

**Terminal 2 - Start Workers:**
```bash
cd worker
go run main.go
```

**Terminal 3 - Test:**
```bash
curl -X POST http://localhost:8081/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "student_001",
    "exam_id": "math_101",
    "answers": {"q1": "A", "q2": "B", "q3": "C"}
  }'
```

You should see response in ~2-10ms (super fast!).

## 4. Load Testing (The Fun Part!)

Install k6:
```bash
brew install k6
```

### Test Direct Approach

**Keep api/direct running**, then:

```bash
cd load-test
k6 run test_direct.js
```

Watch as it struggles under load! You'll see:
- Response times increase to 500-2000ms
- Error rates climb to 20-40%
- Database connections max out

### Test Queue Approach  

**Keep api/queue AND workers running**, then:

```bash
cd load-test
k6 run test_queue.js
```

Watch it handle the load easily! You'll see:
- Response times stay < 20ms
- Error rates stay < 1%
- Smooth processing

## 5. Monitor NATS (Optional)

Visit http://localhost:8222

## 6. Check Database

```bash
docker exec -it exam-db psql -U postgres -d exam_app

-- Inside psql:
SELECT status, COUNT(*) FROM exam_submissions GROUP BY status;
SELECT * FROM exam_submissions ORDER BY submitted_at DESC LIMIT 5;
```

## Common Issues

### "Connection refused" to database
```bash
docker-compose ps
# If postgres is not running:
docker-compose up -d postgres
sleep 5
```

### "Failed to connect to NATS"
```bash
docker-compose up -d nats
sleep 3
```

### Workers not processing
Make sure:
1. NATS is running (`docker-compose ps`)
2. Queue API published messages first
3. Workers are actually started

## Clean Up

```bash
# Stop everything
docker-compose down

# Remove all data
docker-compose down -v
```

## Expected Timeline

- Setup: 2 minutes
- Simple tests: 2 minutes  
- Load tests: 8 minutes (4 min each)
- Total: ~15 minutes to see the full comparison!

---

**Next Step**: Check [README.md](README.md) for detailed documentation.
