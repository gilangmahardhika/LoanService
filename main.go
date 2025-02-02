package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	// Middleware
	app.Use(cors.New())     // Enable CORS
	app.Use(logger.New())   // Logging middleware
	app.Use(recover.New())  // Panic recovery middleware

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

	// Start server
	port := ":" + os.Getenv("SERVER_PORT")
	if port == ":" {
		port = ":8080" // default port
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}
