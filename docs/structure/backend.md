# Backend Architecture (Go)

## Overview
TeraRush backend is a modular monolith built with Go, using the Echo web framework, GORM for database operations, and PostgreSQL as the primary data store. The architecture follows Domain-Driven Design (DDD) principles with clear separation of concerns.

## Technology Stack
- **Framework**: Echo v4 (lightweight and performant)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator
- **Logging**: Uber's Zap logger
- **Configuration**: Viper
- **Caching**: go-cache (in-memory)

## Directory Structure

```
project-root/
├── main.go                 # Application entry point
├── go.mod & go.sum         # Dependency management
├── .env & .env.example     # Environment variables
├── Dockerfile              # Container image definition
├── docker-compose.yml      # Multi-container setup
│
├── internal/
│   ├── app/                # Core application setup
│   │   ├── app.go         # App initialization
│   │   ├── module.go      # Module interface
│   │   └── config.go      # Configuration loading
│   │
│   └── pkg/                # Shared packages
│       ├── database/       # GORM setup, migrations
│       ├── logger/         # Zap logging wrapper
│       ├── jwt/            # JWT token generation/validation
│       ├── middleware/     # Echo middleware (auth, CORS, etc)
│       └── errors/         # Custom error types
│
├── modules/                # Feature modules (DDD)
│   ├── user/
│   │   ├── domain/
│   │   │   ├── user.go    # User entity
│   │   │   └── repository.go  # Repository interface
│   │   ├── dto/
│   │   │   ├── request.go  # Request DTOs
│   │   │   └── response.go # Response DTOs
│   │   ├── handler/
│   │   │   └── user_handler.go  # HTTP handlers
│   │   ├── repository/
│   │   │   └── user_repo.go    # GORM implementation
│   │   ├── service/
│   │   │   └── user_service.go # Business logic
│   │   └── module.go       # Module registration
│   │
│   └── blog/
│       ├── domain/
│       ├── dto/
│       ├── handler/
│       ├── repository/
│       ├── service/
│       └── module.go
│
└── public/                 # Static files, migrations scripts

```

## Module Pattern (DDD Approach)

Each feature module contains these layers:

### 1. Domain Layer (`domain/`)
- **Purpose**: Defines business entities and interfaces
- **Contains**:
  - Entity structs (e.g., `User`, `Blog`)
  - Repository interfaces (contracts)
  - Domain errors

```go
// domain/user.go
type User struct {
    ID        string    `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// domain/repository.go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### 2. DTO Layer (`dto/`)
- **Purpose**: Request/Response transfer objects
- **Contains**: Input/output structures with validation tags

```go
// dto/request.go
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

// dto/response.go
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}
```

### 3. Handler Layer (`handler/`)
- **Purpose**: HTTP request handling (controllers)
- **Contains**: Echo handlers with request/response management

```go
// handler/user_handler.go
type UserHandler struct {
    service domain.UserService
    logger  *zap.Logger
}

func (h *UserHandler) GetUser(c echo.Context) error {
    id := c.Param("id")
    user, err := h.service.GetUser(c.Request().Context(), id)
    if err != nil {
        return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, user)
}
```

### 4. Repository Layer (`repository/`)
- **Purpose**: Data persistence (database operations)
- **Contains**: GORM implementations of repository interfaces

```go
// repository/user_repo.go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 5. Service Layer (`service/`)
- **Purpose**: Business logic implementation
- **Contains**: Service implementations with domain operations

```go
// service/user_service.go
type UserService struct {
    repo   domain.UserRepository
    logger *zap.Logger
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        s.logger.Error("Failed to get user", zap.Error(err), zap.String("userID", id))
        return nil, err
    }
    return user, nil
}
```

### 6. Module Registration (`module.go`)
- **Purpose**: Module initialization and route registration

```go
// module.go
type Module struct {
    handler    *handler.UserHandler
    repository domain.UserRepository
    service    domain.UserService
    logger     *zap.Logger
}

func (m *Module) Name() string {
    return "user"
}

func (m *Module) Initialize(db *gorm.DB, logger *zap.Logger) error {
    m.logger = logger
    m.repository = repository.NewUserRepository(db)
    m.service = service.NewUserService(m.repository, logger)
    m.handler = handler.NewUserHandler(m.service, logger)
    return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
    group := e.Group(basePath)
    
    group.GET("/:id", m.handler.GetUser)
    group.POST("", m.handler.CreateUser, middlewares.AuthMiddleware)
    group.PUT("/:id", m.handler.UpdateUser, middlewares.AuthMiddleware)
    group.DELETE("/:id", m.handler.DeleteUser, middlewares.AuthMiddleware)
}

func (m *Module) Migrations() []interface{} {
    return []interface{}{
        &domain.User{},
    }
}
```

## Core Packages (`internal/pkg/`)

### Database (`pkg/database/`)
- Initialize PostgreSQL connection
- GORM configuration
- Auto-migration system
- Connection pooling

### Logger (`pkg/logger/`)
- Structured logging with Zap
- Module-aware log levels
- Production-ready output

### JWT (`pkg/jwt/`)
- Token generation
- Token validation
- Claims management
- Refresh token support

### Middleware (`pkg/middleware/`)
- Authentication verification
- Error handling
- CORS configuration
- Request logging

### Errors (`pkg/errors/`)
- Custom error types
- Error wrapping and context
- HTTP error mapping

## Request/Response Flow

```
HTTP Request
    ↓
Middleware (Auth, CORS, Logging)
    ↓
Handler Layer (Parse input, validate)
    ↓
DTO Conversion (Request → Domain)
    ↓
Service Layer (Business logic)
    ↓
Repository Layer (Database operations)
    ↓
Service returns result
    ↓
DTO Conversion (Domain → Response)
    ↓
Handler returns JSON Response
    ↓
HTTP Response
```

## Database Schema Management

- Migrations in GORM struct tags
- Auto-migration on startup
- Migration scripts in `public/migrations/`
- Use transactions for data consistency

## Best Practices

### Error Handling
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Never return nil error with nil data
- Log errors with context (user ID, action, etc)

### Validation
- Validate at DTO level using struct tags
- Validate business rules in Service layer
- Return meaningful error messages

### Logging
- Use appropriate log levels (Debug, Info, Warn, Error)
- Include relevant context (user ID, request ID, etc)
- Avoid logging sensitive data (passwords)

### Concurrency
- Use context for request cancellation
- Implement timeouts for external operations
- Use goroutines judiciously

### Testing
- Unit test services and repositories
- Mock external dependencies
- Use table-driven tests for multiple cases
- Aim for >80% coverage

## Adding a New Module

1. Create `modules/[feature-name]/` directory
2. Implement layers: `domain/`, `dto/`, `handler/`, `repository/`, `service/`
3. Create `module.go` implementing the Module interface
4. Register module in `main.go`
5. Run auto-migrations
6. Test all endpoints

## Configuration

Environment variables in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=terrarust
DB_USER=postgres
DB_PASSWORD=password

JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

LOG_LEVEL=INFO
PORT=8080
```

## Deployment

- Docker image: Multi-stage build
- Docker Compose: App + PostgreSQL
- Health check endpoint: `GET /health`
- Graceful shutdown: Signal handling
