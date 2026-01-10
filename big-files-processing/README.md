# Big Files Processing - CSV to PostgreSQL Importer

A high-performance Go application that demonstrates two approaches for processing large CSV files and importing them into PostgreSQL: a single-threaded method and a concurrent worker-pool method.

## Overview

This project reads large CSV files (specifically sales transaction data) and inserts the data into a PostgreSQL database. It provides a comparison between:
- **Non-concurrent method**: Sequential reading and inserting
- **Concurrent method**: Multi-worker parallel processing with goroutines

The application is designed to handle millions of records efficiently, with the example using a 5 million row sales dataset.

## Features

- ✅ Two processing modes: single-threaded vs concurrent
- ✅ Configurable worker pool size (default: 30 workers)
- ✅ PostgreSQL connection pooling with pgxpool
- ✅ Database migration support
- ✅ Progress logging during processing
- ✅ Configurable record limits for testing
- ✅ Performance timing and metrics

## Project Structure

```
big-files-processing/
├── main.go                     # Main application entry point
├── single_read.go              # Non-concurrent CSV processing
├── concurrent_read_process.go  # Concurrent CSV processing with workers
├── docker-compose.yaml         # PostgreSQL container setup
├── env/
│   ├── env.yaml               # Configuration file (you create this)
│   └── sample.env.yaml        # Configuration template
├── migrate/
│   ├── migrate.go             # Database migration tool
│   └── migrations/
│       ├── 000001_create_transaction_table.up.sql
│       └── 000001_create_transaction_table.down.sql
└── store/
    └── db.go                  # Database operations and interface
```

## Prerequisites

- Go 1.23.6 or higher
- Docker and Docker Compose (for PostgreSQL)
- A CSV file with sales data (expected at `../resource/sales_5000000.csv`)

## Installation & Setup

### 1. Clone and Navigate

```bash
cd big-files-processing
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start PostgreSQL Database

```bash
docker-compose up -d
```

This will start a PostgreSQL 16 container:
- **Container name**: bgp-transactions
- **Port**: 5411 (host) → 5432 (container)
- **Database**: bgp-transactions
- **User**: postgres
- **Password**: password
- **Memory limit**: 512MB
- **CPU limit**: 1 core

### 4. Configure Database Connection

Create your configuration file:

```bash
cp env/sample.env.yaml env/env.yaml
```

Edit `env/env.yaml`:

```yaml
database:
  host: localhost
  port: 5411
  user: postgres
  password: password
  dbname: bgp-transactions
```

### 5. Run Database Migrations

Navigate to the migrate directory and run migrations:

```bash
cd migrate
go run migrate.go up
```

**Migration commands:**
- `go run migrate.go up` - Apply migrations
- `go run migrate.go down` - Rollback migrations (requires confirmation)
- `go run migrate.go force VERSION` - Force migration to a specific version

### 6. Prepare Your CSV File

Ensure your CSV file is located at `../resource/sales_5000000.csv` relative to the main.go file.

**Expected CSV format:**
```
Region,Country,Item Type,Sales Channel,Order Priority,Order Date,Order ID,Ship Date,Units Sold,Unit Price,Unit Cost,Total Revenue,Total Cost,Total Profit
Sub-Saharan Africa,South Africa,Fruits,Offline,M,2014-07-07,443368995,2014-07-27,1593,9.33,6.92,14862.69,11023.56,3839.13
```

## Usage

### Running the Application

From the `big-files-processing` directory:

```bash
go run .
```

### Configuration Options

Edit [main.go](main.go) to customize:

```go
var (
    worker = 30        // Number of concurrent workers (for concurrent mode)
    sep    = 1000      // Progress logging interval
)

func main() {
    // ...
    
    limit := 400000           // Max records to process (-1 for unlimited)
    concurrentWork := true    // true = concurrent, false = single-threaded
    
    // ...
}
```

### Connection Pool Configuration

The application uses PostgreSQL connection pooling:

```go
configPool.MaxConns = 60              // Maximum connections
configPool.MaxConnIdleTime = 20 * time.Second
configPool.MinConns = 10              // Minimum idle connections
```

## How It Works

### Non-Concurrent Method ([single_read.go](single_read.go))

1. Opens CSV file and creates a reader
2. Skips the header row
3. Reads records one by one sequentially
4. Parses each record into a `DataSales` struct
5. Inserts into database immediately
6. Logs progress every 1000 records

**Pros:**
- Simple and straightforward
- Predictable memory usage
- Easy to debug

**Cons:**
- Slower for large datasets
- Single-threaded bottleneck

### Concurrent Method ([concurrent_read_process.go](concurrent_read_process.go))

1. Creates a buffered channel for CSV records (capacity: 100,000)
2. Spawns a reader goroutine that:
   - Reads CSV records
   - Sends them to the channel
   - Logs progress every 50,000 records
3. Spawns 30 worker goroutines that:
   - Consume records from the channel
   - Parse and insert into database
   - Track individual worker progress
4. Waits for all workers to complete
5. Aggregates total inserted records

**Pros:**
- Much faster for large datasets
- Utilizes multiple CPU cores
- Efficient parallel processing

**Cons:**
- More complex code
- Higher memory usage
- Requires careful synchronization

### Data Flow (Concurrent Mode)

```
CSV File 
   ↓
