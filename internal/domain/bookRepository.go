package domain

import (
	"database/sql"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) SaveBook(book Book) error {
	_, err := r.DB.Exec(
		"INSERT INTO books (user_id, name, file_data, uploaded_at) VALUES ($1, $2, $3, $4)",
		book.UserID, book.Name, book.PDFData, book.UploadedAt,
	)
	return err
}

func (r *BookRepository) GetBookByID(bookID int, userID int) (*Book, error) {
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

func (r *BookRepository) GetBooksMetaByUser(UserId int) ([]BookMetaData, error) {
	var books []BookMetaData
	rows, err := r.DB.Query("SELECT id, name FROM books WHERE user_id = $1", UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book BookMetaData
		if err := rows.Scan(&book.ID, &book.Name); err != nil {
			return nil, err
		}
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
