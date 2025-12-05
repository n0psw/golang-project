package repository

import (
	"context"
	"database/sql"
	"fmt"

	"movie-review-api/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MovieRepository interface {
	Create(ctx context.Context, movie *models.Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error)
	GetAll(ctx context.Context, filters models.MovieFilters, pagination models.PaginationParams) ([]models.Movie, int, error)
	Update(ctx context.Context, movie *models.Movie) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateAverageRating(ctx context.Context, movieID uuid.UUID, rating float64) error
	GetGenres(ctx context.Context, movieID uuid.UUID) ([]models.Genre, error)
	SetGenres(ctx context.Context, movieID uuid.UUID, genreIDs []uuid.UUID) error
}

type GenreRepository interface {
	Create(ctx context.Context, genre *models.Genre) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error)
	GetByName(ctx context.Context, name string) (*models.Genre, error)
	GetAll(ctx context.Context) ([]models.Genre, error)
	Update(ctx context.Context, genre *models.Genre) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetAll(ctx context.Context, pagination models.PaginationParams) ([]models.Review, int, error)
	GetByMovieID(ctx context.Context, movieID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error)
	GetByMovieAndUser(ctx context.Context, movieID, userID uuid.UUID) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAverageRating(ctx context.Context, movieID uuid.UUID) (float64, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.PasswordHash, user.Role)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, role, created_at, updated_at FROM users WHERE username = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = $2, username = $3, password_hash = $4, role = $5, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.PasswordHash, user.Role)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
