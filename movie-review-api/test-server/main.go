package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	log.Println("Test client server starting on :3000")
	log.Println("Open http://localhost:3000/test-client.html in your browser")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

