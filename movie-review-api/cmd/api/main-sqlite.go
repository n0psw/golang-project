package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-review-api/internal/handler"
	"movie-review-api/internal/middleware"
	"movie-review-api/internal/repository"
	"movie-review-api/internal/service"
	"movie-review-api/pkg/jwt"

	"database/sql"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using environment variables")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := initDatabase(db); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	jwtMgr := jwt.NewManager()

	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	authService := service.NewAuthService(userRepo, jwtMgr)
	movieService := service.NewMovieService(movieRepo, genreRepo)
	genreService := service.NewGenreService(genreRepo)
	reviewService := service.NewReviewService(reviewRepo, movieRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	movieHandler := handler.NewMovieHandler(movieService)
	genreHandler := handler.NewGenreHandler(genreService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	router := setupRoutes(authHandler, movieHandler, genreHandler, reviewHandler, jwtMgr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}

func initDatabase(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS genres (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS movies (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			release_year INTEGER NOT NULL,
			director TEXT,
			duration_minutes INTEGER,
			average_rating REAL DEFAULT 0.0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS movie_genres (
			movie_id TEXT NOT NULL,
			genre_id TEXT NOT NULL,
			PRIMARY KEY (movie_id, genre_id),
			FOREIGN KEY (movie_id) REFERENCES movies(id),
			FOREIGN KEY (genre_id) REFERENCES genres(id)
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id TEXT PRIMARY KEY,
			movie_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 10),
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(movie_id, user_id),
			FOREIGN KEY (movie_id) REFERENCES movies(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	seedData := []string{
		`INSERT OR IGNORE INTO genres (id, name) VALUES 
			('genre-1', 'Action'),
			('genre-2', 'Comedy'),
			('genre-3', 'Drama'),
			('genre-4', 'Horror'),
			('genre-5', 'Sci-Fi')`,
		`INSERT OR IGNORE INTO users (id, email, username, password_hash, role) VALUES 
			('admin-1', 'admin@example.com', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')`,
		`INSERT OR IGNORE INTO movies (id, title, description, release_year, director, duration_minutes) VALUES 
			('movie-1', 'The Matrix', 'A computer hacker learns about the true nature of reality', 1999, 'Lana Wachowski', 136),
			('movie-2', 'Inception', 'A thief who steals corporate secrets through dream-sharing technology', 2010, 'Christopher Nolan', 148)`,
	}

	for _, query := range seedData {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func setupRoutes(authHandler *handler.AuthHandler, movieHandler *handler.MovieHandler, genreHandler *handler.GenreHandler, reviewHandler *handler.ReviewHandler, jwtMgr *jwt.Manager) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.CORSMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtMgr))
	protected.HandleFunc("/users/me", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/users/me/reviews", reviewHandler.GetByUserID).Methods("GET")

	movies := api.PathPrefix("/movies").Subrouter()
	movies.HandleFunc("", movieHandler.GetAll).Methods("GET")
	movies.HandleFunc("/{id}", movieHandler.GetByID).Methods("GET")
	movies.HandleFunc("/{id}/reviews", reviewHandler.GetByMovieID).Methods("GET")

	adminMovies := movies.PathPrefix("").Subrouter()
	adminMovies.Use(middleware.AuthMiddleware(jwtMgr))
	adminMovies.Use(middleware.AdminMiddleware)
	adminMovies.HandleFunc("", movieHandler.Create).Methods("POST")
	adminMovies.HandleFunc("/{id}", movieHandler.Update).Methods("PUT")
	adminMovies.HandleFunc("/{id}", movieHandler.Delete).Methods("DELETE")

	genres := api.PathPrefix("/genres").Subrouter()
	genres.HandleFunc("", genreHandler.GetAll).Methods("GET")
	genres.HandleFunc("/{id}", genreHandler.GetByID).Methods("GET")

	adminGenres := genres.PathPrefix("").Subrouter()
	adminGenres.Use(middleware.AuthMiddleware(jwtMgr))
	adminGenres.Use(middleware.AdminMiddleware)
	adminGenres.HandleFunc("", genreHandler.Create).Methods("POST")
	adminGenres.HandleFunc("/{id}", genreHandler.Update).Methods("PUT")
	adminGenres.HandleFunc("/{id}", genreHandler.Delete).Methods("DELETE")

	reviews := api.PathPrefix("/reviews").Subrouter()
	reviews.Use(middleware.AuthMiddleware(jwtMgr))
	reviews.HandleFunc("/{id}", reviewHandler.GetByID).Methods("GET")
	reviews.HandleFunc("/{id}", reviewHandler.Update).Methods("PUT")
	reviews.HandleFunc("/{id}", reviewHandler.Delete).Methods("DELETE")

	reviewCreate := api.PathPrefix("/movies/{id}/reviews").Subrouter()
	reviewCreate.Use(middleware.AuthMiddleware(jwtMgr))
	reviewCreate.HandleFunc("", reviewHandler.Create).Methods("POST")

	return router
}
