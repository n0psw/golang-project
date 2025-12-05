package service

import (
	"context"
	"fmt"

	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"

	"github.com/google/uuid"
)

type reviewService struct {
	reviewRepo   repository.ReviewRepository
	movieRepo    repository.MovieRepository
	userRepo     repository.UserRepository
	ratingWorker *RatingWorker
}

func NewReviewService(reviewRepo repository.ReviewRepository, movieRepo repository.MovieRepository, userRepo repository.UserRepository, ratingWorker *RatingWorker) ReviewService {
	return &reviewService{
		reviewRepo:   reviewRepo,
		movieRepo:    movieRepo,
		userRepo:     userRepo,
		ratingWorker: ratingWorker,
	}
}

func (s *reviewService) Create(ctx context.Context, movieID uuid.UUID, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error) {
	_, err := s.movieRepo.GetByID(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found")
	}

	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	existingReview, _ := s.reviewRepo.GetByMovieAndUser(ctx, movieID, userID)
	if existingReview != nil {
		return nil, fmt.Errorf("user has already reviewed this movie")
	}

	review := &models.Review{
		ID:      uuid.New(),
		MovieID: movieID,
		UserID:  userID,
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	s.ratingWorker.QueueRatingUpdate(movieID)

	return review, nil
}

func (s *reviewService) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

func (s *reviewService) GetAll(ctx context.Context, pagination models.PaginationParams) ([]models.Review, int, error) {
	return s.reviewRepo.GetAll(ctx, pagination)
}

func (s *reviewService) GetByMovieID(ctx context.Context, movieID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	return s.reviewRepo.GetByMovieID(ctx, movieID, pagination)
}

func (s *reviewService) GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	return s.reviewRepo.GetByUserID(ctx, userID, pagination)
}

func (s *reviewService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateReviewRequest) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if review.UserID != userID {
		return nil, fmt.Errorf("unauthorized: you can only update your own reviews")
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Title != "" {
		review.Title = req.Title
	}
	if req.Content != "" {
		review.Content = req.Content
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	s.ratingWorker.QueueRatingUpdate(review.MovieID)

	return review, nil
}

func (s *reviewService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if review.UserID != userID && user.Role != "admin" {
		return fmt.Errorf("unauthorized: you can only delete your own reviews")
	}

	movieID := review.MovieID
	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.ratingWorker.QueueRatingUpdate(movieID)

	return nil
}
