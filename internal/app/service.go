package app

import (
	"errors"
	"time"

	domain "github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        domain.UserRepository
	logger      utility.Logger
	refreshRepo domain.RefreshTokenRepository
}

func NewUserService(repo domain.UserRepository, logger *utility.Logger, refreshRepo domain.RefreshTokenRepository) *UserService {
	return &UserService{
		repo:        repo,
		refreshRepo: refreshRepo,
		logger:      *logger,
	}
}

func (s *UserService) RegisterUser(user domain.User) (int64, string, string, error) {
	s.logger.Info("Registering user")
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil || err != nil {
		if err != nil {
			s.logger.Error("Error getting user: %s", err.Error())
			return -1, "", "", err
		}
		s.logger.Error("User already exists with email: %s", user.Email)
		return -1, "", "", errors.New("user already exists")
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		s.logger.Error("Error hashing password: %s", hashErr.Error())
		return -1, "", "", errors.New("internal server error")
	}
	user.Password = string(hashedPassword)

	s.logger.Info("Creating a new user with email: %s", user.Email)
	id, create_err := s.repo.CreateUser(&user)
	if create_err != nil {
		s.logger.Error("Error creating user: %s", create_err.Error())
		return -1, "", "", create_err
	}

	token, token_err := utility.GenerateToken(int(id), user.Email)
	s.logger.Info("Generated token for user with ID: %d", id)
	if token_err != nil {
		s.logger.Error("Error generating token: %s", token_err.Error())
		return -1, "", "", token_err
	}
	refreshToken, refreshErr := utility.GenerateRefreshToken()
	if refreshErr != nil {
		s.logger.Error("Error generating refresh token: %s", refreshErr.Error())
		return -1, "", "", refreshErr
	}
	err = s.refreshRepo.Store(int(id), utility.HashToken(refreshToken), time.Now().Add(time.Hour*24*7))
	if err != nil {
		s.logger.Error("Error storing refresh token: %s", err.Error())
		return -1, "", "", err // discard access token too
	}

	s.logger.Info("User Registered successfully")
	return id, token, refreshToken, nil

}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	s.logger.Info("Getting user by email: %s", email)
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.logger.Error("Error getting user: %s", err.Error())
		return nil, err
	}
	s.logger.Info("User retrieved successfully")
	return user, nil
}

func (s *UserService) RefreshSession(rawToken string) (string, string, error) {
	hashedToken := utility.HashToken(rawToken)
	userID, emailID, err := s.refreshRepo.Validate(hashedToken)
	if err != nil {
		s.logger.Error("Error validating refresh token: %s", err.Error())
		return "", "", err
	}
	accessToken, err := utility.GenerateToken(userID, emailID)
	if err != nil {
		s.logger.Error("Error generating access token: %s", err.Error())
		return "", "", err
	}

	newRefreshToken, err := utility.GenerateRefreshToken()
	if err != nil {
		s.logger.Error("Error generating refresh token: %s", err.Error())
		return "", "", err // discard access token too

	}
	newHash := utility.HashToken(newRefreshToken)
	err = s.refreshRepo.Revoke(hashedToken)
	if err != nil {
		s.logger.Error("Error revoking old refresh token: %s", err.Error())
		return "", "", err
	}
	err = s.refreshRepo.Store(userID, newHash, time.Now().Add(time.Hour*24*7))
	if err != nil {
		s.logger.Error("Error storing new refresh token: %s", err.Error())
		return "", "", err
	}
	return accessToken, newRefreshToken, nil

}
