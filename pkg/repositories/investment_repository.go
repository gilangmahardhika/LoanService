package repositories

import (
	"fmt"

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
}

func (r *investmentRepository) Create(db *gorm.DB, investment *models.Investment) error {
	// find Loan by id
	loan := &models.Loan{}
	if err := db.Where("id = ?", investment.LoanID).First(loan).Error; err != nil {
		return err
	}

	// Can't invest on non approved loan
	if loan.State != "approved" {
		return fmt.Errorf("loan with id %d is not in approved state", loan.ID)
	}

	// Validate the investment before creating
	if err := investment.ValidateInvestedAmount(loan); err != nil {
		return err
	}

	if err := db.Create(&investment).Error; err != nil {
		return err
	}
	return nil
}
