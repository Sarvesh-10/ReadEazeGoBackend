package domain

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"

	models "github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	redis "github.com/redis/go-redis/v9"
)

type BookRepository interface {
	SaveBook(book Book) (int, error)
	GetBookByID(bookID int, userID int) (*Book, error)
	GetBooksMetaByUser(UserId int) ([]BookMetaData, error)
	SaveUserBookProfile(profile models.UserBookProfile) error
	DeleteBook(bookID int, userID int) error
}
type BookRepositoryImpl struct {
	DB    *sql.DB
	cache *redis.Client
}

func (r *BookRepositoryImpl) SaveBookIndexingJob(job models.BookIndexingJob) (int, error) {
	var jobID int
	err := r.DB.QueryRow(
		`INSERT INTO book_indexing_jobs (book_id, user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		job.BookID, job.UserID, string(job.Status), job.CreatedAt, job.UpdatedAt,
	).Scan(&jobID)
	if err != nil {
		return 0, err
	}
	return jobID, nil
}
func NewBookRepository(db *sql.DB, cache *redis.Client) *BookRepositoryImpl {
	return &BookRepositoryImpl{DB: db, cache: cache}
}

func (r *BookRepositoryImpl) SaveBook(book Book) (int, error) {
	var bookID int
	err := r.DB.QueryRow(
		`INSERT INTO books (user_id, name, file_data, uploaded_at, cover_image)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		book.UserID, book.Name, book.PDFData, book.UploadedAt, book.CoverImage,
	).Scan(&bookID)

	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func (r *BookRepositoryImpl) GetBookByID(bookID int, userID int) (*Book, error) {
	book := &Book{}
	err := r.DB.QueryRow(
		"SELECT id, user_id, name, file_data, uploaded_at FROM books WHERE id=$1 AND user_id=$2",
		bookID, userID,
	).Scan(&book.ID, &book.UserID, &book.Name, &book.PDFData, &book.UploadedAt)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepositoryImpl) GetBooksMetaByUser(UserId int) ([]BookMetaData, error) {
	var books []BookMetaData
	rows, err := r.DB.Query("SELECT id, name, cover_image FROM books WHERE user_id = $1", UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book BookMetaData
		var imageData []byte
		if err := rows.Scan(&book.ID, &book.Name, &imageData); err != nil {
			return nil, err
		}
		book.CoverImage = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData)

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return books, nil

}
func (r *BookRepositoryImpl) SaveUserBookProfile(profile models.UserBookProfile) error {
	_, err := r.DB.Exec(
		`INSERT INTO book_user_profiles (user_id, book_id, book_name, mode, total_pages, current_page)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		profile.UserID, profile.BookID, profile.BookName, profile.Mode,
		profile.TotalPages, profile.CurrentPage,
	)
	return err
}

func (r *BookRepositoryImpl) DeleteBook(bookID int, userID int) error {
	_, err := r.DB.Exec("DELETE FROM books WHERE id = $1 AND user_id = $2", bookID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BookRepositoryImpl) GetBookName(userID, bookID int) (string, error) {
	ctx := context.Background()
	cacheKey := r.getBookNameKey(userID, bookID)

	// 1. Check Redis cache
	data, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var profile models.UserBookProfile
		if unmarshalErr := json.Unmarshal([]byte(data), &profile); unmarshalErr == nil {
			return profile.BookName, nil
		}
	}

	// 2. If not found in Redis, check DB
	var nameFromDB string
	err = r.DB.QueryRowContext(ctx, "SELECT name FROM books WHERE id = $1", bookID).Scan(&nameFromDB)

	if err != nil {
		return "", err
	}

	return nameFromDB, nil
}

func (r *BookRepositoryImpl) getBookNameKey(userID, bookID int) string {
	return fmt.Sprintf("user:%d:book:%d:profile", userID, bookID)
}

func (r *BookRepositoryImpl) GetBookIndexingJob(bookID int, userID int) (models.BookIndexingJob, error) {
	var bookIndexingJob models.BookIndexingJob
	err := r.DB.QueryRow(
		`SELECT id, book_id, user_id, status, created_at, updated_at FROM book_indexing_jobs WHERE book_id = $1 AND user_id = $2`,
		bookID, userID,
	).Scan(
		&bookIndexingJob.ID,
		&bookIndexingJob.BookID,
		&bookIndexingJob.UserID,
		&bookIndexingJob.Status,
		&bookIndexingJob.CreatedAt,
		&bookIndexingJob.UpdatedAt,
	)
	if err != nil {
		return models.BookIndexingJob{}, err
	}
	return bookIndexingJob, nil
}
