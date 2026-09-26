package main

import (
	"log"
	"net/http"

	"github.com/CollegeBTK-Projects/my-node-go-app/handlers"
	"github.com/CollegeBTK-Projects/my-node-go-app/middleware"
	"github.com/CollegeBTK-Projects/my-node-go-app/models"
)

func main() {
	store := models.NewBookStore()
	books := handlers.NewBookHandler(store)

	rateLimiter := middleware.NewRateLimiter(100)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /about", handlers.About)
	mux.HandleFunc("GET /api/books", books.GetAll)
	mux.HandleFunc("GET /api/books/search", books.Search)
	mux.HandleFunc("GET /api/books/{id}", books.GetByID)
	mux.HandleFunc("POST /api/books", books.Create)
	mux.HandleFunc("PUT /api/books/{id}", books.Update)
	mux.HandleFunc("DELETE /api/books/{id}", books.Delete)
	mux.HandleFunc("GET /error", handlers.TestError)
	mux.HandleFunc("GET /async-error", handlers.TestAsyncError)

	var handler http.Handler = mux
	handler = middleware.Recoverer(handler)
	handler = middleware.Compression(handler)
	handler = rateLimiter.Middleware(handler)
	handler = middleware.Logger(handler)

	log.Println("Сервер запущен: localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
