package service

import (
	"context"
	"fmt"

	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"

	"github.com/google/uuid"
)

type genreService struct {
	genreRepo repository.GenreRepository
}

func NewGenreService(genreRepo repository.GenreRepository) GenreService {
	return &genreService{
		genreRepo: genreRepo,
	}
}

func (s *genreService) Create(ctx context.Context, req *models.CreateGenreRequest) (*models.Genre, error) {
	existingGenre, _ := s.genreRepo.GetByName(ctx, req.Name)
	if existingGenre != nil {
		return nil, fmt.Errorf("genre with name %s already exists", req.Name)
	}

	genre := &models.Genre{
		ID:   uuid.New(),
		Name: req.Name,
	}

	if err := s.genreRepo.Create(ctx, genre); err != nil {
		return nil, fmt.Errorf("failed to create genre: %w", err)
	}

	return genre, nil
}

func (s *genreService) GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	return s.genreRepo.GetByID(ctx, id)
}

func (s *genreService) GetAll(ctx context.Context) ([]models.Genre, error) {
	return s.genreRepo.GetAll(ctx)
}

func (s *genreService) Update(ctx context.Context, id uuid.UUID, req *models.CreateGenreRequest) (*models.Genre, error) {
	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existingGenre, _ := s.genreRepo.GetByName(ctx, req.Name)
		if existingGenre != nil && existingGenre.ID != id {
			return nil, fmt.Errorf("genre with name %s already exists", req.Name)
		}
		genre.Name = req.Name
	}

	if err := s.genreRepo.Update(ctx, genre); err != nil {
		return nil, fmt.Errorf("failed to update genre: %w", err)
	}

	return genre, nil
}

func (s *genreService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.genreRepo.Delete(ctx, id)
}

