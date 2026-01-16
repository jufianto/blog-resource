#!/usr/bin/env python3
"""
Analyze and compare resource usage between Direct and Queue approaches
Generates charts and summary for blog post
"""

import pandas as pd
import matplotlib.pyplot as plt
import json
from pathlib import Path

def load_resource_data(mode):
    """Load resource monitoring CSV"""
    file_path = f"load-test/results/{mode}_resources.csv"
    try:
        df = pd.read_csv(file_path)
        df['datetime'] = pd.to_datetime(df['timestamp'], unit='s')
        return df
    except FileNotFoundError:
        print(f"❌ File not found: {file_path}")
        print(f"   Run: make test-{mode}-monitored")
        return None

def load_k6_results(mode):
    """Load k6 test results"""
    file_path = f"load-test/results/{mode}_summary.json"
    try:
        with open(file_path, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"⚠️  K6 results not found: {file_path}")
        return None

def analyze_container(df, container_name):
    """Analyze metrics for a specific container"""
    container_df = df[df['container'] == container_name].copy()
    
    if container_df.empty:
        return None
    
    return {
        'avg_cpu': container_df['cpu_percent'].mean(),
        'max_cpu': container_df['cpu_percent'].max(),
        'avg_mem_mb': container_df['mem_usage_mb'].mean(),
        'max_mem_mb': container_df['mem_usage_mb'].max(),
        'avg_mem_percent': container_df['mem_percent'].mean(),
        'max_mem_percent': container_df['mem_percent'].max(),
    }

def generate_comparison_charts(direct_df, queue_df):
    """Generate comparison charts for blog"""
    fig, axes = plt.subplots(2, 2, figsize=(15, 10))
    fig.suptitle('Direct vs Queue: Resource Usage Comparison', fontsize=16, fontweight='bold')
    
    # CPU Usage - PostgreSQL
    ax = axes[0, 0]
    direct_db = direct_df[direct_df['container'] == 'exam-db']
    queue_db = queue_df[queue_df['container'] == 'exam-db']
    
    ax.plot(direct_db['datetime'], direct_db['cpu_percent'], label='Direct', color='#e74c3c', linewidth=2)
    ax.plot(queue_db['datetime'], queue_db['cpu_percent'], label='Queue', color='#3498db', linewidth=2)
    ax.set_title('PostgreSQL CPU Usage (%)')
    ax.set_xlabel('Time')
    ax.set_ylabel('CPU %')
    ax.legend()
    ax.grid(True, alpha=0.3)
    
    # Memory Usage - PostgreSQL
    ax = axes[0, 1]
    ax.plot(direct_db['datetime'], direct_db['mem_usage_mb'], label='Direct', color='#e74c3c', linewidth=2)
    ax.plot(queue_db['datetime'], queue_db['mem_usage_mb'], label='Queue', color='#3498db', linewidth=2)
    ax.set_title('PostgreSQL Memory Usage (MB)')
    ax.set_xlabel('Time')
    ax.set_ylabel('Memory (MB)')
    ax.legend()
    ax.grid(True, alpha=0.3)
    
    # CPU Usage - NATS
    ax = axes[1, 0]
    queue_nats = queue_df[queue_df['container'] == 'exam-nats']
    
    if not queue_nats.empty:
        ax.plot(queue_nats['datetime'], queue_nats['cpu_percent'], label='Queue (NATS)', color='#2ecc71', linewidth=2)
        ax.axhline(y=0, color='#95a5a6', linestyle='--', label='Direct (No NATS)')
        ax.set_title('NATS CPU Usage (%)')
        ax.set_xlabel('Time')
        ax.set_ylabel('CPU %')
        ax.legend()
        ax.grid(True, alpha=0.3)
    
    # Summary Bar Chart
    ax = axes[1, 1]
    
    direct_db_stats = analyze_container(direct_df, 'exam-db')
    queue_db_stats = analyze_container(queue_df, 'exam-db')
    
    categories = ['Avg CPU %', 'Max CPU %', 'Avg Mem %', 'Max Mem %']
    direct_values = [
        direct_db_stats['avg_cpu'],
        direct_db_stats['max_cpu'],
        direct_db_stats['avg_mem_percent'],
        direct_db_stats['max_mem_percent']
    ]
    queue_values = [
        queue_db_stats['avg_cpu'],
        queue_db_stats['max_cpu'],
        queue_db_stats['avg_mem_percent'],
        queue_db_stats['max_mem_percent']
    ]
    
    x = range(len(categories))
    width = 0.35
    
    ax.bar([i - width/2 for i in x], direct_values, width, label='Direct', color='#e74c3c', alpha=0.8)
    ax.bar([i + width/2 for i in x], queue_values, width, label='Queue', color='#3498db', alpha=0.8)
    
    ax.set_title('PostgreSQL: Summary Comparison')
    ax.set_xticks(x)
    ax.set_xticklabels(categories, rotation=15)
    ax.legend()
    ax.grid(True, alpha=0.3, axis='y')
    
    plt.tight_layout()
    plt.savefig('load-test/results/comparison_chart.png', dpi=300, bbox_inches='tight')
    print("📊 Chart saved: load-test/results/comparison_chart.png")

