package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"movie-review-api/internal/models"

	"github.com/google/uuid"
)

type movieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, movie *models.Movie) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO movies (id, title, description, release_year, director, duration_minutes)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(ctx, query, movie.ID, movie.Title, movie.Description, movie.ReleaseYear, movie.Director, movie.DurationMinutes)
	if err != nil {
		return err
	}

	if len(movie.Genres) > 0 {
		genreIDs := make([]uuid.UUID, len(movie.Genres))
		for i, genre := range movie.Genres {
			genreIDs[i] = genre.ID
		}
		if err := r.setGenresTx(ctx, tx, movie.ID, genreIDs); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *movieRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	query := `
		SELECT id, title, description, release_year, director, duration_minutes, average_rating, created_at, updated_at 
		FROM movies WHERE id = $1
	`
	movie := &models.Movie{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear, &movie.Director,
		&movie.DurationMinutes, &movie.AverageRating, &movie.CreatedAt, &movie.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, err
	}

	genres, err := r.GetGenres(ctx, id)
	if err != nil {
		return nil, err
	}
	movie.Genres = genres

	return movie, nil
}

func (r *movieRepository) GetAll(ctx context.Context, filters models.MovieFilters, pagination models.PaginationParams) ([]models.Movie, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filters.Search != "" {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	if filters.Year > 0 {
		whereClause += fmt.Sprintf(" AND release_year = $%d", argIndex)
		args = append(args, filters.Year)
		argIndex++
	}

	if filters.MinRating > 0 {
		whereClause += fmt.Sprintf(" AND average_rating >= $%d", argIndex)
		args = append(args, filters.MinRating)
		argIndex++
	}

	if filters.Genre != "" {
		whereClause += fmt.Sprintf(`
			AND id IN (
				SELECT mg.movie_id FROM movie_genres mg 
				JOIN genres g ON mg.genre_id = g.id 
				WHERE g.name ILIKE $%d
			)`, argIndex)
		args = append(args, "%"+filters.Genre+"%")
		argIndex++
	}
	if filters.GenreID != nil {
		whereClause += fmt.Sprintf(`
			AND id IN (
				SELECT movie_id FROM movie_genres 
				WHERE genre_id = $%d
			)`, argIndex)
		args = append(args, *filters.GenreID)
		argIndex++
	}

	countQuery := "SELECT COUNT(*) FROM movies " + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := `
		SELECT id, title, description, release_year, director, duration_minutes, average_rating, created_at, updated_at 
		FROM movies ` + whereClause + `
		ORDER BY created_at DESC 
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, pagination.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear, &movie.Director,
			&movie.DurationMinutes, &movie.AverageRating, &movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		genres, err := r.GetGenres(ctx, movie.ID)
		if err != nil {
			return nil, 0, err
		}
		movie.Genres = genres

		movies = append(movies, movie)
	}

	return movies, total, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *models.Movie) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE movies 
		SET title = $2, description = $3, release_year = $4, director = $5, duration_minutes = $6, updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, query, movie.ID, movie.Title, movie.Description, movie.ReleaseYear, movie.Director, movie.DurationMinutes)
	if err != nil {
		return err
	}

	if len(movie.Genres) > 0 {
		genreIDs := make([]uuid.UUID, len(movie.Genres))
		for i, genre := range movie.Genres {
			genreIDs[i] = genre.ID
		}
		if err := r.setGenresTx(ctx, tx, movie.ID, genreIDs); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *movieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM movies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *movieRepository) UpdateAverageRating(ctx context.Context, movieID uuid.UUID, rating float64) error {
	query := `UPDATE movies SET average_rating = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, movieID, rating)
	return err
}

func (r *movieRepository) GetGenres(ctx context.Context, movieID uuid.UUID) ([]models.Genre, error) {
	query := `
		SELECT g.id, g.name, g.created_at 
		FROM genres g 
		JOIN movie_genres mg ON g.id = mg.genre_id 
		WHERE mg.movie_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, movieID)
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

func (r *movieRepository) SetGenres(ctx context.Context, movieID uuid.UUID, genreIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.setGenresTx(ctx, tx, movieID, genreIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *movieRepository) setGenresTx(ctx context.Context, tx *sql.Tx, movieID uuid.UUID, genreIDs []uuid.UUID) error {
	deleteQuery := `DELETE FROM movie_genres WHERE movie_id = $1`
	_, err := tx.ExecContext(ctx, deleteQuery, movieID)
	if err != nil {
		return err
	}

	if len(genreIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(genreIDs))
	args := make([]interface{}, len(genreIDs)*2)
	for i, genreID := range genreIDs {
		placeholders[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		args[i*2] = movieID
		args[i*2+1] = genreID
	}

	insertQuery := `INSERT INTO movie_genres (movie_id, genre_id) VALUES ` + strings.Join(placeholders, ", ")
	_, err = tx.ExecContext(ctx, insertQuery, args...)
	return err
}
