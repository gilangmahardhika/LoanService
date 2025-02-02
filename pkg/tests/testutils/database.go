package repositories_test

import (
	"testing"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	// Get the test database URL
	databaseURL := GetTestDatabaseURL()

	// Open connection to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database: %v", err)

	// Migrate the schema
	err = db.AutoMigrate(&models.Loan{}, &models.Investment{})
	require.NoError(t, err, "Failed to migrate database schema")

	// Truncate the table before each test
	TruncateTable(t, db)

	return db
}

// TruncateTables truncates specified database tables and restarts identity
func TruncateTable(t *testing.T, db *gorm.DB) {
	// Truncate the loans and investments tables
	err := db.Exec("TRUNCATE TABLE investments, loans RESTART IDENTITY").Error
	require.NoError(t, err, "Failed to truncate tables")
}