def generate_markdown_report(direct_df, queue_df, direct_k6, queue_k6):
    """Generate markdown summary for blog"""
    
    direct_db = analyze_container(direct_df, 'exam-db')
    queue_db = analyze_container(queue_df, 'exam-db')
    queue_nats = analyze_container(queue_df, 'exam-nats')
    
    report = f"""# Load Test Results: Direct vs Queue Approach

## Test Configuration
- **Duration**: 5 minutes
- **Peak Load**: 3000 concurrent users
- **Total Requests**: ~413,000

## 📊 Resource Usage Comparison

### PostgreSQL (exam-db)

| Metric | Direct Approach | Queue Approach | Difference |
|--------|----------------|----------------|------------|
| **Avg CPU** | {direct_db['avg_cpu']:.2f}% | {queue_db['avg_cpu']:.2f}% | {queue_db['avg_cpu'] - direct_db['avg_cpu']:+.2f}% |
| **Max CPU** | {direct_db['max_cpu']:.2f}% | {queue_db['max_cpu']:.2f}% | {queue_db['max_cpu'] - direct_db['max_cpu']:+.2f}% |
| **Avg Memory** | {direct_db['avg_mem_mb']:.0f} MB | {queue_db['avg_mem_mb']:.0f} MB | {queue_db['avg_mem_mb'] - direct_db['avg_mem_mb']:+.0f} MB |
| **Max Memory** | {direct_db['max_mem_mb']:.0f} MB | {queue_db['max_mem_mb']:.0f} MB | {queue_db['max_mem_mb'] - direct_db['max_mem_mb']:+.0f} MB |
| **Avg Mem %** | {direct_db['avg_mem_percent']:.2f}% | {queue_db['avg_mem_percent']:.2f}% | {queue_db['avg_mem_percent'] - direct_db['avg_mem_percent']:+.2f}% |

### NATS (exam-nats) - Queue Only

| Metric | Value |
|--------|-------|
| **Avg CPU** | {queue_nats['avg_cpu']:.2f}% |
| **Max CPU** | {queue_nats['max_cpu']:.2f}% |
| **Avg Memory** | {queue_nats['avg_mem_mb']:.0f} MB |
| **Max Memory** | {queue_nats['max_mem_mb']:.0f} MB |

"""
    
    if direct_k6 and queue_k6:
        direct_http = direct_k6['metrics']['http_req_duration']['values']
        queue_http = queue_k6['metrics']['http_req_duration']['values']
        
        report += f"""## ⚡ Performance Metrics

### Response Times

| Metric | Direct | Queue | Difference |
|--------|--------|-------|------------|
| **Average** | {direct_http['avg']:.2f} ms | {queue_http['avg']:.2f} ms | {queue_http['avg'] - direct_http['avg']:+.2f} ms |
| **P95** | {direct_http['p(95)']:.2f} ms | {queue_http['p(95)']:.2f} ms | {queue_http['p(95)'] - direct_http['p(95)']:+.2f} ms |
| **P99** | {direct_http['p(99)']:.2f} ms | {queue_http['p(99)']:.2f} ms | {queue_http['p(99)'] - direct_http['p(99)']:+.2f} ms |
| **Max** | {direct_http['max']:.2f} ms | {queue_http['max']:.2f} ms | {queue_http['max'] - direct_http['max']:+.2f} ms |

### Throughput

| Metric | Direct | Queue |
|--------|--------|-------|
| **Requests/sec** | {direct_k6['metrics']['http_reqs']['values']['rate']:.2f} | {queue_k6['metrics']['http_reqs']['values']['rate']:.2f} |
| **Total Requests** | {direct_k6['metrics']['http_reqs']['values']['count']} | {queue_k6['metrics']['http_reqs']['values']['count']} |

"""
    
    report += """## 🎯 Key Insights

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
"""
    
    with open('load-test/results/COMPARISON.md', 'w') as f:
        f.write(report)
    
    print("📝 Report saved: load-test/results/COMPARISON.md")

def main():
    print("🔍 Analyzing load test results...\n")
    
    # Load data
    direct_df = load_resource_data('direct')
    queue_df = load_resource_data('queue')
    
    if direct_df is None or queue_df is None:
        print("\n❌ Missing resource data. Run monitored tests first:")
        print("   make clean-db && make test-direct-monitored")
        print("   make clean-db && make test-queue-monitored")
        return
    
    # Load K6 results
    direct_k6 = load_k6_results('direct')
    queue_k6 = load_k6_results('queue')
    
    # Generate analysis
    print("\n📊 Generating charts...")
    generate_comparison_charts(direct_df, queue_df)
    
    print("\n📝 Generating markdown report...")
    generate_markdown_report(direct_df, queue_df, direct_k6, queue_k6)
    
    print("\n✅ Analysis complete!")
    print("\nResults:")
    print("  - load-test/results/comparison_chart.png")
    print("  - load-test/results/COMPARISON.md")
    print("\nOpen COMPARISON.md to see the full report!")

if __name__ == '__main__':
    main()
