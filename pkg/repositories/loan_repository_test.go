package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/testutils"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Get the test database URL
	databaseURL := testutils.GetTestDatabaseURL()

	// Open connection to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database: %v", err)

	// Migrate the schema
	err = db.AutoMigrate(&models.Loan{})
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

func TestLoanRepository_Create(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := NewLoanRepository(db)

	// Prepare a test loan
	testLoan := &models.Loan{
		BorrowerIDNumber: "12345",
		PrincipalAmount:  10000,
		Rate:             5.5,
	}

	// Execute the Create method
	err := repo.Create(db, testLoan)

	// Assertions
	assert.NoError(t, err, "Creating a loan should not return an error")
	assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")
	assert.Equal(t, "proposed", testLoan.State, "Newly created loan should have 'proposed' status")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, testLoan.BorrowerIDNumber, savedLoan.BorrowerIDNumber, "Saved loan should match the original loan's borrower ID")
	assert.Equal(t, "proposed", savedLoan.State, "Saved loan should have 'proposed' status")
}

func TestLoanRepository_CreateWithExistingStatus(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := NewLoanRepository(db)

	// Prepare a test loan with a pre-set status
	testLoan := &models.Loan{
		BorrowerIDNumber: "12345",
		PrincipalAmount:  10000,
		Rate:             5.5,
		State:            "approved", // Explicitly set a different status
	}

	// Execute the Create method
	err := repo.Create(db, testLoan)

	// Assertions
	assert.NoError(t, err, "Creating a loan should not return an error")
	assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")
	assert.Equal(t, "proposed", testLoan.State, "Loan status should be overridden to 'proposed'")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, "proposed", savedLoan.State, "Saved loan should have 'proposed' status")
}
