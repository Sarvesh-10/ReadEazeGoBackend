package domain

import (
	"database/sql"
	"encoding/base64"

	models "github.com/Sarvesh-10/ReadEazeBackend/internal/models"
)

type BookRepository interface {
	SaveBook(book Book) (int, error)
	GetBookByID(bookID int, userID int) (*Book, error)
	GetBooksMetaByUser(UserId int) ([]BookMetaData, error)
	SaveUserBookProfile(profile models.UserBookProfile) error
	DeleteBook(bookID int, userID int) error
}
type BookRepositoryImpl struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepositoryImpl {
	return &BookRepositoryImpl{DB: db}
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
