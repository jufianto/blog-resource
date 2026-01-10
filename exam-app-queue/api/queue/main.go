package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jufianto/blog-resource/exam-app-queue/store"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

var (
	nc *nats.Conn
	js nats.JetStreamContext
)

func init() {
	viper.SetConfigFile("../../env/env.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper config not found: %v", err)
	}
}

func main() {
	// Setup NATS connection
	natsURL := viper.GetString("nats.url")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	js, err = nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}

	// Create or update stream
	streamName := "EXAM_SUBMISSIONS"
	streamSubject := "exam.submit"

	stream, err := js.StreamInfo(streamName)
	if err != nil {
		// Stream doesn't exist, create it
		log.Printf("Creating stream: %s", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubject},
			Storage:  nats.FileStorage,
			MaxAge:   24 * time.Hour, // Keep messages for 24 hours
		})
		if err != nil {
			log.Fatalf("Failed to create stream: %v", err)
		}
		log.Printf("Stream created successfully")
	} else {
		log.Printf("Stream already exists: %s (Messages: %d)", streamName, stream.State.Msgs)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/submit", handleSubmit)
	mux.HandleFunc("/health", handleHealth)

	port := viper.GetInt("api.queue_port")
	if port == 0 {
		port = 8081
	}

	log.Printf("Starting Queue-based API Server on :%d", port)
	log.Printf("Connected to NATS at %s", natsURL)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	var req store.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.UserID == "" || req.ExamID == "" || len(req.Answers) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create submission
	submission := store.ExamSubmission{
		ID:          uuid.New(),
		UserID:      req.UserID,
		ExamID:      req.ExamID,
		Answers:     req.Answers,
		SubmittedAt: time.Now(),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Publish to NATS (non-blocking, fast)
	data, err := json.Marshal(submission)
	if err != nil {
		http.Error(w, "Failed to encode submission", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = js.Publish("exam.submit", data, nats.Context(ctx))
	if err != nil {
		log.Printf("Failed to publish to NATS: %v", err)
		http.Error(w, "Failed to queue submission", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	log.Printf("Request queued in %v - User: %s, Exam: %s", duration, req.UserID, req.ExamID)

	// Return immediate response (fast)
	resp := store.SubmitResponse{
		Status:       "accepted",
		SubmissionID: submission.ID,
		Message:      "Exam submitted and being processed",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Check NATS connection
	status := "healthy"
	if nc.Status() != nats.CONNECTED {
		status = "unhealthy"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":      status,
		"mode":        "queue",
		"nats_status": nc.Status().String(),
	})
}
