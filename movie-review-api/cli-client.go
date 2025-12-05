package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const API_BASE = "http://localhost:8080/api/v1"

type Client struct {
	token string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) makeRequest(method, url string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, API_BASE+url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return map[string]interface{}{"data": result}, nil
}

func (c *Client) Register(email, username, password string) error {
	data := map[string]string{
		"email":    email,
		"username": username,
		"password": password,
	}

	result, err := c.makeRequest("POST", "/auth/register", data)
	if err != nil {
		return err
	}

	if token, ok := result["token"].(string); ok {
		c.SetToken(token)
		fmt.Println("✅ Регистрация успешна!")
		fmt.Printf("Токен: %s\n", token)
	} else {
		fmt.Println("❌ Ошибка получения токена")
	}

	return nil
}

func (c *Client) Login(email, password string) error {
	data := map[string]string{
		"email":    email,
		"password": password,
	}

	result, err := c.makeRequest("POST", "/auth/login", data)
	if err != nil {
		return err
	}

	if token, ok := result["token"].(string); ok {
		c.SetToken(token)
		fmt.Println("✅ Логин успешен!")
		fmt.Printf("Токен: %s\n", token)
	} else {
		fmt.Println("❌ Ошибка получения токена")
	}

	return nil
}

func (c *Client) GetProfile() error {
	result, err := c.makeRequest("GET", "/users/me", nil)
	if err != nil {
		return err
	}

	fmt.Println("👤 Профиль пользователя:")
	printJSON(result)
	return nil
}

func (c *Client) GetMovies() error {
	result, err := c.makeRequest("GET", "/movies", nil)
	if err != nil {
		return err
	}

	fmt.Println("🎬 Фильмы:")
	printJSON(result)
	return nil
}

func (c *Client) GetGenres() error {
	result, err := c.makeRequest("GET", "/genres", nil)
	if err != nil {
		return err
	}

	fmt.Println("🎭 Жанры:")
	printJSON(result)
	return nil
}

func (c *Client) CreateReview(movieID string, rating int, title, content string) error {
	data := map[string]interface{}{
		"rating":  rating,
		"title":   title,
		"content": content,
	}

	result, err := c.makeRequest("POST", "/movies/"+movieID+"/reviews", data)
	if err != nil {
		return err
	}

	fmt.Println("⭐ Отзыв создан:")
	printJSON(result)
	return nil
}

func (c *Client) GetMyReviews() error {
	result, err := c.makeRequest("GET", "/users/me/reviews", nil)
	if err != nil {
		return err
	}

	fmt.Println("📝 Мои отзывы:")
	printJSON(result)
	return nil
}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

func main() {
	client := NewClient()

	if token := os.Getenv("AUTH_TOKEN"); token != "" {
		client.SetToken(token)
	}

	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  go run cli-client.go login <email> <password>")
		fmt.Println("  go run cli-client.go register <email> <username> <password>")
		fmt.Println("  go run cli-client.go profile")
		fmt.Println("  go run cli-client.go movies")
		fmt.Println("  go run cli-client.go genres")
		fmt.Println("  go run cli-client.go review <movie_id> <rating> <title> <content>")
		fmt.Println("  go run cli-client.go my-reviews")
		return
	}

	command := os.Args[1]

	switch command {
	case "login":
		if len(os.Args) < 4 {
			fmt.Println("Использование: login <email> <password>")
			return
		}
		err := client.Login(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Printf("❌ Ошибка логина: %v\n", err)
		}

	case "register":
		if len(os.Args) < 5 {
			fmt.Println("Использование: register <email> <username> <password>")
			return
		}
		err := client.Register(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("❌ Ошибка регистрации: %v\n", err)
		}

	case "profile":
		err := client.GetProfile()
		if err != nil {
			fmt.Printf("❌ Ошибка получения профиля: %v\n", err)
		}

	case "movies":
		err := client.GetMovies()
		if err != nil {
			fmt.Printf("❌ Ошибка получения фильмов: %v\n", err)
		}

	case "genres":
		err := client.GetGenres()
		if err != nil {
			fmt.Printf("❌ Ошибка получения жанров: %v\n", err)
		}

	case "review":
		if len(os.Args) < 6 {
			fmt.Println("Использование: review <movie_id> <rating> <title> <content>")
			return
		}
		err := client.CreateReview(os.Args[2], 8, os.Args[4], os.Args[5])
		if err != nil {
			fmt.Printf("❌ Ошибка создания отзыва: %v\n", err)
		}

	case "my-reviews":
		err := client.GetMyReviews()
		if err != nil {
			fmt.Printf("❌ Ошибка получения отзывов: %v\n", err)
		}

	default:
		fmt.Printf("❌ Неизвестная команда: %s\n", command)
	}
}
