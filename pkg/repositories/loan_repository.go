package repositories

import (
	"github.com/amartha/LoanService/pkg/models"
	"gorm.io/gorm"
)

// LoanRepository handles database operations for loans
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new instance of LoanRepository
func NewLoanRepository(db *gorm.DB) *loanRepository {
	return &loanRepository{
		db: db,
	}
}

type LoanRepository interface {
	Create(db *gorm.DB, loan *models.Loan) error
}

// Create inserts a new loan into the database
func (r *loanRepository) Create(db *gorm.DB, loan *models.Loan) error {
	// Set state to proposed
	loan.SetStateToProposed()

	// Set remaining investment amount to principal amount
	loan.SetDefaultRemainingInvestmentAmount()

	// Create a new loan using model struct
	if err := db.Create(&loan).Error; err != nil {
		return err
	}

	return nil
}
