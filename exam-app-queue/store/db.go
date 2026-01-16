package store

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreInterface interface {
	InsertSubmission(ctx context.Context, submission ExamSubmission) error
	BatchInsertSubmissions(ctx context.Context, submissions []ExamSubmission) error
	GetSubmission(ctx context.Context, id uuid.UUID) (*ExamSubmission, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) InsertSubmission(ctx context.Context, submission ExamSubmission) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Printf("error acquiring connection: %v", err)
		return err
	}
	defer conn.Release()

	answersJSON, err := json.Marshal(submission.Answers)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO exam_submissions 
		(id, user_id, exam_id, answers, score, submitted_at, processed_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`

	_, err = conn.Exec(
		ctx, query,
		submission.ID,
		submission.UserID,
		submission.ExamID,
		answersJSON,
		submission.Score,
		submission.SubmittedAt,
		submission.ProcessedAt,
		submission.Status,
		submission.CreatedAt,
	)

	if err != nil {
		log.Printf("error inserting submission: %v", err)
		return err
	}

	return nil
}

func (s *Store) BatchInsertSubmissions(ctx context.Context, submissions []ExamSubmission) error {
	if len(submissions) == 0 {
		return nil
	}

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Printf("error acquiring connection: %v", err)
		return err
	}
	defer conn.Release()

	batch := &pgx.Batch{}
	query := `
		INSERT INTO exam_submissions 
		(id, user_id, exam_id, answers, score, submitted_at, processed_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`

	for _, submission := range submissions {
		answersJSON, err := json.Marshal(submission.Answers)
		if err != nil {
			log.Printf("error marshaling answers: %v", err)
			continue
		}

		batch.Queue(
			query,
			submission.ID,
			submission.UserID,
			submission.ExamID,
			answersJSON,
			submission.Score,
			submission.SubmittedAt,
			submission.ProcessedAt,
			submission.Status,
			submission.CreatedAt,
		)
	}

	br := conn.SendBatch(ctx, batch)
	defer br.Close()

	inserted := 0
	skipped := 0

	for i := 0; i < len(submissions); i++ {
		cmdTag, err := br.Exec()
		if err != nil {
			log.Printf("error in batch insert at index %d: %v", i, err)
			continue
		}

		// Check if row was actually inserted
		if cmdTag.RowsAffected() > 0 {
			inserted++
		} else {
			skipped++
		}
	}

	if skipped > 0 {
		log.Printf("Batch insert: %d inserted, %d skipped (duplicates)", inserted, skipped)
	}

	return nil
}

func (s *Store) GetSubmission(ctx context.Context, id uuid.UUID) (*ExamSubmission, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
		SELECT id, user_id, exam_id, answers, score, submitted_at, processed_at, status, created_at
		FROM exam_submissions
		WHERE id = $1
	`

	var submission ExamSubmission
	var answersJSON []byte

	err = conn.QueryRow(ctx, query, id).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ExamID,
		&answersJSON,
		&submission.Score,
		&submission.SubmittedAt,
		&submission.ProcessedAt,
		&submission.Status,
		&submission.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(answersJSON, &submission.Answers)
	if err != nil {
		return nil, err
	}

	return &submission, nil
}

func (s *Store) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `
		UPDATE exam_submissions
		SET status = $1, processed_at = $2
		WHERE id = $3
	`

	_, err = conn.Exec(ctx, query, status, time.Now(), id)
	return err
}
