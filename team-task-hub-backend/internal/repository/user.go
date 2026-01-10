package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// UserRepository defines user data access operations
type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, email string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	const query = `
		INSERT INTO users (email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, '', NOW(), NOW())
		RETURNING id, email, name, password_hash, created_at, updated_at
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, apperrors.NewConflictError(apperrors.ErrEmailExists, "email already exists")
		}
		return nil, apperrors.NewDatabaseError("failed to create user", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrUserNotFound, "user not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get user", err)
	}

	user.Name = "" // Set default empty name

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrUserNotFound, "user not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get user by email", err)
	}

	user.Name = "" // Set default empty name

	return user, nil
}

// ListUsers retrieves all users
func (r *userRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		ORDER BY email ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to list users", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to scan user", err)
		}
		u.Name = "" // Set default empty name
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.NewDatabaseError("error iterating users", err)
	}

	return users, nil
}

// UpdateUser updates a user's profile
func (r *userRepository) UpdateUser(ctx context.Context, id int, name string) (*domain.User, error) {
	const query = `
		UPDATE users
		SET name = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, name, password_hash, created_at, updated_at
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id, name).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrUserNotFound, "user not found")
		}
		if err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, apperrors.NewConflictError(apperrors.ErrEmailExists, "email already exists")
		}
		return nil, apperrors.NewDatabaseError("failed to update user", err)
	}

	return user, nil
}
