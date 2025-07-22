package services

import (
	"context"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain/repository"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type UserBookProfileService struct {
	BookRepo *domain.BookRepositoryImpl
	logger   utility.Logger
	cache    repository.CacheRepository
}

func NewUserBookProfileService(bookRepo *domain.BookRepositoryImpl, logger *utility.Logger, cache repository.CacheRepository) *UserBookProfileService {
	return &UserBookProfileService{
		BookRepo: bookRepo,
		logger:   *logger,
		cache:    cache,
	}
}

func (s *UserBookProfileService) CreateUserBookProfile(Book domain.Book, mode string, totalPages int) models.UserBookProfile {
	return models.UserBookProfile{
		UserID:            Book.UserID,
		BookID:            Book.ID,
		BookName:          Book.Name,
		Mode:              mode,
		TotalPages:        totalPages,
		CurrentPage:       0,
		ReadingPercentage: 0.0,
	}
}
func (s *UserBookProfileService) SaveUserBookProfileInCache(profile models.UserBookProfile) error {
	ctx := context.Background()
	err := s.cache.SaveUserBookProfile(ctx, profile.UserID, profile.BookID, profile)
	if err != nil {
		s.logger.Error("Failed to save user book profile in cache: %s", err.Error())
		return err
	}
	return nil
}

func (s *UserBookProfileService) SaveUserBookProfileInDB(profile models.UserBookProfile) error {
	err := s.BookRepo.SaveUserBookProfile(profile)
	if err != nil {
		s.logger.Error("Failed to save user book profile in DB: %s", err.Error())
		return err
	}
	return nil
}
