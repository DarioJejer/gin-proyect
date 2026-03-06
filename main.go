package main

import (
	_ "Gin/docs"
	"Gin/initializers"
	"Gin/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables("./.env")
	log.Println("Environment variables loaded successfully")
	initializers.ConnectToDB()
	if err := initializers.AutoMigrate(); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
}

// @title Gin API
// @version 1.0
// @description This is a sample server for a Gin application.
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Listen on all interfaces so the server is reachable from the host (e.g. Docker port mapping)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal(err)
	}
}
