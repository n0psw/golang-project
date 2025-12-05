package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"movie-review-api/internal/database"
	"movie-review-api/internal/models"
	"movie-review-api/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var email, username, password string
	flag.StringVar(&email, "email", "", "Admin email (required)")
	flag.StringVar(&username, "username", "", "Admin username (required)")
	flag.StringVar(&password, "password", "", "Admin password (required)")
	flag.Parse()

	if email == "" || username == "" || password == "" {
		fmt.Println("Usage: go run cmd/admin/create_admin.go -email <email> -username <username> -password <password>")
		fmt.Println("Or set environment variables: ADMIN_EMAIL, ADMIN_USERNAME, ADMIN_PASSWORD")
		flag.PrintDefaults()
		os.Exit(1)
	}

	envEmail := os.Getenv("ADMIN_EMAIL")
	envUsername := os.Getenv("ADMIN_USERNAME")
	envPassword := os.Getenv("ADMIN_PASSWORD")

	if envEmail != "" {
		email = envEmail
	}
	if envUsername != "" {
		username = envUsername
	}
	if envPassword != "" {
		password = envPassword
	}

	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	existingUser, _ := userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		log.Fatalf("User with email %s already exists", email)
	}

	existingUser, _ = userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		log.Fatalf("Username %s is already taken", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("Admin user created successfully!\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Role: admin\n")
}
