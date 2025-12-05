package repository

import (
	"context"
	"database/sql"
	"fmt"

	"movie-review-api/internal/models"

	"github.com/google/uuid"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (id, movie_id, user_id, rating, title, content)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, review.ID, review.MovieID, review.UserID, review.Rating, review.Title, review.Content)
	return err
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	query := `
		SELECT r.id, r.movie_id, r.user_id, r.rating, r.title, r.content, r.created_at, r.updated_at,
		       u.id, u.email, u.username, u.role, u.created_at, u.updated_at,
		       m.id, m.title, m.description, m.release_year, m.director, m.duration_minutes, m.average_rating, m.created_at, m.updated_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN movies m ON r.movie_id = m.id
		WHERE r.id = $1
	`
	review := &models.Review{}
	user := &models.User{}
	movie := &models.Movie{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt,
		&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear, &movie.Director, &movie.DurationMinutes, &movie.AverageRating, &movie.CreatedAt, &movie.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, err
	}

	review.User = user
	review.Movie = movie
	return review, nil
}

func (r *reviewRepository) GetByMovieID(ctx context.Context, movieID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	countQuery := `SELECT COUNT(*) FROM reviews WHERE movie_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, movieID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := `
		SELECT r.id, r.movie_id, r.user_id, r.rating, r.title, r.content, r.created_at, r.updated_at,
		       u.id, u.email, u.username, u.role, u.created_at, u.updated_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.movie_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, movieID, pagination.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var user models.User
		err := rows.Scan(
			&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		review.User = &user
		reviews = append(reviews, review)
	}

	return reviews, total, nil
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	countQuery := `SELECT COUNT(*) FROM reviews WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := `
		SELECT r.id, r.movie_id, r.user_id, r.rating, r.title, r.content, r.created_at, r.updated_at,
		       m.id, m.title, m.description, m.release_year, m.director, m.duration_minutes, m.average_rating, m.created_at, m.updated_at
		FROM reviews r
		JOIN movies m ON r.movie_id = m.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, pagination.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var movie models.Movie
		err := rows.Scan(
			&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt,
			&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear, &movie.Director, &movie.DurationMinutes, &movie.AverageRating, &movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		review.Movie = &movie
		reviews = append(reviews, review)
	}

	return reviews, total, nil
}

func (r *reviewRepository) GetAll(ctx context.Context, pagination models.PaginationParams) ([]models.Review, int, error) {
	countQuery := `SELECT COUNT(*) FROM reviews`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := `
		SELECT r.id, r.movie_id, r.user_id, r.rating, r.title, r.content, r.created_at, r.updated_at,
		       u.id, u.email, u.username, u.role, u.created_at, u.updated_at,
		       m.id, m.title, m.description, m.release_year, m.director, m.duration_minutes, m.average_rating, m.created_at, m.updated_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN movies m ON r.movie_id = m.id
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var user models.User
		var movie models.Movie
		err := rows.Scan(
			&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt,
			&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear, &movie.Director, &movie.DurationMinutes, &movie.AverageRating, &movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		review.User = &user
		review.Movie = &movie
		reviews = append(reviews, review)
	}

	return reviews, total, nil
}

func (r *reviewRepository) GetByMovieAndUser(ctx context.Context, movieID, userID uuid.UUID) (*models.Review, error) {
	query := `
		SELECT id, movie_id, user_id, rating, title, content, created_at, updated_at
		FROM reviews WHERE movie_id = $1 AND user_id = $2
	`
	review := &models.Review{}
	err := r.db.QueryRowContext(ctx, query, movieID, userID).Scan(
		&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, err
	}
	return review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews 
		SET rating = $2, title = $3, content = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, review.ID, review.Rating, review.Title, review.Content)
	return err
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, movieID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE movie_id = $1`
	var avgRating float64
	err := r.db.QueryRowContext(ctx, query, movieID).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}
