package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jufianto/blog-resource/exam-app-queue/store"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

var (
	db          *store.Store
	workerCount int
	batchSize   int
)

func init() {
	viper.SetConfigFile("../env/env.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper config not found: %v", err)
	}

	workerCount = viper.GetInt("worker.count")
	if workerCount == 0 {
		workerCount = 10
	}

	batchSize = viper.GetInt("worker.batch_size")
	if batchSize == 0 {
		batchSize = 100
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup database connection
	sqlConnStr := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
	)

	configPool, err := pgxpool.ParseConfig(sqlConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// Connection pool settings - should accommodate all workers
	configPool.MaxConns = int32(workerCount + 5)
	configPool.MaxConnIdleTime = 30 * time.Second
	configPool.MinConns = int32(workerCount / 2)

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	db = store.NewStore(pool)

	// Setup NATS connection
	natsURL := viper.GetString("nats.url")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}

	log.Printf("Starting %d workers with batch size %d", workerCount, batchSize)
	log.Printf("Database pool: MaxConns=%d, MinConns=%d", configPool.MaxConns, configPool.MinConns)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, js, i, &wg)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	cancel() // Signal all workers to stop
	log.Println("Shutting down workers...")
	wg.Wait()
	log.Println("All workers stopped")
}

func worker(ctx context.Context, js nats.JetStreamContext, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Worker %d started", workerID)

	// ALL workers share the SAME consumer for work distribution
	sub, err := js.PullSubscribe("exam.submit", "exam-workers", nats.ManualAck())
	if err != nil {
		log.Fatalf("Worker %d failed to subscribe: %v", workerID, err)
	}
	defer sub.Unsubscribe()

	batch := make([]store.ExamSubmission, 0, batchSize)
	batchMsgs := make([]*nats.Msg, 0, batchSize) // Track messages for acking
	ticker := time.NewTicker(2 * time.Second)    // Flush batch every 2 seconds
	defer ticker.Stop()

	processed := 0

	for {
		select {
		case <-ctx.Done():
			// Flush remaining batch before stopping
			if len(batch) > 0 {
				if processBatch(ctx, batch, workerID) {
					for _, msg := range batchMsgs {
						msg.Ack()
					}
				} else {
					for _, msg := range batchMsgs {
						msg.Nak()
					}
				}
			}
			log.Printf("Worker %d stopped. Total processed: %d", workerID, processed)
			return

		case <-ticker.C:
			// Periodic flush
			if len(batch) > 0 {
				if processBatch(ctx, batch, workerID) {
					// Ack all messages
					for _, msg := range batchMsgs {
						msg.Ack()
					}
					processed += len(batch)
				} else {
					// Nak all messages to retry
					for _, msg := range batchMsgs {
						msg.Nak()
					}
				}
				batch = batch[:0]
				batchMsgs = batchMsgs[:0]
			}

			// After ticker, try to fetch messages
			msgs, err := sub.Fetch(batchSize, nats.MaxWait(100*time.Millisecond))
			if err != nil {
				if err != nats.ErrTimeout {
					log.Printf("Worker %d fetch error: %v", workerID, err)
				}
				continue
			}

			for _, msg := range msgs {
				// Check context before processing
				select {
				case <-ctx.Done():
					msg.Nak()
					return
				default:
				}

				var submission store.ExamSubmission
				if err := json.Unmarshal(msg.Data, &submission); err != nil {
					log.Printf("Worker %d failed to unmarshal: %v", workerID, err)
					msg.Nak()
					continue
				}

				// Update status to processing
				submission.Status = "processed"
				now := time.Now()
				submission.ProcessedAt = &now

				batch = append(batch, submission)
				batchMsgs = append(batchMsgs, msg)

				// If batch is full, process it
				if len(batch) >= batchSize {
					if processBatch(ctx, batch, workerID) {
						// Ack all messages in batch
						for _, m := range batchMsgs {
							m.Ack()
						}
						processed += len(batch)
					} else {
						// Nak all messages to retry
						for _, m := range batchMsgs {
							m.Nak()
						}
					}
					batch = batch[:0]
					batchMsgs = batchMsgs[:0]
					break
				}
			}
		}
	}
}

func processBatch(ctx context.Context, batch []store.ExamSubmission, workerID int) bool {
	if len(batch) == 0 {
		return true
	}

	start := time.Now()
	err := db.BatchInsertSubmissions(ctx, batch)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Worker %d failed to insert batch of %d: %v", workerID, len(batch), err)
		return false
	}

	log.Printf("Worker %d processed batch of %d in %v", workerID, len(batch), duration)
	return true
}
