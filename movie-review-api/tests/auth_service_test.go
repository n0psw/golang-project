package tests

import (
	"context"
	"fmt"
	"testing"

	"movie-review-api/internal/models"
	"movie-review-api/internal/service"
	"movie-review-api/pkg/jwt"

	"github.com/google/uuid"
)

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	jwtMgr := jwt.NewManager()
	authService := service.NewAuthService(mockUserRepo, jwtMgr)

	tests := []struct {
		name    string
		req     *models.CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid registration",
			req: &models.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			req: &models.CreateUserRequest{
				Email:    "invalid-email",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short password",
			req: &models.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := authService.Register(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AuthService.Register() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if user == nil {
				t.Errorf("AuthService.Register() user = nil, want user")
				return
			}

			if token == "" {
				t.Errorf("AuthService.Register() token = empty, want token")
				return
			}

			if user.Email != tt.req.Email {
				t.Errorf("AuthService.Register() user.Email = %v, want %v", user.Email, tt.req.Email)
			}

			if user.Username != tt.req.Username {
				t.Errorf("AuthService.Register() user.Username = %v, want %v", user.Username, tt.req.Username)
			}
		})
	}
}

type MockUserRepository struct {
	users map[string]*models.User
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
