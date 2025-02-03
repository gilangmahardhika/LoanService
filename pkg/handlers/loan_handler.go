package handlers

import (
	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LoanHandler struct {
	db             *gorm.DB
	loanRepository repositories.LoanRepository
}

func NewLoanHandler(db *gorm.DB, loanRepository repositories.LoanRepository) *LoanHandler {
	return &LoanHandler{
		db:             db,
		loanRepository: loanRepository,
	}
}

type CreateLoanRequest struct {
	BorrowerIDNumber string  `json:"borrower_id_number"`
	PrincipalAmount  float64 `json:"principal_amount"`
	Rate             float64 `json:"rate"`
}

func (h *LoanHandler) CreateLoan(c *fiber.Ctx) error {
	var req CreateLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate request
	if req.BorrowerIDNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Borrower ID number is required",
		})
	}
	if req.PrincipalAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Principal amount must be greater than 0",
		})
	}
	if req.Rate <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rate must be greater than 0",
		})
	}

	// Create loan model
	loan := &models.Loan{
		BorrowerIDNumber: req.BorrowerIDNumber,
		PrincipalAmount:  req.PrincipalAmount,
		Rate:             req.Rate,
	}

	// Create loan in database
	if err := h.loanRepository.Create(h.db, loan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create loan: " + err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(loan)
}
