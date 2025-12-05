package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-review-api/internal/database"
	"movie-review-api/internal/handler"
	"movie-review-api/internal/middleware"
	"movie-review-api/internal/repository"
	"movie-review-api/internal/service"
	"movie-review-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using environment variables")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	jwtMgr := jwt.NewManager()

	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	authService := service.NewAuthService(userRepo, jwtMgr)
	movieService := service.NewMovieService(movieRepo, genreRepo)
	genreService := service.NewGenreService(genreRepo)

	ratingWorker := service.NewRatingWorker(reviewRepo, movieRepo, 3)
	ratingWorker.Start()
	defer ratingWorker.Stop()

	reviewService := service.NewReviewService(reviewRepo, movieRepo, userRepo, ratingWorker)

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

func setupRoutes(authHandler *handler.AuthHandler, movieHandler *handler.MovieHandler, genreHandler *handler.GenreHandler, reviewHandler *handler.ReviewHandler, jwtMgr *jwt.Manager) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtMgr))
	protected.GET("/users/me", authHandler.GetProfile)
	protected.GET("/users/me/reviews", reviewHandler.GetByUserID)

	movies := api.Group("/movies")
	movies.GET("", movieHandler.GetAll)
	movies.GET("/:id", movieHandler.GetByID)
	movies.GET("/:id/reviews", reviewHandler.GetByMovieID)

	adminMovies := movies.Group("")
	adminMovies.Use(middleware.AuthMiddleware(jwtMgr))
	adminMovies.Use(middleware.AdminMiddleware())
	adminMovies.POST("", movieHandler.Create)
	adminMovies.PUT("/:id", movieHandler.Update)
	adminMovies.DELETE("/:id", movieHandler.Delete)

	genres := api.Group("/genres")
	genres.GET("", genreHandler.GetAll)
	genres.GET("/:id", genreHandler.GetByID)

	adminGenres := genres.Group("")
	adminGenres.Use(middleware.AuthMiddleware(jwtMgr))
	adminGenres.Use(middleware.AdminMiddleware())
	adminGenres.POST("", genreHandler.Create)
	adminGenres.PUT("/:id", genreHandler.Update)
	adminGenres.DELETE("/:id", genreHandler.Delete)

	reviews := api.Group("/reviews")
	reviews.Use(middleware.AuthMiddleware(jwtMgr))
	reviews.GET("/:id", reviewHandler.GetByID)
	reviews.PUT("/:id", reviewHandler.Update)
	reviews.DELETE("/:id", reviewHandler.Delete)

	adminReviews := reviews.Group("")
	adminReviews.Use(middleware.AdminMiddleware())
	adminReviews.GET("", reviewHandler.GetAll)

	reviewCreate := api.Group("/movies/:id/reviews")
	reviewCreate.Use(middleware.AuthMiddleware(jwtMgr))
	reviewCreate.POST("", reviewHandler.Create)

	return router
}
