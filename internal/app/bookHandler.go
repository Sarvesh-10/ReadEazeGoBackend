package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/middleware"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/services"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/gorilla/mux"
	"rsc.io/pdf"
)

type BookHandler struct {
	BookService     *BookService
	UserBookProfile *services.UserBookProfileService
	logger          utility.Logger
}

func NewBookHandler(service *BookService, UserBookProfile *services.UserBookProfileService, logger *utility.Logger) *BookHandler {
	return &BookHandler{
		BookService:     service,
		UserBookProfile: UserBookProfile,
		logger:          *logger,
	}
}

func (h *BookHandler) UploadBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	// ctx := r.Context()
	h.logger.Info("HERE IN UPLoad books")
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	mode := r.FormValue("mode")
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

	book, err := h.BookService.UploadBook(userID, header.Filename, pdfData)
	if err != nil {
		http.Error(w, "Failed to save book", http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	// h.logger.Info("Book uploaded successfully with ID: %d", bookID)
	totalPageCount, err := h.getPDFPageCount(pdfData)
	if err != nil {
		// http.Error(w, "Failed to get page count", http.StatusInternalServerError)
		totalPageCount = -1
		h.logger.Error("Failed to get page count: %s", err.Error())

	}
	userBookProfile := h.UserBookProfile.CreateUserBookProfile(book, mode, totalPageCount)
	err = h.UserBookProfile.SaveUserBookProfileInCache(userBookProfile)
	if err != nil {
		http.Error(w, "Failed to save user book profile", http.StatusInternalServerError)
		h.logger.Error("Failed to save user book profile: %s", err.Error())
		// return
	}
	err = h.UserBookProfile.SaveUserBookProfileInDB(userBookProfile)
	if err != nil {
		http.Error(w, "Failed to save user book profile in DB", http.StatusInternalServerError)
		h.logger.Error("Failed to save user book profile in DB: %s", err.Error())
		//delete book from repository
		errDelete := h.BookService.DeleteBook(book.ID, userID)
		if errDelete != nil {
			h.logger.Error("Failed to delete book from repository after profile save failure: %s", errDelete.Error())
		}
		return
	}

	//save this in the redis with ttl of 2 hours
	// save in the postgres db

	if userBookProfile.Mode == "study" {
		// Create the job
		job := models.BookIndexingJob{
			BookID: book.ID,
			UserID: userID,
			Status: models.JobStatusPending,
		}
		// Save job to DB as PENDING
		jobID, err := h.BookService.BookRepo.SaveBookIndexingJob(job)
		if err != nil {
			h.logger.Error("Failed to save indexing job: %s", err.Error())
			http.Error(w, "Failed to create indexing job", http.StatusInternalServerError)
			return
		}
		job.ID = jobID

		// Push job to Redis queue "books"
		jobBytes, err := json.Marshal(job)
		if err != nil {
			h.logger.Error("Failed to marshal job: %s", err.Error())
			http.Error(w, "Failed to queue indexing job", http.StatusInternalServerError)
			return
		}
		err = h.BookService.cache.PushToQueue("book_indexing_queue", jobBytes)
		if err != nil {
			h.logger.Error("Failed to push job to Redis queue: %s", err.Error())
			http.Error(w, "Failed to queue indexing job", http.StatusInternalServerError)
			return
		}
		h.BookService.cache.PushToQueue("output_jobs_queue", jobBytes)

	}

	response := map[string]string{"message": "Book uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "201 Created")
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

func (h *BookHandler) getPDFPageCount(data []byte) (int, error) {
	h.logger.Info("pdf data", string(data[:10]))
	reader := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return 0, err
	}
	return pdfReader.NumPage(), nil
}
