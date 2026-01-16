#!/bin/bash

# Monitor Docker container resources
# Usage: ./monitor-resources.sh direct|queue

MODE=$1
DURATION=${2:-300}  # Default 5 minutes
OUTPUT_FILE="load-test/results/${MODE}_resources.csv"

if [ -z "$MODE" ]; then
    echo "Usage: ./monitor-resources.sh <direct|queue> [duration_seconds]"
    exit 1
fi

echo "Monitoring resources for $MODE mode for $DURATION seconds..."
echo "Output: $OUTPUT_FILE"

# CSV Header
echo "timestamp,container,cpu_percent,mem_usage_mb,mem_limit_mb,mem_percent,net_input_mb,net_output_mb,block_read_mb,block_write_mb" > "$OUTPUT_FILE"

# Monitor for specified duration
END_TIME=$((SECONDS + DURATION))

while [ $SECONDS -lt $END_TIME ]; do
    TIMESTAMP=$(date +%s)
    
    # Get stats for postgres and nats
    docker stats --no-stream --format "{{.Container}},{{.CPUPerc}},{{.MemUsage}},{{.MemPerc}},{{.NetIO}},{{.BlockIO}}" exam-db exam-nats | while read line; do
        # Parse docker stats output
        CONTAINER=$(echo $line | cut -d',' -f1)
        CPU=$(echo $line | cut -d',' -f2 | sed 's/%//')
        MEM_USAGE=$(echo $line | cut -d',' -f3 | awk '{print $1}' | sed 's/MiB//')
        MEM_LIMIT=$(echo $line | cut -d',' -f3 | awk '{print $3}' | sed 's/MiB//')
        MEM_PERCENT=$(echo $line | cut -d',' -f4 | sed 's/%//')
        NET_IN=$(echo $line | cut -d',' -f5 | awk '{print $1}' | sed 's/MB//' | sed 's/kB//')
        NET_OUT=$(echo $line | cut -d',' -f5 | awk '{print $3}' | sed 's/MB//' | sed 's/kB//')
        BLOCK_READ=$(echo $line | cut -d',' -f6 | awk '{print $1}' | sed 's/MB//' | sed 's/GB//')
        BLOCK_WRITE=$(echo $line | cut -d',' -f6 | awk '{print $3}' | sed 's/MB//' | sed 's/GB//')
        
        echo "$TIMESTAMP,$CONTAINER,$CPU,$MEM_USAGE,$MEM_LIMIT,$MEM_PERCENT,$NET_IN,$NET_OUT,$BLOCK_READ,$BLOCK_WRITE" >> "$OUTPUT_FILE"
    done
    
    sleep 1
done

echo "Monitoring complete! Results saved to $OUTPUT_FILE"
