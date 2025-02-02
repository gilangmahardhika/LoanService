// Investment represents the investment details for a loan
package models

import (
	"fmt"
	"time"
)

type Investment struct {
	ID             uint      `gorm:"not null;primary_key" json:"id"`
	LoanID         uint      `gorm:"not null" json:"loan_id"`
	Loan           Loan      `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
	InvestorID     string    `gorm:"not null" json:"investor_id"`
	InvestedAmount float64   `gorm:"not null" json:"invested_amount"`
	ROI            float64   `gorm:"not null" json:"roi"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

// Function for calculating the ROI
func (i Investment) CalculateROI() float64 {
	return i.InvestedAmount * i.Loan.Rate / 100
}

// Function for updating the ROI
func (i Investment) UpdateROI() {
	i.ROI = i.CalculateROI()
}

// Validate invested amount should be more than 0 and can't be more than remaining investment amount
func (i Investment) ValidateInvestedAmount() error {
	if i.InvestedAmount <= 0 {
		return fmt.Errorf("invested amount must be more than 0")
	}
	if i.InvestedAmount > i.Loan.CalculateRemainingInvestmentAmount() {
		return fmt.Errorf("invested amount can't be more than remaining investment amount")
	}
	return nil
}
