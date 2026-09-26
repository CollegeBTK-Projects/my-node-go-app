package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CollegeBTK-Projects/my-node-go-app/handlers"
	"github.com/CollegeBTK-Projects/my-node-go-app/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", middleware.Logging(handlers.Handler))
	mux.HandleFunc("GET /api/users", middleware.Logging(handlers.GetAllUsers))
	mux.HandleFunc("POST /api/users", middleware.Logging(handlers.PostUser))
	mux.HandleFunc("PUT /api/users/{id}", middleware.Logging(handlers.PutUser))
	mux.HandleFunc("DELETE /api/users/{id}", middleware.Logging(handlers.DeleteUser))
	mux.HandleFunc("GET /error", middleware.Logging(handlers.ErrorHandler))
	mux.HandleFunc("GET /protected", middleware.AuthHandler(handlers.Auth))
	fmt.Println("Сервер запущен: localhost:3000")
	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
