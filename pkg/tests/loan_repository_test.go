package repositories

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"strings"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
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

	// Truncate the table before each test
	truncateLoansTable(t, db)

	return db
}

func truncateLoansTable(t *testing.T, db *gorm.DB) {
	// Truncate the loans table
	err := db.Exec("TRUNCATE TABLE loans RESTART IDENTITY").Error
	require.NoError(t, err, "Failed to truncate loans table")
}

func TestLoanRepository_Create(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

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
	assert.Equal(t, "proposed", testLoan.State, "Loan status should be 'proposed'")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, testLoan.BorrowerIDNumber, savedLoan.BorrowerIDNumber, "Saved loan should match the original loan's borrower ID")
	assert.Equal(t, "proposed", savedLoan.State, "Saved loan should have 'proposed' status")

	// Truncate the table after each test
	truncateLoansTable(t, db)
}

func TestLoanRepository_CreateWithExistingStatus(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

	// Prepare a test loan with an existing status
	testLoan := &models.Loan{
		BorrowerIDNumber: "12345",
		PrincipalAmount:  10000,
		Rate:             5.5,
		State:            "proposed", // Explicitly set to a valid state
	}

	// Execute the Create method
	err := repo.Create(db, testLoan)

	// Assertions
	assert.NoError(t, err, "Creating a loan should not return an error")
	assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")
	assert.Equal(t, "proposed", testLoan.State, "Loan status should be 'proposed'")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, "proposed", savedLoan.State, "Saved loan should have 'proposed' status")

	// Truncate the table after each test
	truncateLoansTable(t, db)
}

func TestLoanRepository_CreateWithDifferentStatus(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

	// Test cases with different initial statuses
	testCases := []models.Loan{
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            "proposed",
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            "approved",
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            "invested",
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            "disbursed",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Create loan with %s status", tc.State), func(t *testing.T) {
			// Prepare a test loan with a non-proposed status
			testLoan := &tc

			// Execute the Create method
			err := repo.Create(db, testLoan)

			// Assertions
			assert.NoError(t, err, "Creating a loan should not return an error")
			assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")
			assert.Equal(t, tc.State, testLoan.State, "Loan status should remain the same")

			// Verify loan was actually saved in the database
			var savedLoan models.Loan
			result := db.First(&savedLoan, testLoan.ID)
			assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
			assert.Equal(t, tc.State, savedLoan.State, "Saved loan should have the original status")

			// Truncate the table after each test
			truncateLoansTable(t, db)
		})
	}
}

func TestLoanRepository_CreateRemainingInvestmentAmount(t *testing.T) {
	// Test cases with different principal amounts
	testCases := []struct {
		name            string
		principalAmount float64
	}{
		{
			name:            "Create loan with small principal amount",
			principalAmount: 1000,
		},
		{
			name:            "Create loan with medium principal amount",
			principalAmount: 10000,
		},
		{
			name:            "Create loan with large principal amount",
			principalAmount: 100000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup in-memory test database
			db := setupTestDB(t)

			// Create a repository instance
			repo := repositories.NewLoanRepository(db)

			// Prepare a test loan
			testLoan := &models.Loan{
				BorrowerIDNumber:          "12345",
				PrincipalAmount:           tc.principalAmount,
				Rate:                      5.5,
				RemainingInvestmentAmount: 0,          // Intentionally set to 0
				State:                     "proposed", // Explicitly set to a valid state
			}

			// Execute the Create method
			err := repo.Create(db, testLoan)

			// Assertions
			assert.NoError(t, err, "Creating a loan should not return an error")
			assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")

			// Check remaining investment amount
			assert.Equal(t, tc.principalAmount, testLoan.RemainingInvestmentAmount,
				"Remaining investment amount should be equal to principal amount")

			// Verify loan was actually saved in the database with correct remaining investment amount
			var savedLoan models.Loan
			result := db.First(&savedLoan, testLoan.ID)
			assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
			assert.Equal(t, tc.principalAmount, savedLoan.RemainingInvestmentAmount,
				"Saved loan's remaining investment amount should be equal to principal amount")

			// Truncate the table after each test
			truncateLoansTable(t, db)
		})
	}
}

func TestLoanRepository_CreateErrorScenarios(t *testing.T) {
	// Setup in-memory test database
	db := setupTestDB(t)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

	testCases := []struct {
		name        string
		loan        *models.Loan
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name: "Create loan with missing required borrower ID",
			loan: &models.Loan{
				BorrowerIDNumber: "", // Empty borrower ID
				PrincipalAmount:  10000,
				Rate:             5.5,
				State:            "proposed", // Explicitly set to a valid state
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return err != nil && strings.Contains(err.Error(), "BorrowerIDNumber")
			},
		},
		{
			name: "Create loan with negative principal amount",
			loan: &models.Loan{
				BorrowerIDNumber: "12345",
				PrincipalAmount:  -1000, // Negative principal amount
				Rate:             5.5,
				State:            "proposed", // Explicitly set to a valid state
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return err != nil && strings.Contains(err.Error(), "PrincipalAmount")
			},
		},
		{
			name: "Create loan with negative rate",
			loan: &models.Loan{
				BorrowerIDNumber: "23456",
				PrincipalAmount:  10000,
				Rate:             -5.5,       // Negative rate
				State:            "proposed", // Explicitly set to a valid state
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return err != nil && strings.Contains(err.Error(), "Rate")
			},
		},
		{
			name: "Create loan with empty state",
			loan: &models.Loan{
				BorrowerIDNumber: "16652",
				PrincipalAmount:  10000,
				Rate:             5.5,
				State:            "", // Empty state
			},
			expectError: false,
		},
		{
			name: "Create valid loan",
			loan: &models.Loan{
				BorrowerIDNumber: "76743",
				PrincipalAmount:  10000,
				Rate:             5.5,
				State:            "proposed", // Explicitly set to a valid state
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the Create method
			err := repo.Create(db, tc.loan)

			if tc.expectError {
				// Check that an error was returned
				assert.Error(t, err, "Expected an error during loan creation")

				// If a specific error check is provided, use it
				if tc.errorCheck != nil {
					assert.True(t, tc.errorCheck(err), "Error did not match expected conditions")
				}

				// Ensure no loan was created
				var count int64
				result := db.Model(&models.Loan{}).Where("borrower_id_number = ?", tc.loan.BorrowerIDNumber).Count(&count)
				assert.NoError(t, result.Error, "Should be able to count loans")
				assert.Zero(t, count, "No loan should be created when validation fails")
				assert.Zero(t, tc.loan.ID, "Loan ID should not be set when creation fails")
			} else {
				// Check that no error was returned
				assert.NoError(t, err, "Creating a valid loan should not return an error")

				// Verify the loan was created
				var savedLoan models.Loan
				result := db.Where("borrower_id_number = ?", tc.loan.BorrowerIDNumber).First(&savedLoan)
				assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
				assert.Equal(t, tc.loan.BorrowerIDNumber, savedLoan.BorrowerIDNumber, "Saved loan should match the original loan")
			}

			// Truncate the table after each test
			truncateLoansTable(t, db)
		})
	}
}
