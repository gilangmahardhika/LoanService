package repositories_test

import (
	"testing"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
	"github.com/amartha/LoanService/pkg/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvestmentRepositoryCreate(t *testing.T) {
	// Setup test database
	db := testutils.SetupTestDB(t)

	defer testutils.TruncateTable(t, db)

	// Create repository
	repo := repositories.NewInvestmentRepository(db)

	// Prepare a loan for testing
	loan := &models.Loan{
		BorrowerIDNumber:          "87876",
		Rate:                      10.0,
		RemainingInvestmentAmount: 10000.0,
		PrincipalAmount:           10000.0,
		State:                     "proposed",
	}

	// First, create the loan in the database
	err := db.Create(loan).Error
	require.NoError(t, err, "Failed to create test loan")

	loanID := loan.ID
	loan.State = "approved"
	approvedBy := uint(21)
	loan.ApprovedBy = &approvedBy
	visitProof := "visit_proof"
	loan.VisitProof = &visitProof
	err = db.Save(&loan).Error
	require.NoError(t, err, "Failed to update test loan")

	testCases := []struct {
		name          string
		investment    *models.Investment
		expectedError bool
		errorMessage  string
	}{
		{
			name: "Successful Investment Creation",
			investment: &models.Investment{
				LoanID:         loanID,
				InvestorID:     1,
				InvestedAmount: 5000.0,
			},
			expectedError: false,
		},
		{
			name: "Investment Amount Exceeds Remaining Amount",
			investment: &models.Investment{
				LoanID:         loanID,
				InvestorID:     2,
				InvestedAmount: 15000.0,
			},
			expectedError: true,
			errorMessage:  "invested amount can't be more than remaining investment amount",
		},
		{
			name: "Negative Investment Amount",
			investment: &models.Investment{
				LoanID:         loanID,
				InvestorID:     3,
				InvestedAmount: -1000.0,
			},
			expectedError: true,
			errorMessage:  "invested amount must be more than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to create investment
			err := repo.Create(db, tc.investment)

			// Assert expectations
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}

				// Ensure no investment was created
				var count int64
				err := db.Model(&models.Investment{}).Where("investor_id = ?", tc.investment.InvestorID).Count(&count).Error
				require.NoError(t, err)
				assert.Zero(t, count, "No investment should be created when validation fails")
			} else {
				assert.NoError(t, err)

				// Verify investment was created
				var savedInvestment models.Investment
				err := db.Where("investor_id = ?", tc.investment.InvestorID).First(&savedInvestment).Error
				require.NoError(t, err)
				assert.Equal(t, tc.investment.InvestedAmount, savedInvestment.InvestedAmount)
			}
		})
	}
}
