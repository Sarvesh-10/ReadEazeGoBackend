// internal/services/status_service.go
package services

import (
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
)

type StatusService struct {
	bookRepo *domain.BookRepositoryImpl
}

func NewStatusService(bookRepo *domain.BookRepositoryImpl) *StatusService {
	return &StatusService{bookRepo: bookRepo}
}

func (s *StatusService) ProcessJob(job models.BookIndexingJob) domain.SSEMessage {
	msg := domain.SSEMessage{
		BookID: job.BookID,
		Status: domain.JobStatus(job.Status),
	}
	bookName, err := s.bookRepo.GetBookName(job.UserID, job.BookID)

	if err == nil {
		msg.Name = bookName
	}
	return msg
}
