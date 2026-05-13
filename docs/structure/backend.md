# Backend Structure (Go)

The backend follows a Modular Monolith approach.

## Core Directories
- **/internal/app**: Application bootstrapping and dependency injection.
- **/internal/pkg**: Shared infrastructure (database, jwt, logger, middleware).
- **/modules**: Domain-driven feature modules.

## Module Pattern
Each module in `/modules/[name]` typically contains:
- `domain/`: Entities and interface definitions.
- `handler/`: HTTP controllers (Echo handlers).
- `dto/`: Request/Response structures.
- `repository/`: Database implementation (GORM).
- `module.go`: Module initialization logic.
