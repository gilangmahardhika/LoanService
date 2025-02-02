package repositories

import (
	"fmt"
	"time"

	"github.com/amartha/LoanService/pkg/models"
	"gorm.io/gorm"
)

// LoanRepository handles database operations for loans
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new instance of LoanRepository
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

type LoanRepository interface {
	Create(db *gorm.DB, loan *models.Loan) error
	SetStateToApproved(db *gorm.DB, id uint, approvedBy uint, visitProof string) error
	SetStateToDisbursed(db *gorm.DB, id uint, disbursedBy uint) error
	GetByID(db *gorm.DB, id uint) (*models.Loan, error)
}

// Create inserts a new loan into the database
func (r *loanRepository) Create(db *gorm.DB, loan *models.Loan) error {
	// Set state to proposed
	loan.SetStateToProposed()

	// Set remaining investment amount to principal amount
	loan.SetDefaultRemainingInvestmentAmount()

	// Create a new loan using model struct
	if err := db.Create(&loan).Error; err != nil {
		return fmt.Errorf("failed to create loan: %w", err)
	}

	return nil
}

func (r *loanRepository) GetByID(db *gorm.DB, id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := db.Where("id = ?", id).First(&loan).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve loan with id %d: %w", id, err)
	}
	return &loan, nil
}

// Function for set the state to approved after the loan is approved
func (r *loanRepository) SetStateToApproved(db *gorm.DB, id uint, approvedBy uint, visitProof string) error {
	// If the loan is already approved, return an error
	// Get the loan by id
	loan := &models.Loan{}
	if err := db.Where("id = ?", id).First(loan).Error; err != nil {
		return fmt.Errorf("failed to find loan with id %d: %w", id, err)
	}

	// Return error if the loan state is not proposed
	if loan.State != models.LoanStatusProposed {
		return fmt.Errorf("loan with id %d is not in proposed state", id)
	}

	// Update the loan state to approved
	loan.State = models.LoanStatusApproved
	loan.ApprovedBy = &approvedBy
	loan.VisitProof = &visitProof
	loan.ApprovedAt = &[]time.Time{time.Now()}[0]

	if err := db.Save(loan).Error; err != nil {
		return fmt.Errorf("failed to update loan state to approved: %w", err)
	}

	return nil
}

// Set state to disbursed if remaining investment amount is 0
func (r *loanRepository) SetStateToDisbursed(db *gorm.DB, id uint, disbursedBy uint) error {
	// Get the loan by id
	loan := &models.Loan{}
	if err := db.Where("id = ?", id).First(loan).Error; err != nil {
		return fmt.Errorf("failed to find loan with id %d: %w", id, err)
	}

	// Return error if the loan state is not approved
	if loan.State != models.LoanStatusInvested {
		return fmt.Errorf("loan with id %d is not in invested state", id)
	}

	// Set the state to disbursed
	loan.State = models.LoanStatusDisbursed
	disbursedByStr := disbursedBy
	loan.DisbursedBy = &disbursedByStr
	now := time.Now()
	loan.DisbursedAt = &now

	if err := db.Save(loan).Error; err != nil {
		return fmt.Errorf("failed to update loan state to disbursed: %w", err)
	}

	return nil
}
