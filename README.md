# Employee Management API

## Features

- **Dual Transport**: gRPC on separate ports
- JWT-based authentication (gRPC metadata)
- Employee CRUD operations with soft delete
- Salary calculations with country-based tax rules
- Salary metrics (min/max/avg by country, avg by job title)
- Clean Architecture with clear layer separation
- Docker and Docker Compose support

## Technology Stack

- **Language**: Go 1.24.3
- **Framework**: Go Kit (service-oriented architecture)
- **Router**: gRPC
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Logging**: go-kit/log
- **Serialization**: Protocol Buffers (gRPC)
- **Containerization**: Docker + Docker Compose

## Project Structure

```
employee-api/
├── cmd/server/            # Application entry point
├── internal/
│   ├── domain/            # Business entities and repository interfaces
│   │   ├── entity/        # User, Employee
│   │   ├── valueobject/   # Country, Salary
│   │   └── repository/    # Repository interfaces
│   ├── usecase/           # Business logic (services)
│   │   ├── auth/          # Authentication service
│   │   ├── employee/      # Employee service
│   │   └── salary/        # Salary calculation service
│   ├── transport/
│   │   └── grpc/          # gRPC server implementations
│   ├── infrastructure/    # External concerns
│   │   ├── persistence/   # Database implementations
│   │   ├── auth/          # JWT manager
│   │   └── config/        # Configuration
│   └── pkg/               # Shared utilities
├── proto/             # Protocol Buffer definitions
├── docs/                  # Documentation
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.24 or later
- PostgreSQL 16+ (or Docker)
- Make (optional)
- (Bloomrpc, for testing gRPC)

### Running Locally

1. Clone the repository and navigate to the project directory

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Start PostgreSQL (if not using Docker):
   ```bash
   # Ensure PostgreSQL is running with:
   # - Database: employee_db
   # - User: postgres
   # - Password: postgres
   ```

4. Run the application:
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

   This starts both servers:
   - gRPC server on `localhost:50051`

### Running with Docker

```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## gRPC API

The gRPC server runs on port `50051` by default and provides the same functionality as the REST API.

### Services

| Service | Methods |
|---------|---------|
| `auth.v1.AuthService` | `Register`, `Login` |
| `employee.v1.EmployeeService` | `CreateEmployee`, `GetEmployee`, `ListEmployees`, `UpdateEmployee`, `DeleteEmployee` |
| `salary.v1.SalaryService` | `CalculateNetSalary`, `GetSalaryStatsByCountry`, `GetAvgSalaryByJobTitle` |

### Authentication

For authenticated endpoints, pass the JWT token in gRPC metadata:
```
authorization: Bearer <token>
```

## Tax Rules

| Country | Tax Rate |
|---------|----------|
| India | 10% |
| United States | 12% |
| All Others | 0% |

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| GRPC_PORT | 50051 | gRPC server port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | employee_db | Database name |
| DB_SSLMODE | disable | SSL mode |
| JWT_SECRET | (required) | JWT signing secret |
| JWT_EXPIRATION | 24h | Token expiration |
| JWT_ISSUER | employee-api | Token issuer |

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Development

```bash
# Download dependencies
make deps

# Build binary
make build

# Regenerate proto files (requires protoc, protoc-gen-go, protoc-gen-go-grpc)
make gen-proto
```

## Architecture

This project follows Clean Architecture principles:

1. **Domain Layer** (innermost): Entities, value objects, repository interfaces
2. **Use Case Layer**: Business logic services
3. **Interface Adapters**: Go Kit, DTOs
4. **Frameworks Layer** (outermost): database implementations
