package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DatabaseEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("an error .env: did not find .env file: %s", err)
	}

	buffer := os.Getenv("DB_URI")

	fmt.Println(buffer)

	return buffer
}
