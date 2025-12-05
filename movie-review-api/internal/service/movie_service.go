package service

import (
	"context"
	"fmt"

	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"

	"github.com/google/uuid"
)

type movieService struct {
	movieRepo repository.MovieRepository
	genreRepo repository.GenreRepository
}

func NewMovieService(movieRepo repository.MovieRepository, genreRepo repository.GenreRepository) MovieService {
	return &movieService{
		movieRepo: movieRepo,
		genreRepo: genreRepo,
	}
}

func (s *movieService) Create(ctx context.Context, req *models.CreateMovieRequest) (*models.Movie, error) {
	genreIDs := make([]uuid.UUID, len(req.GenreIDs))
	for i, genreIDStr := range req.GenreIDs {
		genreID, err := uuid.Parse(genreIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid genre ID: %s", genreIDStr)
		}
		genreIDs[i] = genreID
	}

	for _, genreID := range genreIDs {
		_, err := s.genreRepo.GetByID(ctx, genreID)
		if err != nil {
			return nil, fmt.Errorf("genre with ID %s not found", genreID)
		}
	}

	movie := &models.Movie{
		ID:              uuid.New(),
		Title:           req.Title,
		Description:     req.Description,
		ReleaseYear:     req.ReleaseYear,
		Director:        req.Director,
		DurationMinutes: req.DurationMinutes,
		AverageRating:   0.0,
	}

	genres := make([]models.Genre, len(genreIDs))
	for i, genreID := range genreIDs {
		genre, err := s.genreRepo.GetByID(ctx, genreID)
		if err != nil {
			return nil, err
		}
		genres[i] = *genre
	}
	movie.Genres = genres

	if err := s.movieRepo.Create(ctx, movie); err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	return movie, nil
}

func (s *movieService) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	return s.movieRepo.GetByID(ctx, id)
}

func (s *movieService) GetAll(ctx context.Context, filters models.MovieFilters, pagination models.PaginationParams) ([]models.Movie, int, error) {
	return s.movieRepo.GetAll(ctx, filters, pagination)
}

func (s *movieService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateMovieRequest) (*models.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.ReleaseYear > 0 {
		movie.ReleaseYear = req.ReleaseYear
	}
	if req.Director != "" {
		movie.Director = req.Director
	}
	if req.DurationMinutes > 0 {
		movie.DurationMinutes = req.DurationMinutes
	}

	if len(req.GenreIDs) > 0 {
		genreIDs := make([]uuid.UUID, len(req.GenreIDs))
		for i, genreIDStr := range req.GenreIDs {
			genreID, err := uuid.Parse(genreIDStr)
			if err != nil {
				return nil, fmt.Errorf("invalid genre ID: %s", genreIDStr)
			}
			genreIDs[i] = genreID
		}

		for _, genreID := range genreIDs {
			_, err := s.genreRepo.GetByID(ctx, genreID)
			if err != nil {
				return nil, fmt.Errorf("genre with ID %s not found", genreID)
			}
		}

		genres := make([]models.Genre, len(genreIDs))
		for i, genreID := range genreIDs {
			genre, err := s.genreRepo.GetByID(ctx, genreID)
			if err != nil {
				return nil, err
			}
			genres[i] = *genre
		}
		movie.Genres = genres
	}

	if err := s.movieRepo.Update(ctx, movie); err != nil {
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	return movie, nil
}

func (s *movieService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.movieRepo.Delete(ctx, id)
}

func (s *movieService) UpdateAverageRating(ctx context.Context, movieID uuid.UUID) error {
	return s.movieRepo.UpdateAverageRating(ctx, movieID, 0.0)
}

