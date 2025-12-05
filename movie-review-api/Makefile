.PHONY: migrate-up migrate-down migrate-create test build run

migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

test:
	go test ./...

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

docker-build:
	docker build -t movie-review-api .

docker-run:
	docker run -p 8080:8080 movie-review-api

