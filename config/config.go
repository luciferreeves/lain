package config

import (
	"lain/utils/env"
	"log"

	"github.com/joho/godotenv"
)

var (
	Server     server
	MailServer mail
	Database   database
	MinIO      minio
	AIServer   ai
	Session    session
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := env.Parse(&Server); err != nil {
		log.Fatalf("Failed to parse ServerConfig: %v", err)
	}

	if err := env.Parse(&MailServer); err != nil {
		log.Fatalf("Failed to parse MailServerConfig: %v", err)
	}

	if err := env.Parse(&Database); err != nil {
		log.Fatalf("Failed to parse DatabaseConfig: %v", err)
	}

	if err := env.Parse(&MinIO); err != nil {
		log.Fatalf("Failed to parse MinIOConfig: %v", err)
	}

	if err := env.Parse(&AIServer); err != nil {
		log.Fatalf("Failed to parse AIServerConfig: %v", err)
	}

	if err := env.Parse(&Session); err != nil {
		log.Fatalf("Failed to parse SessionConfig: %v", err)
	}
}
