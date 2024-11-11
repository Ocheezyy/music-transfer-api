package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// Load environment variables so they are accessible in os.GetEnv()
func LoadEnvs() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env: %s", err)
	}
}
