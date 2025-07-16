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
