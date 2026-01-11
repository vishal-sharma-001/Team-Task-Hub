package service

import (
	"context"
	"log"

	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/repository"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

// UserService defines user-related business logic operations
type UserService interface {
	SignUp(ctx context.Context, email, password string) (*domain.User, string, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, email string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// SignUp creates a new user account
func (s *userService) SignUp(ctx context.Context, email, password string) (*domain.User, string, error) {
	// Validate inputs
	log.Printf("[Service.SignUp] Starting signup for email: %s", email)

	if appErr := utils.ValidateEmail(email); appErr != nil {
		log.Printf("[Service.SignUp] Email validation failed: %v", appErr)
		return nil, "", appErr
	}

	if appErr := utils.ValidatePassword(password); appErr != nil {
		log.Printf("[Service.SignUp] Password validation failed: %v", appErr)
		return nil, "", appErr
	}

	// Hash password
	log.Printf("[Service.SignUp] Hashing password")
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("[Service.SignUp] Hash error: %v", err)
		return nil, "", apperrors.NewInternalError("failed to hash password", err)
	}

	// Create user in database
	log.Printf("[Service.SignUp] Creating user in database")
	user, err := s.userRepo.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		log.Printf("[Service.SignUp] CreateUser error: %v, Type: %T", err, err)
		return nil, "", err
	}

	// Generate JWT token
	log.Printf("[Service.SignUp] Generating token")
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("[Service.SignUp] Token generation error: %v", err)
		return nil, "", apperrors.NewInternalError("failed to generate token", err)
	}

	log.Printf("[Service.SignUp] SignUp successful for user: %s", user.ID)
	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	// Validate inputs
	if appErr := utils.ValidateEmail(email); appErr != nil {
		return nil, "", appErr
	}

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	// Verify password
	if !utils.VerifyPassword(user.PasswordHash, password) {
		return nil, "", apperrors.NewAuthError(apperrors.ErrInvalidPassword, "invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", apperrors.NewInternalError("failed to generate token", err)
	}

	return user, token, nil
}

// GetProfile retrieves the current user's profile
func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile updates the current user's profile
func (s *userService) UpdateProfile(ctx context.Context, userID string, name string) (*domain.User, error) {
	// Update user in database
	user, err := s.userRepo.UpdateUser(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers retrieves all users (for assignee selection)
func (s *userService) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
