package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DatabaseEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("an error .env: did not find .env file: %s", err)
	}

	return os.Getenv("DB_URI")
}
