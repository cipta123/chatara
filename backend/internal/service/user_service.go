package service

import (
	"errors"
	"strings"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/repository"
	"github.com/umess/backend/pkg/auth"
	"github.com/umess/backend/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Clear password before returning
	user.Password = ""

	return &model.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	// Trim whitespace from email
	email := strings.TrimSpace(req.Email)
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Clear password before returning
	user.Password = ""

	return &model.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) GetUser(userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}

func (s *UserService) SearchUsers(query string, excludeUserID int) ([]*model.User, error) {
	if query == "" {
		return []*model.User{}, nil
	}
	
	users, err := s.userRepo.SearchUsers(query, excludeUserID, 20)
	if err != nil {
		return nil, err
	}

	// Clear passwords
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}