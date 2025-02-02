package models

import (
	"testing"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateROI(t *testing.T) {
	testCases := []struct {
		name           string
		investedAmount float64
		loanRate       float64
		expectedROI    float64
	}{
		{
			name:           "Basic ROI calculation",
			investedAmount: 10000,
			loanRate:       10,
			expectedROI:    1000,
		},
		{
			name:           "Zero invested amount",
			investedAmount: 0,
			loanRate:       15,
			expectedROI:    0,
		},
		{
			name:           "Zero loan rate",
			investedAmount: 5000,
			loanRate:       0,
			expectedROI:    0,
		},
		{
			name:           "Fractional invested amount",
			investedAmount: 7500.50,
			loanRate:       12.5,
			expectedROI:    937.5625,
		},
		{
			name:           "Negative loan rate",
			investedAmount: 10000,
			loanRate:       -5,
			expectedROI:    -500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock Loan with the specified rate
			mockLoan := models.Loan{
				Rate: tc.loanRate,
			}

			// Create an Investment with the invested amount and mock loan
			investment := models.Investment{
				InvestedAmount: tc.investedAmount,
				Loan:           mockLoan,
			}

			// Calculate ROI
			roi := investment.CalculateROI()

			// Use InDelta for floating-point comparison
			assert.InDelta(t, tc.expectedROI, roi, 0.0001,
				"ROI calculation should be accurate")
		})
	}
}

func TestUpdateROI(t *testing.T) {
	testCases := []struct {
		name           string
		investedAmount float64
		loanRate       float64
		expectedROI    float64
	}{
		{
			name:           "Basic ROI update",
			investedAmount: 10000,
			loanRate:       10,
			expectedROI:    1000,
		},
		{
			name:           "Zero invested amount",
			investedAmount: 0,
			loanRate:       15,
			expectedROI:    0,
		},
		{
			name:           "Zero loan rate",
			investedAmount: 5000,
			loanRate:       0,
			expectedROI:    0,
		},
		{
			name:           "Fractional invested amount",
			investedAmount: 7500.50,
			loanRate:       12.5,
			expectedROI:    937.5625,
		},
		{
			name:           "Negative loan rate",
			investedAmount: 10000,
			loanRate:       -5,
			expectedROI:    -500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock Loan with the specified rate
			mockLoan := models.Loan{
				Rate: tc.loanRate,
			}

			// Create an Investment with the invested amount and mock loan
			investment := &models.Investment{
				InvestedAmount: tc.investedAmount,
				Loan:           mockLoan,
				ROI:            0, // Initialize ROI to 0 to ensure it gets updated
			}

			// Update ROI
			investment.UpdateROI()

			// Use InDelta for floating-point comparison
			assert.InDelta(t, tc.expectedROI, investment.ROI, 0.0001,
				"ROI should be updated correctly")
		})
	}
}

func TestValidateInvestedAmount(t *testing.T) {
	testCases := []struct {
		name               string
		principalAmount    float64
		investedAmount     float64
		existingInvestment float64
		expectedError      bool
		errorMessage       string
	}{
		{
			name:               "Valid investment amount",
			principalAmount:    10000,
			investedAmount:     5000,
			existingInvestment: 2000,
			expectedError:      false,
		},
		{
			name:               "Zero invested amount",
			principalAmount:    10000,
			investedAmount:     0,
			existingInvestment: 0,
			expectedError:      true,
			errorMessage:       "invested amount must be more than 0",
		},
		{
			name:               "Negative invested amount",
			principalAmount:    10000,
			investedAmount:     -1000,
			existingInvestment: 0,
			expectedError:      true,
			errorMessage:       "invested amount must be more than 0",
		},
		{
			name:               "Investment exceeding remaining amount",
			principalAmount:    10000,
			investedAmount:     7000,
			existingInvestment: 4000,
			expectedError:      true,
			errorMessage:       "invested amount can't be more than remaining investment amount",
		},
		{
			name:               "Investment equal to remaining amount",
			principalAmount:    10000,
			investedAmount:     3000,
			existingInvestment: 7000,
			expectedError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock Loan with the specified principal amount
			mockLoan := models.Loan{
				PrincipalAmount: tc.principalAmount,
				Investments: []models.Investment{
					{InvestedAmount: tc.existingInvestment},
				},
			}

			// Create an Investment with the invested amount and mock loan
			investment := models.Investment{
				InvestedAmount: tc.investedAmount,
				Loan:           mockLoan,
			}

			// Validate invested amount
			err := investment.ValidateInvestedAmount()

			if tc.expectedError {
				assert.Error(t, err, "Expected an error for invalid investment amount")
				if err != nil {
					assert.EqualError(t, err, tc.errorMessage, "Error message should match expected")
				}
			} else {
				assert.NoError(t, err, "Expected no error for valid investment amount")
			}
		})
	}
}
