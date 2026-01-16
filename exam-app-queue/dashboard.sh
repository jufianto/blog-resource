#!/bin/bash

# Real-time resource monitoring dashboard
# Usage: ./dashboard.sh

watch -n 1 '
echo "=== RESOURCE USAGE DASHBOARD ==="
echo ""
echo "Container Stats:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}" exam-db exam-nats

echo ""
echo "PostgreSQL Connections:"
docker exec exam-db psql -U postgres -d exam_app -t -c "SELECT count(*) FROM pg_stat_activity WHERE datname = '"'exam_app'"';" 2>/dev/null || echo "N/A"

echo ""
echo "Database Stats:"
docker exec exam-db psql -U postgres -d exam_app -t -c "SELECT status, COUNT(*) FROM exam_submissions GROUP BY status;" 2>/dev/null || echo "N/A"

echo ""
echo "NATS Messages (if queue mode):"
curl -s http://localhost:8222/jsz 2>/dev/null | grep -o '"'"'messages":[0-9]*' | head -1 || echo "N/A"

echo ""
echo "Press Ctrl+C to stop monitoring"
'
