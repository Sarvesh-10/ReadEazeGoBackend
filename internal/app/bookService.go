package app

import (
	"bytes"
	"image/jpeg"
	"os"
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain/repository"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/gen2brain/go-fitz"
)

type BookService struct {
	BookRepo *domain.BookRepositoryImpl
	logger   utility.Logger
	cache    repository.CacheRepository
}

func NewBookService(bookRepo *domain.BookRepositoryImpl, logger *utility.Logger, cache repository.CacheRepository) *BookService {
	return &BookService{BookRepo: bookRepo, logger: *logger, cache: cache}
}

func (s *BookService) UploadBook(userID int, name string, pdfData []byte) (domain.Book, error) {
	tmfFile, err := os.CreateTemp("", "book_*.pdf")
	if err != nil {
		s.logger.Error("Error creating temporary file: %s", err.Error())
		return domain.Book{}, err
	}
	defer os.Remove(tmfFile.Name()) // Clean up the temp file after use
	defer tmfFile.Close()

	if _, err := tmfFile.Write(pdfData); err != nil {
		s.logger.Error("Error writing PDF to temp file: %s", err.Error())
		return domain.Book{}, err
	}

	doc, err := fitz.New(tmfFile.Name())
	if err != nil {
		s.logger.Error("Error creating fitz document: %s", err.Error())
		return domain.Book{}, err
	}
	defer doc.Close()
	coverImage, err := doc.Image(0)
	if err != nil {
		s.logger.Error("Error extracting cover image: %s", err.Error())
		return domain.Book{}, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, coverImage, &jpeg.Options{Quality: 100})
	if err != nil {
		s.logger.Error("Error encoding cover image: %s", err.Error())
		return domain.Book{}, err
	}
	coverImageData := buf.Bytes()
	book := domain.Book{
		UserID:     userID,
		Name:       name,
		PDFData:    pdfData,
		UploadedAt: time.Now(),
		CoverImage: coverImageData,
	}
	book.ID, err = s.BookRepo.SaveBook(book)
	if err != nil {
		s.logger.Error("Error saving book to repository: %s", err.Error())
		return domain.Book{}, err
	}
	return book, nil
}

func (s *BookService) GetBook(bookID, userID int) (*domain.Book, error) {
	return s.BookRepo.GetBookByID(bookID, userID)
}
func (s *BookService) GetBookIndexingJob(bookID, userID int) (models.BookIndexingJob, error) {
	return s.BookRepo.GetBookIndexingJob(bookID, userID)
}

func (s *BookService) GetBookMetadataByUser(userID int) ([]domain.BookMetaData, error) {
	books, err := s.BookRepo.GetBooksMetaByUser(userID)
	if err != nil {
		s.logger.Error("Some error in getting the books this is isn service layer")
		return []domain.BookMetaData{}, err
	}
	return books, nil
}

// Delete book
func (s *BookService) DeleteBook(bookID, userID int) error {
	err := s.BookRepo.DeleteBook(bookID, userID)
	if err != nil {
		s.logger.Error("Failed to delete book: %s", err.Error())
		return err
	}
	// Optionally, delete the user book profile from cache
	err = s.cache.DeleteUserBookProfile(userID, bookID)
	if err != nil {
		s.logger.Error("Failed to delete user book profile from cache: %s", err.Error())
	}
	return nil
}
