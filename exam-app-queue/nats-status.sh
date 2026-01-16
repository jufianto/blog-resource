#!/bin/bash
# Simple NATS monitoring script

echo "=== NATS JetStream Status ==="
echo ""

# Get basic stats
STATS=$(curl -s http://localhost:8222/jsz)

MESSAGES=$(echo $STATS | jq -r '.messages')
BYTES=$(echo $STATS | jq -r '.bytes')
STREAMS=$(echo $STATS | jq -r '.streams')
CONSUMERS=$(echo $STATS | jq -r '.consumers')

echo "📊 Overview:"
echo "  Streams:   $STREAMS"
echo "  Consumers: $CONSUMERS"
echo "  Messages:  $(printf "%'d" $MESSAGES)"
echo "  Size:      $(numfmt --to=iec-i --suffix=B $BYTES 2>/dev/null || echo "$BYTES bytes")"
echo ""

# Get stream details
echo "📦 Stream Details:"
curl -s "http://localhost:8222/jsz?streams=1" | jq -r '.streams[] | "  Name: \(.name)\n  Messages: \(.state.messages)\n  Bytes: \(.state.bytes)\n  First Seq: \(.state.first_seq)\n  Last Seq: \(.state.last_seq)\n  Consumers: \(.state.consumer_count)"'
echo ""

# Try to get consumer info
echo "👥 Consumer Info:"
CONSUMER_DATA=$(curl -s "http://localhost:8222/jsz?consumers=1")

# Check if consumer_detail exists
if echo "$CONSUMER_DATA" | jq -e '.streams[0].consumer_detail' > /dev/null 2>&1; then
    echo "$CONSUMER_DATA" | jq -r '.streams[0].consumer_detail[] | "  Name: \(.name)\n  Pending: \(.num_pending)\n  Ack Pending: \(.num_ack_pending)\n  Delivered: \(.delivered.consumer_seq)"'
else
    echo "  No consumer details available"
    echo "  (Consumers exist but details not in response)"
fi

echo ""
echo "🔗 Web Interfaces:"
echo "  NATS Monitor: http://localhost:8222"
echo "  Surveyor UI:  http://localhost:7777 (if running)"
