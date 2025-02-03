package handlers

import (
	"log"

	"gorm.io/gorm"

	"github.com/amartha/LoanService/pkg/models"
	"github.com/amartha/LoanService/pkg/repositories"
	"github.com/gofiber/fiber/v2"
)

type InvestmentHandler struct {
	db                   *gorm.DB
	investmentRepository repositories.InvestmentRepository
}

func NewInvestmentHandler(db *gorm.DB, investmentRepository repositories.InvestmentRepository) *InvestmentHandler {
	return &InvestmentHandler{
		db:                   db,
		investmentRepository: investmentRepository,
	}
}

type InvestRequest struct {
	LoanID         uint    `json:"loan_id"`
	InvestorID     uint    `json:"investor_id"`
	InvestedAmount float64 `json:"invested_amount"`
}

func (h *InvestmentHandler) Invest(c *fiber.Ctx) error {
	var req InvestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate request
	if req.LoanID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Loan ID is required",
		})
	}
	if req.InvestorID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investor ID is required",
		})
	}
	if req.InvestedAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invested amount must be greater than 0",
		})
	}

	// Create investment
	investment := &models.Investment{
		LoanID:         req.LoanID,
		InvestorID:     req.InvestorID,
		InvestedAmount: req.InvestedAmount,
	}

	// Attempt to create investment
	if err := h.investmentRepository.Create(h.db, investment); err != nil {
		log.Printf("Failed to create investment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return investment details
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Investment created successfully",
		"investment_id":  investment.ID,
		"loan_id":        investment.LoanID,
		"investor_id":    investment.InvestorID,
		"invested_amount": investment.InvestedAmount,
	})
}
