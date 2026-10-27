package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CollegeBTK-Projects/my-node-go-app/db"
	"github.com/CollegeBTK-Projects/my-node-go-app/handlers"
	"github.com/CollegeBTK-Projects/my-node-go-app/repository"
	"github.com/CollegeBTK-Projects/my-node-go-app/validators"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла: %v", err)
	}

	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	clientMongo, err := db.Connect(ctx, os.Getenv("MONGO_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения к БД: %v", err)
		return
	}

	coll := clientMongo.Database("478").Collection("users")

	validators.Register()
	userRepo := repository.NewUserRepository(coll)
	userHand := handlers.NewUserHandler(userRepo)

	if err := userRepo.CreateIndexes(ctx); err != nil {
		log.Fatal(err)
	}

	api := r.Group("/api/users")
	{
		api.GET("", userHand.List)
		api.GET("/search", userHand.Search)
		api.GET("/stats", userHand.Stats)
		api.GET("/export", userHand.Export)
		api.GET("/:id", userHand.FindById)
		api.POST("", userHand.InsertUser)
		api.PUT("/:id", userHand.UpdateById)
		api.DELETE("/:id", userHand.DeleteById)
	}

	if err := r.Run(":3000"); err != nil {
		log.Fatal("Сервер упал")
	}
}
