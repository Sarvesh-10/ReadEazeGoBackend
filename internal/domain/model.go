package domain

import "time"

type User struct {
	ID       int
	Name     string
	Password string
	Email    string
}

type Book struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	PDFData    []byte    `json:"-"` // Exclude from JSON response
	UploadedAt time.Time `json:"uploaded_at"`
	CoverImage []byte    `json:"cover_image"` // Base64 encoded image data
}

type BookMetaData struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CoverImage string `json:"cover_image"` // base64 string
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
)

type BookIndexingJob struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	UserID    int       `json:"user_id"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SSEMessage struct {
	BookID  int       `json:"book_id"`
	Status  JobStatus `json:"status"`
	Name    string    `json:"name"`
	Message string    `json:"message,omitempty"`
}
