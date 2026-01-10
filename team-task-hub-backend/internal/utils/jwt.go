package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/launchventures/team-task-hub-backend/internal/errors"
)

const (
	TokenExpiration = 24 * time.Hour
	JWTSecret       = "your-secret-key-change-in-production"
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func GenerateToken(userID int, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns claims
func ValidateToken(tokenString string) (*JWTClaims, *errors.AppError) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, errors.NewAuthError(errors.ErrInvalidToken, "invalid token")
	}

	if !token.Valid {
		return nil, errors.NewAuthError(errors.ErrInvalidToken, "token is not valid")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewAuthError(errors.ErrTokenExpired, "token has expired")
	}

	return claims, nil
}

// ExtractUserIDFromToken extracts user ID from JWT token string
func ExtractUserIDFromToken(tokenString string) (int, *errors.AppError) {
	claims, appErr := ValidateToken(tokenString)
	if appErr != nil {
		return 0, appErr
	}
	return claims.UserID, nil
}

// ExtractUserIDFromContext extracts user ID from request context
func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return 0, errors.NewAuthError(errors.ErrUnauthorized, "user ID not found in context")
	}
	return userID, nil
}
