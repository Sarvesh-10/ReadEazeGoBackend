package app

import (
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type BookService struct {
	BookRepo *domain.BookRepository
	logger   utility.Logger
}

func NewBookService(bookRepo *domain.BookRepository, logger *utility.Logger) *BookService {
	return &BookService{BookRepo: bookRepo, logger: *logger}
}

func (s *BookService) UploadBook(userID int, name string, pdfData []byte) error {
	book := domain.Book{
		UserID:     userID,
		Name:       name,
		PDFData:    pdfData,
		UploadedAt: time.Now(),
	}
	return s.BookRepo.SaveBook(book)
}

func (s *BookService) GetBook(bookID, userID int) (*domain.Book, error) {
	return s.BookRepo.GetBookByID(bookID, userID)
}

func (s *BookService) GetBookMetadataByUser(userID int) ([]domain.BookMetaData, error) {
	books, err := s.BookRepo.GetBooksMetaByUser(userID)
	if err != nil {
		s.logger.Error("Some error in getting the books this is isn service layer")
		return []domain.BookMetaData{}, err
	}
	return books, nil
}
