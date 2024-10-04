package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/handler"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/infrastructure/cache"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/infrastructure/database"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/middleware"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/repository"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}

	redis, err := cache.NewRedis()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := router.Group("/v1")

	repoUser := repository.NewUser(db, redis)
	useCaseUser := usecase.NewUser(repoUser)
	handlerUser := handler.NewUser(useCaseUser)
	handlerUser.RegisterRoutes(v1.Group("/user"))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server started on port %s\n", os.Getenv("PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
