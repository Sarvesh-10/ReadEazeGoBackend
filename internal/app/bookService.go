package app

import (
	"bytes"
	"image/jpeg"
	"os"
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/gen2brain/go-fitz"
)

type BookService struct {
	BookRepo *domain.BookRepository
	logger   utility.Logger
}

func NewBookService(bookRepo *domain.BookRepository, logger *utility.Logger) *BookService {
	return &BookService{BookRepo: bookRepo, logger: *logger}
}

func (s *BookService) UploadBook(userID int, name string, pdfData []byte) error {
	tmfFile, err := os.CreateTemp("", "book_*.pdf")
	if err != nil {
		s.logger.Error("Error creating temporary file: %s", err.Error())
		return err
	}
	defer os.Remove(tmfFile.Name()) // Clean up the temp file after use
	defer tmfFile.Close()

	if _, err := tmfFile.Write(pdfData); err != nil {
		s.logger.Error("Error writing PDF to temp file: %s", err.Error())
		return err
	}

	doc, err := fitz.New(tmfFile.Name())
	if err != nil {
		s.logger.Error("Error creating fitz document: %s", err.Error())
		return err
	}
	defer doc.Close()
	coverImage, err := doc.Image(0)
	if err != nil {
		s.logger.Error("Error extracting cover image: %s", err.Error())
		return err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, coverImage, &jpeg.Options{Quality: 100})
	if err != nil {
		s.logger.Error("Error encoding cover image: %s", err.Error())
		return err
	}
	coverImageData := buf.Bytes()
	book := domain.Book{
		UserID:     userID,
		Name:       name,
		PDFData:    pdfData,
		UploadedAt: time.Now(),
		CoverImage: coverImageData,
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
