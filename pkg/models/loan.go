package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Loan represents a loan application in the system
type Loan struct {
	ID                        uint         `gorm:"not null;primary_key" json:"id"`
	BorrowerIDNumber          string       `gorm:"not null" json:"borrower_id_number"`
	PrincipalAmount           float64      `gorm:"not null" json:"principal_amount"`
	Rate                      float64      `gorm:"not null" json:"rate"`
	RemainingInvestmentAmount float64      `gorm:"not null" json:"remaining_investment_amount"`
	AgreementLink             string       `gorm:"not null" json:"agreement_link"`
	State                     string       `gorm:"not null;default:'proposed'" json:"state"`
	Investments               []Investment `gorm:"foreignKey:LoanID" json:"investments,omitempty"`
	ApprovedAt                *time.Time   `gorm:"not null;default:null" json:"approved_at"`
	ApprovedBy                *string      `gorm:"not null;default:null" json:"approved_by"`
	DisbursedAt               *time.Time   `gorm:"not null;default:null" json:"disbursed_at"`
	DisbursedBy               *string      `gorm:"not null;default:null" json:"disbursed_by"`
	CreatedAt                 time.Time    `gorm:"not null" json:"created_at"`
	UpdatedAt                 time.Time    `gorm:"not null" json:"updated_at"`
}

// function for calculation remaining investment amount
func (l Loan) CalculateRemainingInvestmentAmount() float64 {
	// Calculate the sum of invested amounts
	var totalInvestedAmount float64
	for _, investment := range l.Investments {
		totalInvestedAmount += investment.InvestedAmount
	}

	// Calculate the remaining investment amount
	return l.PrincipalAmount - totalInvestedAmount
}

// After the loan is approved, the remaining investment amount will be updated
func (l *Loan) UpdateRemainingInvestmentAmount() {
	l.RemainingInvestmentAmount = l.CalculateRemainingInvestmentAmount()
}

// Set status to invested if remaining investment amount is 0
func (l *Loan) SetStatusToInvested() {
	if l.RemainingInvestmentAmount == 0 {
		l.State = "invested"
	}
}

// Validate loan state before saving
func (l *Loan) BeforeSave(tx *gorm.DB) error {
	validStates := map[string]bool{
		"proposed":   true,
		"approved":   true,
		"invested":   true,
		"disbursed":  true,
	}

	if !validStates[l.State] {
		return fmt.Errorf(
			"invalid state '%s', must be one of: proposed, approved, invested, disbursed", 
			l.State)
	}
	return nil
}
