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

	// Connection pool settings
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

	log.Println("Shutting down workers...")
	cancel()
	wg.Wait()
	log.Println("All workers stopped")
}

func worker(ctx context.Context, js nats.JetStreamContext, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Worker %d started", workerID)

	// ALL workers share the SAME consumer with increased ack pending limit
	sub, err := js.PullSubscribe("exam.submit", "exam-workers",
		nats.ManualAck(),
		nats.MaxAckPending(20000), // Allow up to 20,000 unacked messages
	)
	if err != nil {
		log.Fatalf("Worker %d failed to subscribe: %v", workerID, err)
	}
	defer sub.Unsubscribe()

	processed := 0

	for {
		// Check if we should stop
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped. Total processed: %d", workerID, processed)
			return
		default:
		}

		// Fetch a batch of messages (blocks up to 500ms)
		msgs, err := sub.Fetch(batchSize, nats.MaxWait(1000*time.Millisecond))
		if err != nil {
			if err == nats.ErrTimeout {
				// No messages, wait a bit and try again
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Printf("Worker %d fetch error: %v", workerID, err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Worker %d fetched %d messages", workerID, len(msgs))

		// Build batch of submissions
		batch := make([]store.ExamSubmission, 0, len(msgs))
		validMsgs := make([]*nats.Msg, 0, len(msgs))

		for _, msg := range msgs {
			var submission store.ExamSubmission
			if err := json.Unmarshal(msg.Data, &submission); err != nil {
				log.Printf("Worker %d failed to unmarshal: %v", workerID, err)
				msg.Nak()
				continue
			}

			// Update status
			submission.Status = "processed"
			now := time.Now()
			submission.ProcessedAt = &now

			batch = append(batch, submission)
			validMsgs = append(validMsgs, msg)
		}

		// Process the batch
		if len(batch) > 0 {
			start := time.Now()
			err := db.BatchInsertSubmissions(ctx, batch)
			duration := time.Since(start)

			if err != nil {
				log.Printf("Worker %d failed to insert batch of %d: %v", workerID, len(batch), err)
				// NAK all messages for retry
				for _, msg := range validMsgs {
					msg.Nak()
				}
			} else {
				log.Printf("Worker %d processed batch of %d in %v", workerID, len(batch), duration)
				// ACK all messages - SUCCESS!
				ackErrors := 0
				for i, msg := range validMsgs {
					if err := msg.Ack(); err != nil {
						log.Printf("Worker %d ACK error for message %d: %v", workerID, i, err)
						ackErrors++
					}
				}
				if ackErrors == 0 {
					log.Printf("Worker %d successfully ACK'd %d messages", workerID, len(validMsgs))
				} else {
					log.Printf("Worker %d had %d ACK errors out of %d messages", workerID, ackErrors, len(validMsgs))
				}
				processed += len(batch)
			}
		}
	}
}
