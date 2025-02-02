# LoanService

[![codecov](https://codecov.io/gh/gilangmahardhika/LoanService/graph/badge.svg?token=MRVH02QJU3)](https://codecov.io/gh/gilangmahardhika/LoanService)

A Go Fiber-based microservice for loan-related operations.

## Prerequisites

- Go 1.21+
- Go Fiber v2

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to download dependencies
3. Run `go run main.go` to start the server

## Available Endpoints

- `/`: Root endpoint with basic service information
- `/health`: Health check endpoint

## Environment Variables

The application uses `godotenv` to manage environment variables. Create a `.env` file in the project root with the following variables:

- `APP_NAME`: Name of the application
- `APP_ENV`: Environment (development, staging, production)
- `SERVER_PORT`: Port to run the server (default: 8080)
- `DATABASE_URL`: Database connection string
- `LOG_LEVEL`: Logging level

Example `.env` file:
```
APP_NAME=LoanService
APP_ENV=development
SERVER_PORT=8080
DATABASE_URL=postgresql://localhost:5432/loanservice
LOG_LEVEL=debug
```

## Database Setup

The application uses GORM with PostgreSQL. Ensure you have PostgreSQL installed and create a database for the application.

### Database Connection

Connection details are managed through the `DATABASE_URL` environment variable in the `.env` file. 

Example PostgreSQL connection string:
```
DATABASE_URL=postgresql://username:password@localhost:5432/loanservice
```

### Migrations

Migrations will be added in future updates to manage database schema.

## Development

- Use `go fmt ./...` to format code
- Use `go vet ./...` to check for potential issues

## Configuration

Modify `main.go` to adjust server configurations and add more routes as needed.
