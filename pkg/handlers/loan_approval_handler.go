package handlers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/amartha/LoanService/pkg/helpers"
	"github.com/amartha/LoanService/pkg/repositories"
	"github.com/gofiber/fiber/v2"
)

type LoanApprovalHandler struct {
	db             *gorm.DB
	loanRepository repositories.LoanRepository
}

func NewLoanApprovalHandler(db *gorm.DB, loanRepository repositories.LoanRepository) *LoanApprovalHandler {
	return &LoanApprovalHandler{
		db:             db,
		loanRepository: loanRepository,
	}
}

type ApproveLoanRequest struct {
	LoanID     uint   `json:"loan_id"`
	ApprovedBy uint   `json:"approved_by"`
	VisitProof string `json:"visit_proof"`
}

func (h *LoanApprovalHandler) ApproveLoan(c *fiber.Ctx) error {
	var req ApproveLoanRequest
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
	if req.ApprovedBy == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Approved by is required",
		})
	}
	if req.VisitProof == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visit proof is required",
		})
	}

	// Decode base64 visit proof
	decodedVisitProof, err := base64.StdEncoding.DecodeString(req.VisitProof)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid base64 visit proof",
		})
	}

	// Detect file extension
	fileExt := helpers.DetectFileExtension(decodedVisitProof)

	// Create uploads directory if not exists
	uploadsDir := filepath.Join(".", "uploads")
	log.Printf("Attempting to create uploads directory: %s", uploadsDir)
	
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create uploads directory: %v", err),
		})
	}

	// Create loan-specific directory
	loanDir := filepath.Join(uploadsDir, strconv.FormatUint(uint64(req.LoanID), 10), "visit_proof")
	log.Printf("Attempting to create loan-specific directory: %s", loanDir)
	
	if err := os.MkdirAll(loanDir, 0755); err != nil {
		log.Printf("Failed to create loan-specific directory: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create loan-specific directory: %v", err),
		})
	}

	// Generate unique filename with detected extension
	filename := fmt.Sprintf("visit_proof_%s%s", time.Now().Format("20060102_150405"), fileExt)
	filePath := filepath.Join(loanDir, filename)

	// Save the file
	if err := ioutil.WriteFile(filePath, decodedVisitProof, 0644); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save visit proof",
		})
	}

	// Approve loan in database
	if err := h.loanRepository.SetStateToApproved(h.db, req.LoanID, req.ApprovedBy, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve loan: " + err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Loan approved successfully",
		"visit_proof": filePath,
	})
}
