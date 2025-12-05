package service

import (
	"context"
	"log/slog"
	"sync"

	"movie-review-api/internal/repository"

	"github.com/google/uuid"
)

type RatingWorker struct {
	reviewRepo  repository.ReviewRepository
	movieRepo   repository.MovieRepository
	queue       chan uuid.UUID
	workerCount int
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewRatingWorker(reviewRepo repository.ReviewRepository, movieRepo repository.MovieRepository, workerCount int) *RatingWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &RatingWorker{
		reviewRepo:  reviewRepo,
		movieRepo:   movieRepo,
		queue:       make(chan uuid.UUID, 100),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (w *RatingWorker) Start() {
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.worker()
	}
	slog.Info("Rating worker started", "workers", w.workerCount)
}

func (w *RatingWorker) Stop() {
	w.cancel()
	close(w.queue)
	w.wg.Wait()
	slog.Info("Rating worker stopped")
}

func (w *RatingWorker) QueueRatingUpdate(movieID uuid.UUID) {
	select {
	case w.queue <- movieID:
	case <-w.ctx.Done():
		return
	default:
		slog.Warn("Rating update queue is full, skipping", "movie_id", movieID)
	}
}

func (w *RatingWorker) worker() {
	defer w.wg.Done()
	for {
		select {
		case movieID, ok := <-w.queue:
			if !ok {
				return
			}
			w.processRatingUpdate(movieID)
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *RatingWorker) processRatingUpdate(movieID uuid.UUID) {
	ctx := context.Background()

	avgRating, err := w.reviewRepo.GetAverageRating(ctx, movieID)
	if err != nil {
		slog.Error("Failed to get average rating", "movie_id", movieID, "error", err)
		return
	}

	if err := w.movieRepo.UpdateAverageRating(ctx, movieID, avgRating); err != nil {
		slog.Error("Failed to update average rating", "movie_id", movieID, "error", err)
		return
	}

	slog.Debug("Updated average rating", "movie_id", movieID, "rating", avgRating)
}
