package main

import (
	"Gin/initializers"
	"log"
	"os"
)

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Environment not specified. Please provide an environment as the command-line argument [test, dev, prod]. Example: go run autoMigrateDB.go test")
	}
	env := os.Args[1]
	if env != "test" && env != "prod" && env != "dev" {
		log.Fatal("Environment not valid. Please provide an environment as the command-line argument [test, dev, prod]. Example: go run autoMigrateDB.go test")
	}
	log.Printf("Running in environment: %s", env)

	var envDir string
	if env == "test" {
		envDir = "./tests/.env"
	}
	if env == "dev" {
		envDir = "./.env"
	}
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		log.Fatalf("Error: .env file not found in %s", envDir)
	}

	initializers.LoadEnvVariables(envDir)
	log.Println("Environment variables loaded successfully")
	initializers.ConnectToDB()
}

// Run the file selecting the environment like this: go run autoMigrateDB.go test
func main() {
	// initializers.DB.AutoMigrate(
	// &models.User{},
	// &models.Book{},
	// &models.Company{},
	// &models.House{},
	// )
}
