package repository

import (
	"context"
	"fmt"
	"sync"

	"movie-review-api/internal/models"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewMockUserRepository() UserRepository {
	repo := &MockUserRepository{
		users: make(map[string]*models.User),
	}

	repo.users["admin@example.com"] = &models.User{
		ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:        "admin@example.com",
		Username:     "admin",
		PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
		Role:         "admin",
	}

	return repo
}

func (r *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if user, exists := r.users[email]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for email, user := range r.users {
		if user.ID == id {
			delete(r.users, email)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

type MockMovieRepository struct {
	movies map[string]*models.Movie
	mutex  sync.RWMutex
}

func NewMockMovieRepository() MovieRepository {
	repo := &MockMovieRepository{
		movies: make(map[string]*models.Movie),
	}

	movie1 := &models.Movie{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Title:           "The Matrix",
		Description:     "A computer hacker learns about the true nature of reality",
		ReleaseYear:     1999,
		Director:        "Lana Wachowski",
		DurationMinutes: 136,
		AverageRating:   8.5,
	}

	movie2 := &models.Movie{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Title:           "Inception",
		Description:     "A thief who steals corporate secrets through dream-sharing technology",
		ReleaseYear:     2010,
		Director:        "Christopher Nolan",
		DurationMinutes: 148,
		AverageRating:   8.8,
	}

	repo.movies[movie1.ID.String()] = movie1
	repo.movies[movie2.ID.String()] = movie2

	return repo
}

func (r *MockMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.movies[movie.ID.String()] = movie
	return nil
}

func (r *MockMovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if movie, exists := r.movies[id.String()]; exists {
		return movie, nil
	}
	return nil, fmt.Errorf("movie not found")
}

func (r *MockMovieRepository) GetAll(ctx context.Context, filters models.MovieFilters, pagination models.PaginationParams) ([]models.Movie, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var movies []models.Movie
	for _, movie := range r.movies {
		movies = append(movies, *movie)
	}

	return movies, len(movies), nil
}

func (r *MockMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.movies[movie.ID.String()] = movie
	return nil
}

func (r *MockMovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.movies, id.String())
	return nil
}

func (r *MockMovieRepository) UpdateAverageRating(ctx context.Context, movieID uuid.UUID, rating float64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if movie, exists := r.movies[movieID.String()]; exists {
		movie.AverageRating = rating
	}
	return nil
}

func (r *MockMovieRepository) GetGenres(ctx context.Context, movieID uuid.UUID) ([]models.Genre, error) {
	return []models.Genre{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Action"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Sci-Fi"},
	}, nil
}

func (r *MockMovieRepository) SetGenres(ctx context.Context, movieID uuid.UUID, genreIDs []uuid.UUID) error {
	return nil
}

type MockGenreRepository struct {
	genres map[string]*models.Genre
	mutex  sync.RWMutex
}

func NewMockGenreRepository() GenreRepository {
	repo := &MockGenreRepository{
		genres: make(map[string]*models.Genre),
	}

	genres := []*models.Genre{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Action"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Comedy"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Drama"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "Horror"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Name: "Sci-Fi"},
	}

	for _, genre := range genres {
		repo.genres[genre.ID.String()] = genre
	}

	return repo
}

func (r *MockGenreRepository) Create(ctx context.Context, genre *models.Genre) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.genres[genre.ID.String()] = genre
	return nil
}

func (r *MockGenreRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if genre, exists := r.genres[id.String()]; exists {
		return genre, nil
	}
	return nil, fmt.Errorf("genre not found")
}

func (r *MockGenreRepository) GetByName(ctx context.Context, name string) (*models.Genre, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, genre := range r.genres {
		if genre.Name == name {
			return genre, nil
		}
	}
	return nil, fmt.Errorf("genre not found")
}

func (r *MockGenreRepository) GetAll(ctx context.Context) ([]models.Genre, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var genres []models.Genre
	for _, genre := range r.genres {
		genres = append(genres, *genre)
	}

	return genres, nil
}

func (r *MockGenreRepository) Update(ctx context.Context, genre *models.Genre) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.genres[genre.ID.String()] = genre
	return nil
}

func (r *MockGenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.genres, id.String())
	return nil
}

type MockReviewRepository struct {
	reviews map[string]*models.Review
	mutex   sync.RWMutex
}

func NewMockReviewRepository() ReviewRepository {
	return &MockReviewRepository{
		reviews: make(map[string]*models.Review),
	}
}

func (r *MockReviewRepository) Create(ctx context.Context, review *models.Review) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.reviews[review.ID.String()] = review
	return nil
}

func (r *MockReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if review, exists := r.reviews[id.String()]; exists {
		return review, nil
	}
	return nil, fmt.Errorf("review not found")
}

func (r *MockReviewRepository) GetAll(ctx context.Context, pagination models.PaginationParams) ([]models.Review, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var reviews []models.Review
	for _, review := range r.reviews {
		reviews = append(reviews, *review)
	}

	total := len(reviews)
	offset := (pagination.Page - 1) * pagination.Limit
	end := offset + pagination.Limit
	if end > total {
		end = total
	}
	if offset > total {
		return []models.Review{}, total, nil
	}

	return reviews[offset:end], total, nil
}

func (r *MockReviewRepository) GetByMovieID(ctx context.Context, movieID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var reviews []models.Review
	for _, review := range r.reviews {
		if review.MovieID == movieID {
			reviews = append(reviews, *review)
		}
	}

	return reviews, len(reviews), nil
}

func (r *MockReviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var reviews []models.Review
	for _, review := range r.reviews {
		if review.UserID == userID {
			reviews = append(reviews, *review)
		}
	}

	return reviews, len(reviews), nil
}

func (r *MockReviewRepository) GetByMovieAndUser(ctx context.Context, movieID, userID uuid.UUID) (*models.Review, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, review := range r.reviews {
		if review.MovieID == movieID && review.UserID == userID {
			return review, nil
		}
	}
	return nil, fmt.Errorf("review not found")
}

func (r *MockReviewRepository) Update(ctx context.Context, review *models.Review) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.reviews[review.ID.String()] = review
	return nil
}

func (r *MockReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.reviews, id.String())
	return nil
}

func (r *MockReviewRepository) GetAverageRating(ctx context.Context, movieID uuid.UUID) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var total float64
	var count int
	for _, review := range r.reviews {
		if review.MovieID == movieID {
			total += float64(review.Rating)
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return total / float64(count), nil
}
