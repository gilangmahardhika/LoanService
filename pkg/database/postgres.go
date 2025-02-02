package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDatabase initializes the database connection
func InitDatabase() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		// Get DATABASE_URL from environment
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			err = fmt.Errorf("DATABASE_URL environment variable is not set")
			return
		}

		// Configure logging based on environment
		logLevel := logger.Silent
		if os.Getenv("APP_ENV") == "development" {
			logLevel = logger.Info
		}

		// Open database connection
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})

		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return
		}

		// Optional: Connection pool and other configurations
		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			err = sqlErr
			return
		}

		// Parse CONNECTION_POOL from environment
		connectionPool := 25 // default value
		if poolStr := os.Getenv("CONNECTION_POOL"); poolStr != "" {
			if poolVal, err := strconv.Atoi(poolStr); err == nil {
				connectionPool = poolVal
			} else {
				log.Printf("Invalid CONNECTION_POOL value: %v, using default: %d", err, connectionPool)
			}
		}

		sqlDB.SetMaxOpenConns(connectionPool)
		sqlDB.SetMaxIdleConns(connectionPool)
		// sqlDB.SetConnMaxLifetime(5 * time.Minute)
	})

	return db, err
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return db
}
