package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/amartha/LoanService/pkg/database"
	"github.com/amartha/LoanService/pkg/handlers"
	"github.com/amartha/LoanService/pkg/repositories"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
		// Increase body size limit for base64 image uploads
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// Middleware
	app.Use(cors.New())    // Enable CORS
	app.Use(logger.New())  // Logging middleware
	app.Use(recover.New()) // Panic recovery middleware

	// Initialize repositories
	loanRepo := repositories.NewLoanRepository(db)
	investmentRepo := repositories.NewInvestmentRepository(db)

	// Initialize handlers
	loanHandler := handlers.NewLoanHandler(db, loanRepo)
	loanApprovalHandler := handlers.NewLoanApprovalHandler(db, loanRepo)
	investmentHandler := handlers.NewInvestmentHandler(db, investmentRepo)

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to " + os.Getenv("APP_NAME"),
			"env":     os.Getenv("APP_ENV"),
			"status":  "healthy",
		})
	})

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Loan routes
	loans := v1.Group("/loans")
	loans.Post("/", loanHandler.CreateLoan)
	loans.Post("/approve", loanApprovalHandler.ApproveLoan)
	loans.Post("/invest", investmentHandler.Invest)

	// Start server
	port := ":" + os.Getenv("SERVER_PORT")
	if port == ":" {
		port = ":8080" // default port
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}
