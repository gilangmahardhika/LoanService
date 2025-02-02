package repositories_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
	testutils "github.com/amartha/LoanService/pkg/tests/testutils"
)

func TestLoanRepository_Create(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)

	defer testutils.TruncateTable(t, db)

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
	assert.Equal(t, models.LoanStatusProposed, testLoan.State, "Loan status should be 'proposed'")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, testLoan.BorrowerIDNumber, savedLoan.BorrowerIDNumber, "Saved loan should match the original loan's borrower ID")
	assert.Equal(t, models.LoanStatusProposed, savedLoan.State, "Saved loan should have 'proposed' status")
}

func TestLoanRepository_CreateWithExistingStatus(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)

	defer testutils.TruncateTable(t, db)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

	// Prepare a test loan with an existing status
	testLoan := &models.Loan{
		BorrowerIDNumber: "12345",
		PrincipalAmount:  10000,
		Rate:             5.5,
		State:            models.LoanStatusProposed, // Explicitly set to a valid state
	}

	// Execute the Create method
	err := repo.Create(db, testLoan)

	// Assertions
	assert.NoError(t, err, "Creating a loan should not return an error")
	assert.NotZero(t, testLoan.ID, "Loan ID should be set after creation")
	assert.Equal(t, models.LoanStatusProposed, testLoan.State, "Loan status should be 'proposed'")

	// Verify loan was actually saved in the database
	var savedLoan models.Loan
	result := db.First(&savedLoan, testLoan.ID)
	assert.NoError(t, result.Error, "Should be able to retrieve the saved loan")
	assert.Equal(t, models.LoanStatusProposed, savedLoan.State, "Saved loan should have 'proposed' status")
}

func TestLoanRepository_CreateWithDifferentStatus(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)

	// Truncate the table after each test
	defer testutils.TruncateTable(t, db)

	// Create a repository instance
	repo := repositories.NewLoanRepository(db)

	// Test cases with different initial statuses
	testCases := []models.Loan{
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            models.LoanStatusProposed,
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            models.LoanStatusApproved,
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            models.LoanStatusInvested,
		},
		{
			BorrowerIDNumber: "12345",
			PrincipalAmount:  10000,
			Rate:             5.5,
			State:            models.LoanStatusDisbursed,
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
			testutils.TruncateTable(t, db)
		})
	}
}

func TestLoanRepository_CreateRemainingInvestmentAmount(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)
	// Truncate the table after each test
	defer testutils.TruncateTable(t, db)
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

			// Create a repository instance
			repo := repositories.NewLoanRepository(db)

			// Prepare a test loan
			testLoan := &models.Loan{
				BorrowerIDNumber:          "12345",
				PrincipalAmount:           tc.principalAmount,
				Rate:                      5.5,
				RemainingInvestmentAmount: 0,                         // Intentionally set to 0
				State:                     models.LoanStatusProposed, // Explicitly set to a valid state
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
		})
	}
}

func TestLoanRepository_CreateErrorScenarios(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)
	// Truncate the table after each test
	defer testutils.TruncateTable(t, db)

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
				State:            models.LoanStatusProposed, // Explicitly set to a valid state
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
				State:            models.LoanStatusProposed, // Explicitly set to a valid state
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
				Rate:             -5.5,                      // Negative rate
				State:            models.LoanStatusProposed, // Explicitly set to a valid state
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
				State:            models.LoanStatusProposed, // Explicitly set to a valid state
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
		})
	}
}

func TestLoanRepository_SetStateToApproved(t *testing.T) {
	// Setup test database
	db := testutils.SetupTestDB(t)
	defer testutils.TruncateTable(t, db)

	// Create a loan repository
	loanRepo := repositories.NewLoanRepository(db)

	// Test case 1: Successfully approve a proposed loan
	t.Run("Approve Proposed Loan", func(t *testing.T) {
		// Create a proposed loan
		loan := &models.Loan{
			BorrowerIDNumber: "32123",
			PrincipalAmount:  10000.0,
			State:            models.LoanStatusProposed,
			Rate:             5.5,
		}
		err := loanRepo.Create(db, loan)
		require.NoError(t, err)

		// Approve the loan
		approverID := uint(1)
		visitProof := "visit_proof.jpg"
		err = loanRepo.SetStateToApproved(db, loan.ID, approverID, visitProof)
		require.NoError(t, err)

		// Retrieve the updated loan
		updatedLoan := &models.Loan{}
		err = db.Where("id = ?", loan.ID).First(updatedLoan).Error
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, models.LoanStatusApproved, updatedLoan.State)
		assert.Equal(t, &approverID, updatedLoan.ApprovedBy)
		assert.Equal(t, &visitProof, updatedLoan.VisitProof)
		assert.NotNil(t, updatedLoan.ApprovedAt)
	})

	// Test case 2: Attempt to approve an already approved loan
	t.Run("Approve Already Approved Loan", func(t *testing.T) {
		// Create an already approved loan
		loan := &models.Loan{
			BorrowerIDNumber: "67890",
			PrincipalAmount:  20000.0,
			ApprovedBy:       &[]uint{1}[0],
			VisitProof:       &[]string{"visit_proof.jpg"}[0],
			ApprovedAt:       &[]time.Time{time.Now()}[0],
			Rate:             5.5,
		}
		err := loanRepo.Create(db, loan)
		require.NoError(t, err)

		// Set loan state to approved
		loan.State = models.LoanStatusApproved
		err = db.Save(loan).Error
		require.NoError(t, err)

		// Attempt to approve the already approved loan
		approverID := uint(2)
		visitProof := "another_proof.jpg"
		err = loanRepo.SetStateToApproved(db, loan.ID, approverID, visitProof)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not in proposed state")
	})

	// Test case 3: Attempt to approve a non-existent loan
	t.Run("Approve Non-Existent Loan", func(t *testing.T) {
		approverID := uint(3)
		visitProof := "non_existent_proof.jpg"
		err := loanRepo.SetStateToApproved(db, 9999, approverID, visitProof)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})
}
