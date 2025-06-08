package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/middleware"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/gorilla/mux"
)

type BookHandler struct {
	BookService *BookService
	logger      utility.Logger
}

func NewBookHandler(service *BookService, logger *utility.Logger) *BookHandler {
	return &BookHandler{
		BookService: service,
		logger:      *logger,
	}
}

func (h *BookHandler) UploadBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	h.logger.Info("HERE IN UPLoad books")
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		h.logger.Info(err.Error())
		return
	}
	defer file.Close()

	pdfData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}

	err = h.BookService.UploadBook(userID, header.Filename, pdfData)
	if err != nil {
		http.Error(w, "Failed to save book", http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	response := map[string]string{"message": "Book uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.BookService.GetBook(bookID, userID)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(book.PDFData)
}

func (h *BookHandler) GetBooksMeta(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		h.logger.Error("Invalid User Id")
		return
	}
	books, err := h.BookService.GetBookMetadataByUser(userId)
	if err != nil {
		h.logger.Error("SOme error while fetching list of meta books")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if books == nil {
		books = []domain.BookMetaData{} // Ensure we return an empty array if no books found
	}
	json.NewEncoder(w).Encode(books)
}
