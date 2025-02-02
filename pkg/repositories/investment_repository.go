package repositories

import (
	"fmt"
	"log"

	"github.com/amartha/LoanService/pkg/models"
	"gorm.io/gorm"
)

type investmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) *investmentRepository {
	return &investmentRepository{db}
}

type InvestmentRepository interface {
	Create(db *gorm.DB, investment *models.Investment) error
	createInvestmentWithTransaction(db *gorm.DB, loan *models.Loan, investment *models.Investment) error
	checkInvestorAlreadyInvested(db *gorm.DB, loanID, investorID uint) error
}

func (r *investmentRepository) Create(db *gorm.DB, investment *models.Investment) error {
	// find Loan by id
	loan := &models.Loan{}
	if err := db.First(&loan, investment.LoanID).Error; err != nil {
		return err
	}

	// Check if investor already invested on this loan
	if err := r.checkInvestorAlreadyInvested(db, loan.ID, investment.InvestorID); err != nil {
		return err
	}

	// Can't invest on non approved loan
	if loan.State != models.LoanStatusApproved {
		return fmt.Errorf("loan with id %d is not in approved state", loan.ID)
	}

	// Validate the investment before creating
	if err := investment.ValidateInvestedAmount(loan); err != nil {
		return err
	}

	return r.createInvestmentWithTransaction(db, loan, investment)
}

// checkInvestorAlreadyInvested checks if the investor has already invested in the given loan
func (r *investmentRepository) checkInvestorAlreadyInvested(db *gorm.DB, loanID, investorID uint) error {
	// Check if investor already invested on this loan
	if err := db.Where("loan_id = ? AND investor_id = ?", loanID, investorID).First(&models.Investment{}).Error; err == nil {
		return fmt.Errorf("investor with id %d has already invested on loan with id %d", investorID, loanID)
	}
	return nil
}

// createInvestmentWithTransaction handles the transactional logic for creating an investment
func (r *investmentRepository) createInvestmentWithTransaction(db *gorm.DB, loan *models.Loan, investment *models.Investment) error {
	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&investment).Error; err != nil {
		tx.Rollback()
		return err
	}

	log.Printf("Remaining investment amount of Loan after investment: %v", loan.RemainingInvestmentAmount)
	log.Printf("Investment amount: %v", investment.InvestedAmount)
	log.Printf("Is loan fully invested: %v", isLoanFullyInvested(loan, investment))

	if err := updateLoan(tx, loan, investment); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// isLoanFullyInvested checks if the loan's remaining investment amount
// matches the current investment amount
func isLoanFullyInvested(loan *models.Loan, investment *models.Investment) bool {
	return loan.RemainingInvestmentAmount == investment.InvestedAmount
}

// updateLoanToInvested sets the loan status to invested and saves it to the database
func updateLoan(tx *gorm.DB, loan *models.Loan, investment *models.Investment) error {
	if isLoanFullyInvested(loan, investment) {
		loan.SetStatusToInvested()
	}
	loan.ReduceRemainingInvestmentAmount(investment.InvestedAmount)
	return tx.Save(loan).Error
}
