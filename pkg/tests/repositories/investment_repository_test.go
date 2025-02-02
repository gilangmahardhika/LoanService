package repositories_test

import (
	"log"
	"testing"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
	testutils "github.com/amartha/LoanService/pkg/tests/testutils"
	"gorm.io/gorm"

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
		State:                     models.LoanStatusProposed,
	}

	// First, create the loan in the database
	err := db.Create(loan).Error
	require.NoError(t, err, "Failed to create test loan")

	loanID := loan.ID
	err = prepareLoanForInvestment(db, loan)
	require.NoError(t, err, "Failed to update loan state")

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

// Create test for make status == invested
func TestInvestmentRepository_CreateInvestmentWithLoanStatusInvested(t *testing.T) {
	// Setup in-memory test database
	db := testutils.SetupTestDB(t)
	defer testutils.TruncateTable(t, db)

	// Create repository
	repo := repositories.NewInvestmentRepository(db)

	// Create loan
	loan := &models.Loan{
		BorrowerIDNumber:          "12345",
		PrincipalAmount:           10000,
		Rate:                      5.5,
		RemainingInvestmentAmount: 10000,
		State:                     models.LoanStatusProposed,
	}

	err := db.Create(loan).Error
	require.NoError(t, err, "Failed to create loan")

	err = prepareLoanForInvestment(db, loan)
	require.NoError(t, err, "Failed to update loan state")

	// Create investment
	investments := []models.Investment{
		{
			LoanID:         loan.ID,
			InvestorID:     223,
			InvestedAmount: 5000,
		},
		{
			LoanID:         loan.ID,
			InvestorID:     221,
			InvestedAmount: 5000,
		},
	}

	for _, investment := range investments {
		err := repo.Create(db, &investment)
		require.NoError(t, err, "Failed to create investment")
	}

	// Update loan status to invested
	// fetch updated loan
	reloadedLoan := &models.Loan{}
	err = db.Preload("Investments").First(reloadedLoan, loan.ID).Error
	require.NoError(t, err, "Failed to fetch updated loan")

	// Assert
	assert.Equal(t, 2, len(reloadedLoan.Investments))
	assert.Equal(t, models.LoanStatusInvested, reloadedLoan.State)
	assert.Equal(t, float64(0), reloadedLoan.RemainingInvestmentAmount)
}

func prepareLoanForInvestment(db *gorm.DB, loan *models.Loan) error {
	loan.State = models.LoanStatusApproved
	approvedBy := uint(21)
	loan.ApprovedBy = &approvedBy
	visitProof := "visit-proof"
	loan.VisitProof = &visitProof

	log.Println("State is ", loan.State)

	return db.Save(loan).Error
}
