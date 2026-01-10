CREATE TABLE IF NOT EXISTS exam_submissions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    exam_id VARCHAR(50) NOT NULL,
    answers JSONB NOT NULL,
    score DECIMAL(5, 2),
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_user_exam ON exam_submissions(user_id, exam_id);
CREATE INDEX IF NOT EXISTS idx_status ON exam_submissions(status);
CREATE INDEX IF NOT EXISTS idx_submitted_at ON exam_submissions(submitted_at);
CREATE INDEX IF NOT EXISTS idx_exam_id ON exam_submissions(exam_id);