readCSVRAW (goroutine)
   ↓
recordChan (buffered channel: 100k capacity)
   ↓
30 × parseCSV workers (goroutines)
   ↓
PostgreSQL Database
   ↓
totalChan (aggregates results)
```

## Database Schema

The application creates a `transactions` table:

```sql
CREATE TABLE IF NOT EXISTS transactions (
    id varchar(36) PRIMARY KEY,
    region VARCHAR(100),
    country VARCHAR(100),
    item_type VARCHAR(100),
    sales_channel VARCHAR(50),
    order_priority CHAR(1),
    order_date DATE,
    order_id BIGINT,
    ship_date DATE,
    units_sold INT,
    unit_price DECIMAL(10, 2),
    unit_cost DECIMAL(10, 2),
    total_revenue DECIMAL(15, 2),
    total_cost DECIMAL(15, 2),
    total_profit DECIMAL(15, 2)
);
```

## Performance Tuning

### 1. Worker Count
Adjust based on your system:
```go
var worker = 30  // Increase for more CPU cores, decrease for limited resources
```

### 2. Channel Buffer Size
In [concurrent_read_process.go](concurrent_read_process.go#L16):
```go
recordChan := make(chan []string, 100000)  // Larger = more memory, better throughput
```

### 3. Database Connection Pool
```go
configPool.MaxConns = 60  // Should be >= worker count
```

### 4. Record Limit
For testing:
```go
limit := 400000  // Process only first 400k records
limit := -1      // Process entire file
```

### 5. Progress Logging
Reduce logging overhead:
```go
var sep = 10000  // Log every 10k instead of 1k
```

## Monitoring

The application provides runtime information:

```
runtime 8                                    # Number of CPU cores
starting with concurrent method              # Selected mode
starting worker 0                            # Worker initialization
starting worker 1
...
50000 has read from csv                      # CSV reading progress
1000 record has been saved on worker 3       # Worker progress
total inserted 400000                        # Final count
end in 45.23 seconds                         # Total execution time
```

## Troubleshooting

### Connection Refused
```
failed to fetch pool connection: unable to connect to database
```
**Solution**: Ensure PostgreSQL is running:
```bash
docker-compose ps
docker-compose up -d
```

### File Not Found
```
failed to open: no such file or directory
```
**Solution**: Verify CSV file path in [main.go](main.go#L37):
```go
file, err := os.Open("../resource/sales_5000000.csv")
```

### Migration Failed
```
Failed to get migration version
```
**Solution**: 
1. Check database is running
2. Verify `env/env.yaml` configuration
3. Run from `migrate/` directory

### Out of Memory
**Solution**: Reduce worker count or channel buffer size:
```go
var worker = 10  // Reduce workers
recordChan := make(chan []string, 10000)  // Smaller buffer
```

## Dependencies

```
github.com/golang-migrate/migrate/v4  # Database migrations
github.com/google/uuid                # UUID generation
github.com/jackc/pgx/v5              # PostgreSQL driver
github.com/spf13/viper               # Configuration management
```

## Performance Comparison

Approximate performance (5 million records):
- **Non-concurrent**: ~300-400 seconds
- **Concurrent (30 workers)**: ~45-90 seconds

*Results vary based on hardware, database configuration, and disk I/O.*

## Advanced Usage

### Custom CSV Format

Modify [main.go](main.go#L130-L188) `getResult()` function to match your CSV structure:

```go
func getResult(record []string) (DataSales, error) {
    // Adjust field indices and parsing logic
    // based on your CSV column order
}
```

### Different Database

Implement the `StoreInterface` in [store/db.go](store/db.go):

```go
type StoreInterface interface {
    InsertSales(ctx context.Context, sales DataSales) error
}
```

### Batch Inserts

For even better performance, modify `InsertSales` to use batch inserts with `pgx.Batch`.

## License

This is a demo/educational project. Use as needed for your own purposes.

## Contributing

This is a blog resource repository. Feel free to fork and adapt for your needs.

## Contact

For questions or issues, please refer to the blog post associated with this code.
