package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func NewPostgres() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env found")
	}

	dburl := os.Getenv("DATABASE_URL")
	fmt.Println("Loaded db URL", dburl)

	if dburl == "" {
		log.Println("DATABASE_URL is not set in .env")
	}
}