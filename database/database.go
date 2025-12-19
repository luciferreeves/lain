package database

import (
	"fmt"
	"lain/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	var err error

	DSN := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Name,
		config.Database.SSLMode,
	)

	if config.Database.Password != "" {
		DSN += fmt.Sprintf(" password=%s", config.Database.Password)
	}

	loglevel := logger.Silent
	if config.Server.DevMode {
		loglevel = logger.Info
	}

	dialector := postgres.Open(DSN)

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err = migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("database connection established")
}
