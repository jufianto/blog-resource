# 🎉 Exam App Queue - Project Complete!

## ✅ What's Been Built

A complete, production-ready demonstration comparing two architectural approaches for handling high-concurrency database writes.

## 📦 Project Contents

### Core Application
- ✅ **Direct Insert API** (`api/direct/`) - Port 8080
- ✅ **Queue-based API** (`api/queue/`) - Port 8081  
- ✅ **Background Workers** (`worker/`) - 10 concurrent workers with batch processing
- ✅ **Database Layer** (`store/`) - PostgreSQL operations with connection pooling
- ✅ **Migrations** (`migrate/`) - Database schema management

### Infrastructure
- ✅ **Docker Compose** - PostgreSQL + NATS with resource limits:
  - PostgreSQL: 1 CPU, 512MB RAM (simulates real server)
  - NATS: 0.5 CPU, 256MB RAM
- ✅ **Database Schema** - Exam submissions table with indexes
- ✅ **Configuration** - Viper-based config management

### Testing & Monitoring
- ✅ **Load Tests** (`load-test/`) - k6 scripts for both approaches
- ✅ **Metrics & Logging** - Built-in performance monitoring
- ✅ **NATS Monitoring** - Web UI at http://localhost:8222

### Documentation
- ✅ **README.md** - Complete documentation
- ✅ **QUICKSTART.md** - 15-minute quick start guide
- ✅ **COMPARISON.md** - Detailed performance analysis
- ✅ **Makefile** - Convenient command shortcuts

## 🚀 Quick Commands

```bash
# Setup (one time)
make setup              # Start infra + run migrations

# Run Direct Approach
make run-direct         # Terminal 1

# Run Queue Approach  
make run-queue          # Terminal 1
make run-worker         # Terminal 2

# Load Testing
make test-direct        # Test direct approach
make test-queue         # Test queue approach

# Utilities
make check-db          # Check submission counts
make logs              # View container logs
make clean             # Clean up everything
```

## 📊 Expected Results

### Direct Insert (Will Struggle)
- Response time: 200-800ms
- Success rate: 60-80%
- Errors under load: timeout, connection exhaustion

### Queue-based (Will Excel)  
- Response time: < 10ms
- Success rate: 99%+
- Handles 3000+ req/s smoothly

## 🎓 Learning Objectives

This project demonstrates:

1. **pgxpool** - Connection pooling importance
2. **NATS JetStream** - Message queue for high throughput
3. **Worker Pools** - Concurrent batch processing
4. **Load Testing** - k6 for performance validation
5. **Resource Constraints** - Simulating real production limits
6. **Architecture Trade-offs** - When to use queues vs direct insert

## 📁 File Structure

```
exam-app-queue/
├── README.md              # Main documentation
├── QUICKSTART.md          # Fast setup guide
├── COMPARISON.md          # Performance analysis
├── Makefile               # Helper commands
├── docker-compose.yaml    # Infrastructure (limited resources)
├── go.mod                 # Dependencies
│
├── api/
│   ├── direct/main.go    # Direct insert (port 8080)
│   └── queue/main.go     # Queue-based (port 8081)
│
├── worker/main.go         # Background consumers
│
├── store/
│   ├── models.go         # Data structures
│   └── db.go             # Database operations
│
├── migrate/
│   ├── migrate.go
│   └── migrations/
│       ├── 001_create_exam_submissions.up.sql
│       └── 001_create_exam_submissions.down.sql
│
├── load-test/
│   ├── test_direct.js    # k6 test for direct
│   ├── test_queue.js     # k6 test for queue
│   └── results/          # Test outputs
│
└── env/
    ├── env.yaml          # Your config
    └── sample.env.yaml   # Template
```

## 🔧 Technologies Used

- **Go 1.23.6** - Application language
- **PostgreSQL 16** - Database (resource-limited)
- **NATS JetStream** - Message queue
- **pgx v5** - PostgreSQL driver with pooling
- **k6** - Load testing tool
- **Docker Compose** - Container orchestration
- **Viper** - Configuration management

## 💡 Blog Post Ideas

This project is perfect for writing about:

1. **"When Database Connection Pools Fail"** - Why pooling isn't enough
2. **"Queue vs Direct Insert: A Performance Showdown"** - This comparison
3. **"Handling 10,000 Concurrent Users"** - Architecture patterns
4. **"Message Queues Explained"** - NATS JetStream tutorial
5. **"Load Testing with k6"** - Performance validation
6. **"Batch Processing for Performance"** - Worker pool patterns
7. **"Backpressure in Distributed Systems"** - How queues help

## 🎯 Next Steps

1. **Start infrastructure**: `make setup`
2. **Read**: [QUICKSTART.md](QUICKSTART.md)
3. **Run both approaches**: Compare performance
4. **Load test**: See the difference under stress
5. **Analyze**: Check [COMPARISON.md](COMPARISON.md)
6. **Write your blog post!** 📝

## 🐛 Troubleshooting

See [README.md](README.md) "Troubleshooting" section.

Common issues:
- Connection refused → `make start-infra`
- Workers not processing → Check NATS is running
- Load test fails → Ensure correct API is running

## 📸 Screenshots for Blog

Capture these for your blog post:
1. k6 output showing response times (direct vs queue)
2. NATS monitoring UI (http://localhost:8222)
3. Database query showing submission counts
4. Side-by-side terminal showing both APIs
5. Load test graphs (direct failing, queue succeeding)

## 🙏 Credits

Built for demonstrating queue-based architecture patterns in high-concurrency scenarios.

## 📄 License

Educational/blog resource - use freely!

---

**Ready to demonstrate the power of message queues! 🚀**

Run `make help` to see all available commands.
