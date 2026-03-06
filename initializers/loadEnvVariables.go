package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(envDir string) {
	err := godotenv.Load(envDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("No .env file found at %q; using environment variables", envDir)
			return
		}
		log.Fatalf("Error loading .env file at %q: %v", envDir, err)
	}
}
