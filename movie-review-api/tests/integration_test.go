package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"movie-review-api/internal/database"
	"movie-review-api/internal/handler"
	"movie-review-api/internal/middleware"
	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"
	"movie-review-api/internal/service"
	"movie-review-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_Register(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	jwtMgr := jwt.NewManager()
	authService := service.NewAuthService(userRepo, jwtMgr)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/auth/register", authHandler.Register)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid registration",
			requestBody: models.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			requestBody: models.CreateUserRequest{
				Email:    "invalid-email",
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			requestBody: models.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("AuthHandler.Register() status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if _, exists := response["user"]; !exists {
					t.Errorf("Response missing 'user' field")
				}

				if _, exists := response["token"]; !exists {
					t.Errorf("Response missing 'token' field")
				}
			}
		})
	}
}

func TestMovieHandler_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieService := service.NewMovieService(movieRepo, genreRepo)
	movieHandler := handler.NewMovieHandler(movieService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/movies", movieHandler.GetAll)

	req := httptest.NewRequest("GET", "/movies", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("MovieHandler.GetAll() status = %v, want %v", rr.Code, http.StatusOK)
	}

	var response models.PaginatedResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Data == nil {
		t.Errorf("Response missing 'data' field")
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	config := &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "movie_review_test_db",
		SSLMode:  "disable",
	}

	db, err := database.Connect(config)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to test database: %v", err)
	}

	return db
}
