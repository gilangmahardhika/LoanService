# LoanService

[![Tests](https://github.com/gilangmahardhika/LoanService/actions/workflows/go-tests.yml/badge.svg)](https://github.com/gilangmahardhika/LoanService/actions/workflows/go-tests.yml)
[![codecov](https://codecov.io/gh/gilangmahardhika/LoanService/graph/badge.svg?token=MRVH02QJU3)](https://codecov.io/gh/gilangmahardhika/LoanService)

A Go Fiber-based microservice for loan and investment operations. This service provides RESTful APIs for managing loans, investments, and related financial transactions.

## Features

- Investment Management
  - Create and manage investment records
  - Handle investment agreements
  - Process investment transactions
- Loan Processing
  - Loan application handling
  - Loan status management
- Email Notifications
  - Agreement notifications
    
## Prerequisites

- Go 1.21+
- Go Fiber v2
- PostgreSQL 12+
- Docker (optional, for containerization)

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to download dependencies
3. Run `go run main.go` to start the server

## Project Structure

```
pkg/
├── models/          # Data models for investments and loans
├── repositories/    # Database operations
├── handlers/        # HTTP request handlers
├── mailers/         # Email notification services
└── tests/          # Test suites
```

## Available Endpoints

- `/`: Root endpoint with basic service information
- `/health`: Health check endpoint
- `/api/v1/loans`: Loan management endpoints

## Environment Variables

The application uses `godotenv` to manage environment variables. Create a `.env` file in the project root with the following variables:

- `APP_NAME`: Name of the application
- `APP_ENV`: Environment (development, staging, production)
- `SERVER_PORT`: Port to run the server (default: 8080)
- `DATABASE_URL`: Database connection string
- `LOG_LEVEL`: Logging level
- `TEST_DATABASE_URL`: Database connection string for testing (default: same as `DATABASE_URL`)

Example `.env` file:
```
APP_NAME=LoanService
APP_ENV=development
SERVER_PORT=8080
DATABASE_URL=postgresql://localhost:5432/loanservice
LOG_LEVEL=debug
TEST_DATABASE_URL=postgresql://localhost:5432/loanservice_test
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

Database migrations are managed using GORM's auto-migration feature. The application will automatically create and update tables based on the model definitions.

To manually run migrations:
```bash
go run cmd/migrate/main.go
```

## Development

### Code Quality

- Use `go fmt ./...` to format code
- Use `go vet ./...` to check for potential issues
- Run tests with `go test ./...`
- Generate test coverage with `go test ./... -coverprofile=coverage.out`

### API Documentation

API documentation is available at `/docs` when running in development mode.

## Configuration

The application configuration is managed through:

1. Environment variables (via `.env` file)
2. Command line flags
3. Configuration files in `configs/`

Modify `main.go` to adjust server configurations and add more routes as needed.

## Testing

Run the test suite:
```bash
go test ./... -v
```

Generate and view test coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Deployment

### Docker

Build the Docker image:
```bash
docker build -t loanservice .
```

Run the container:
```bash
docker run -p 8080:8080 --env-file .env loanservice
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
