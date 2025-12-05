package repository

import (
	"context"
	"database/sql"
	"fmt"

	"movie-review-api/internal/models"

	"github.com/google/uuid"
)

type genreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) GenreRepository {
	return &genreRepository{db: db}
}

func (r *genreRepository) Create(ctx context.Context, genre *models.Genre) error {
	query := `INSERT INTO genres (id, name) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, genre.ID, genre.Name)
	return err
}

func (r *genreRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	query := `SELECT id, name, created_at FROM genres WHERE id = $1`
	genre := &models.Genre{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&genre.ID, &genre.Name, &genre.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("genre not found")
		}
		return nil, err
	}
	return genre, nil
}

func (r *genreRepository) GetByName(ctx context.Context, name string) (*models.Genre, error) {
	query := `SELECT id, name, created_at FROM genres WHERE name = $1`
	genre := &models.Genre{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&genre.ID, &genre.Name, &genre.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("genre not found")
		}
		return nil, err
	}
	return genre, nil
}

func (r *genreRepository) GetAll(ctx context.Context) ([]models.Genre, error) {
	query := `SELECT id, name, created_at FROM genres ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.Name, &genre.CreatedAt)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

func (r *genreRepository) Update(ctx context.Context, genre *models.Genre) error {
	query := `UPDATE genres SET name = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, genre.ID, genre.Name)
	return err
}

func (r *genreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM genres WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

