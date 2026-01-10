package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jufianto/blog-resource/exam-app-queue/store"
	"github.com/spf13/viper"
)

var (
	db *store.Store
)

func init() {
	viper.SetConfigFile("../../env/env.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper config not found: %v", err)
	}
}

func main() {
	// Setup database connection
	sqlConnStr := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
	)

	ctx := context.Background()

	configPool, err := pgxpool.ParseConfig(sqlConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// Connection pool settings
	configPool.MaxConns = 50
	configPool.MaxConnIdleTime = 30 * time.Second
	configPool.MinConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	db = store.NewStore(pool)

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/submit", handleSubmit)
	mux.HandleFunc("/health", handleHealth)

	port := viper.GetInt("api.direct_port")
	if port == 0 {
		port = 8080
	}

	log.Printf("Starting Direct API Server on :%d", port)
	log.Printf("Database pool: MaxConns=%d, MinConns=%d", configPool.MaxConns, configPool.MinConns)

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
		Status:      "completed",
		CreatedAt:   time.Now(),
		ProcessedAt: func() *time.Time { t := time.Now(); return &t }(),
	}

	// Direct insert to database (blocking)
	ctx := r.Context()
	if err := db.InsertSubmission(ctx, submission); err != nil {
		log.Printf("Failed to insert submission: %v", err)
		http.Error(w, "Failed to save submission", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	log.Printf("Request processed in %v - User: %s, Exam: %s", duration, req.UserID, req.ExamID)

	// Return success response
	resp := store.SubmitResponse{
		Status:       "success",
		SubmissionID: submission.ID,
		Message:      "Exam submitted and saved",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"mode":   "direct",
	})
}
