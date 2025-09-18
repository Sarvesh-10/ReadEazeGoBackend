package models

type SSEMessage struct {
	BookID  int       `json:"book_id"`
	Status  JobStatus `json:"status"`
	Name    string    `json:"name"`
	Message string    `json:"message,omitempty"`
}
