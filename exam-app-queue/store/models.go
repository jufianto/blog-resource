package store

import (
	"time"

	"github.com/google/uuid"
)

type ExamSubmission struct {
	ID          uuid.UUID              `json:"id"`
	UserID      string                 `json:"user_id"`
	ExamID      string                 `json:"exam_id"`
	Answers     map[string]interface{} `json:"answers"`
	Score       *float64               `json:"score,omitempty"`
	SubmittedAt time.Time              `json:"submitted_at"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
}

type SubmitRequest struct {
	UserID  string                 `json:"user_id"`
	ExamID  string                 `json:"exam_id"`
	Answers map[string]interface{} `json:"answers"`
}

type SubmitResponse struct {
	Status       string    `json:"status"`
	SubmissionID uuid.UUID `json:"submission_id"`
	Message      string    `json:"message"`
}
