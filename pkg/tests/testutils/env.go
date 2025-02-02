package testutils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// GetDatabaseURL retrieves the database URL from environment variables
// It attempts to load a .env file and provides a default connection URL if not set
func GetDatabaseURL(envVarName string, defaultURL string) string {
	// Load environment variables from .env file
	projectRoot, err := os.Getwd()
	if err != nil {
		projectRoot = "."
	}
	projectRoot = filepath.Join(projectRoot, "..", "..", "..")

	// Load .env file (ignore errors if file doesn't exist)
	_ = godotenv.Load(filepath.Join(projectRoot, ".env"))

	// Retrieve PostgreSQL connection URL from environment variable
	databaseURL := os.Getenv(envVarName)
	if databaseURL == "" {
		databaseURL = defaultURL
	}

	return databaseURL
}

// GetTestDatabaseURL is a convenience function for getting the test database URL
func GetTestDatabaseURL() string {
	return GetDatabaseURL(
		"TEST_DATABASE_URL",
		"postgres://postgres:postgres@localhost:5432/loanservice_test?sslmode=disable&TimeZone=UTC",
	)
}
