package app

import (
	"errors"

	domain "github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   domain.UserRepository
	logger utility.Logger
}

func NewUserService(repo domain.UserRepository, logger *utility.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: *logger,
	}
}

func (s *UserService) RegisterUser(user domain.User) (int64, string, error) {
	s.logger.Info("Registering user")
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil || err != nil {
		if err != nil {
			s.logger.Error("Error getting user: %s", err.Error())
			return -1, "", err
		}
		s.logger.Error("User already exists with email: %s", user.Email)
		return -1, "", errors.New("user already exists")
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		s.logger.Error("Error hashing password: %s", hashErr.Error())
		return -1, "", errors.New("internal server error")
	}
	user.Password = string(hashedPassword)

	s.logger.Info("Creating a new user with email: %s", user.Email)
	id, create_err := s.repo.CreateUser(&user)
	if create_err != nil {
		s.logger.Error("Error creating user: %s", create_err.Error())
		return -1, "", create_err
	}

	token, token_err := utility.GenerateToken(user.ID, user.Email)
	if token_err != nil {
		s.logger.Error("Error generating token: %s", token_err.Error())
		return -1, "", token_err
	}
	s.logger.Info("User Registered successfully")
	return id, token, nil

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
