package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// LoanStatus represents the different states a loan can be in
type LoanStatus string

// Predefined loan status constants
const (
	LoanStatusProposed  LoanStatus = "proposed"
	LoanStatusInvested  LoanStatus = "invested"
	LoanStatusApproved  LoanStatus = "approved"
	LoanStatusDisbursed LoanStatus = "disbursed"
)

// Loan represents a loan application in the system
type Loan struct {
	ID                        uint         `gorm:"not null;primary_key" json:"id"`
	BorrowerIDNumber          string       `gorm:"not null" json:"borrower_id_number"`
	PrincipalAmount           float64      `gorm:"not null" json:"principal_amount"`
	Rate                      float64      `gorm:"not null;" json:"rate"`
	RemainingInvestmentAmount float64      `gorm:"not null;default:0" json:"remaining_investment_amount"`
	State                     LoanStatus   `gorm:"not null;default:'proposed'" json:"state"`
	Investments               []Investment `gorm:"foreignKey:LoanID" json:"investments,omitempty"`
	ApprovedAt                *time.Time   `gorm:"default:null" json:"approved_at"`
	ApprovedBy                *uint        `gorm:"default:null" json:"approved_by"`
	VisitProof                *string      `gorm:"default:null" json:"visit_proof"`
	DisbursedAt               *time.Time   `gorm:"default:null" json:"disbursed_at"`
	DisbursedBy               *uint        `gorm:"default:null" json:"disbursed_by"`
	CreatedAt                 time.Time    `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time    `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

// Set default state to proposed before creating loan
func (l *Loan) SetStateToProposed() {
	l.State = LoanStatusProposed
}

// Set default remaining investment amount to principal amount
func (l *Loan) SetDefaultRemainingInvestmentAmount() {
	l.RemainingInvestmentAmount = l.PrincipalAmount
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
	l.State = LoanStatusInvested
}

// Validate loan state before saving
func (l *Loan) BeforeSave(tx *gorm.DB) error {
	validStates := map[LoanStatus]bool{
		LoanStatusProposed:  true,
		LoanStatusApproved:  true,
		LoanStatusInvested:  true,
		LoanStatusDisbursed: true,
	}

	if !validStates[l.State] {
		return fmt.Errorf(
			"invalid state '%s', must be one of: proposed, approved, invested, disbursed",
			l.State)
	}
	return nil
}

// Reduce Remaining Investment Amount
func (l *Loan) ReduceRemainingInvestmentAmount(amount float64) {
	l.RemainingInvestmentAmount -= amount
}

// BeforeCreate is a GORM hook that runs before creating a new loan
func (l *Loan) BeforeCreate(tx *gorm.DB) error {

	// Validate Borrower ID Number
	if strings.TrimSpace(l.BorrowerIDNumber) == "" {
		return fmt.Errorf("BorrowerIDNumber cannot be empty")
	}

	// Validate Principal Amount
	if l.PrincipalAmount <= 0 {
		return fmt.Errorf("PrincipalAmount must be greater than zero")
	}

	// Validate Rate
	if l.Rate <= 0 {
		return fmt.Errorf("Rate cannot be negative")
	}

	return nil
}
