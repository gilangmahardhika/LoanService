package models

import (
	"testing"

	"github.com/amartha/LoanService/pkg/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBeforeSave(t *testing.T) {
	testCases := []struct {
		name          string
		state         string
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "Valid state: proposed",
			state:         "proposed",
			expectedError: false,
		},
		{
			name:          "Valid state: approved",
			state:         "approved",
			expectedError: false,
		},
		{
			name:          "Valid state: invested",
			state:         "invested",
			expectedError: false,
		},
		{
			name:          "Valid state: disbursed",
			state:         "disbursed",
			expectedError: false,
		},
		{
			name:          "Invalid state: random string",
			state:         "random",
			expectedError: true,
			errorMessage:  "invalid state 'random', must be one of: proposed, approved, invested, disbursed",
		},
		{
			name:          "Invalid state: empty string",
			state:         "",
			expectedError: true,
			errorMessage:  "invalid state '', must be one of: proposed, approved, invested, disbursed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock gorm.DB (nil is okay for this test since we're not using it)
			var mockTx *gorm.DB

			loan := &models.Loan{
				State: tc.state,
			}

			err := loan.BeforeSave(mockTx)

			if tc.expectedError {
				assert.Error(t, err, "Expected an error for invalid state")
				if err != nil {
					assert.EqualError(t, err, tc.errorMessage, "Error message should match expected")
				}
			} else {
				assert.NoError(t, err, "Expected no error for valid state")
			}
		})
	}
}

func TestCalculateRemainingInvestmentAmount(t *testing.T) {
	testCases := []struct {
		name              string
		principalAmount   float64
		investments       []models.Investment
		expectedRemaining float64
	}{
		{
			name:              "No investments",
			principalAmount:   10000,
			investments:       []models.Investment{},
			expectedRemaining: 10000,
		},
		{
			name:            "Partial investment",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 3000},
				{InvestedAmount: 2000},
			},
			expectedRemaining: 5000,
		},
		{
			name:            "Full investment",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 5000},
				{InvestedAmount: 5000},
			},
			expectedRemaining: 0,
		},
		{
			name:            "Over investment (edge case)",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 6000},
				{InvestedAmount: 5000},
			},
			expectedRemaining: -1000,
		},
		{
			name:            "Fractional investment",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 3333.33},
				{InvestedAmount: 3333.33},
			},
			expectedRemaining: 3333.34, // Due to floating-point arithmetic
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loan := models.Loan{
				PrincipalAmount: tc.principalAmount,
				Investments:     tc.investments,
			}

			remaining := loan.CalculateRemainingInvestmentAmount()

			// Use approximate comparison for floating-point values
			assert.InDelta(t, tc.expectedRemaining, remaining, 0.01,
				"Remaining investment amount should match expected value")
		})
	}
}

func TestUpdateRemainingInvestmentAmount(t *testing.T) {
	testCases := []struct {
		name              string
		principalAmount   float64
		investments       []models.Investment
		expectedRemaining float64
	}{
		{
			name:              "No investments",
			principalAmount:   10000,
			investments:       []models.Investment{},
			expectedRemaining: 10000,
		},
		{
			name:            "Partial investment",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 3000},
				{InvestedAmount: 2000},
			},
			expectedRemaining: 5000,
		},
		{
			name:            "Full investment",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 5000},
				{InvestedAmount: 5000},
			},
			expectedRemaining: 0,
		},
		{
			name:            "Over investment (edge case)",
			principalAmount: 10000,
			investments: []models.Investment{
				{InvestedAmount: 6000},
				{InvestedAmount: 5000},
			},
			expectedRemaining: -1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loan := models.Loan{
				PrincipalAmount:           tc.principalAmount,
				Investments:               tc.investments,
				RemainingInvestmentAmount: 0, // Initialize to ensure it's updated
			}

			// Call the method to update remaining investment amount
			loan.UpdateRemainingInvestmentAmount()

			// Use approximate comparison for floating-point values
			assert.InDelta(t, tc.expectedRemaining, loan.RemainingInvestmentAmount, 0.01,
				"Remaining investment amount should be updated correctly")
		})
	}
}

func TestSetStatusToInvested(t *testing.T) {
	testCases := []struct {
		name                      string
		initialState              string
		remainingInvestmentAmount float64
		expectedState             string
	}{
		{
			name:                      "Fully invested loan changes state",
			initialState:              "proposed",
			remainingInvestmentAmount: 0,
			expectedState:             "invested",
		},
		{
			name:                      "Partially invested loan does not change state",
			initialState:              "proposed",
			remainingInvestmentAmount: 1000,
			expectedState:             "proposed",
		},
		{
			name:                      "Already invested loan remains invested",
			initialState:              "invested",
			remainingInvestmentAmount: 0,
			expectedState:             "invested",
		},
		{
			name:                      "Loan with negative remaining amount does not change state",
			initialState:              "proposed",
			remainingInvestmentAmount: -500,
			expectedState:             "proposed",
		},
		{
			name:                      "Loan with zero point zero remaining amount changes state",
			initialState:              "proposed",
			remainingInvestmentAmount: 0.0,
			expectedState:             "invested",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loan := models.Loan{
				State:                     tc.initialState,
				RemainingInvestmentAmount: tc.remainingInvestmentAmount,
			}

			// Call the method to potentially change state
			loan.SetStatusToInvested()

			// Assert the final state
			assert.Equal(t, tc.expectedState, loan.State,
				"Loan state should be updated correctly based on remaining investment amount")
		})
	}
}

func TestSetStateToProposed(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  string
		expectedState string
	}{
		{
			name:          "Set state to proposed from empty state",
			initialState:  "",
			expectedState: "proposed",
		},
		{
			name:          "Set state to proposed from different state",
			initialState:  "approved",
			expectedState: "proposed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loan := &models.Loan{
				State: tc.initialState,
			}

			loan.SetStateToProposed()

			assert.Equal(t, tc.expectedState, loan.State, "Loan state should be set to proposed")
		})
	}
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
