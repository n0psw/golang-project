package service

import (
	"context"
	"fmt"

	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"
	"movie-review-api/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, string, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type MovieService interface {
	Create(ctx context.Context, req *models.CreateMovieRequest) (*models.Movie, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Movie, error)
	GetAll(ctx context.Context, filters models.MovieFilters, pagination models.PaginationParams) ([]models.Movie, int, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateMovieRequest) (*models.Movie, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateAverageRating(ctx context.Context, movieID uuid.UUID) error
}

type GenreService interface {
	Create(ctx context.Context, req *models.CreateGenreRequest) (*models.Genre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Genre, error)
	GetAll(ctx context.Context) ([]models.Genre, error)
	Update(ctx context.Context, id uuid.UUID, req *models.CreateGenreRequest) (*models.Genre, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReviewService interface {
	Create(ctx context.Context, movieID uuid.UUID, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetAll(ctx context.Context, pagination models.PaginationParams) ([]models.Review, int, error)
	GetByMovieID(ctx context.Context, movieID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.PaginationParams) ([]models.Review, int, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateReviewRequest) (*models.Review, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type authService struct {
	userRepo repository.UserRepository
	jwtMgr   *jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwtMgr *jwt.Manager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtMgr:   jwtMgr,
	}
}

func (s *authService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, string, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", req.Email)
	}

	existingUser, _ = s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, "", fmt.Errorf("username %s is already taken", req.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.jwtMgr.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	cleanUser := *user
	cleanUser.PasswordHash = ""
	return &cleanUser, token, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtMgr.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}
