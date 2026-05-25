# Installation & Setup Guide

## Prerequisites

### Backend
- **Go** 1.24.0 or higher
- **PostgreSQL** 15 or higher
- **Git** for version control

### Frontend
- **Node.js** 18+ with npm 9+
- **Git** for version control

### Optional (For Docker)
- **Docker** 20.10+
- **Docker Compose** 1.29+

---

## Local Setup

### 1. Database Setup (PostgreSQL)

#### macOS (Homebrew)
```bash
brew install postgresql@15
brew services start postgresql@15
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

#### Windows
Download and install from [postgresql.org](https://www.postgresql.org/download/windows/)

#### Create Database
```bash
# Login to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE terrarust;

# Create user
CREATE USER terrarust_user WITH PASSWORD 'secure_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE terrarust TO terrarust_user;

# Exit
\q
```

### 2. Backend Setup

```bash
# Clone repository
git clone <repository-url>
cd TeraRush

# Install Go dependencies
go mod download

# Create environment file
cp .env.example .env

# Edit .env with your configuration
nano .env
```

#### .env Configuration
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=terrarust
DB_USER=terrarust_user
DB_PASSWORD=secure_password

# Server
PORT=8080
ENVIRONMENT=development

# JWT
JWT_SECRET=your-super-secret-key-min-32-chars
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=INFO
```

#### Run Backend
```bash
# Run with hot reload (requires 'air')
air

# Or direct run
go run main.go

# Build binary
go build -o app
```

Backend will be available at `http://localhost:8080`

### 3. Frontend Setup

```bash
# Navigate to frontend
cd web

# Install dependencies
npm install

# Create environment file
cp .env.example .env.local

# Edit configuration if needed
# VITE_API_URL=http://localhost:8080
```

#### Run Frontend
```bash
# Development server with HMR
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

Frontend will be available at `http://localhost:5173`

---

## Docker Setup (Recommended for Quick Start)

### Prerequisites
- Docker and Docker Compose installed

### Quick Start
```bash
# From project root
docker-compose up --build

# Check status
docker-compose ps
```

Services will start:
- **API**: http://localhost:8080
- **Database**: localhost:5432
- **Frontend** (if included): http://localhost:5173

### Docker Commands
```bash
# View logs
docker-compose logs -f api
docker-compose logs -f postgres

# Stop services
docker-compose down

# Remove volumes (careful!)
docker-compose down -v

# Rebuild after code changes
docker-compose up --build
```

---

## Verification Checklist

### Backend
- [ ] Go installed: `go version`
- [ ] PostgreSQL running: `psql -U postgres -d terrarust`
- [ ] Dependencies installed: `go mod download`
- [ ] .env configured correctly
- [ ] Server starts: `go run main.go`
- [ ] Health endpoint: `curl http://localhost:8080/health`

### Frontend
- [ ] Node.js installed: `node --version`
- [ ] npm installed: `npm --version`
- [ ] Dependencies installed: `npm install`
- [ ] Dev server runs: `npm run dev`
- [ ] Can access: http://localhost:5173

### Database
- [ ] PostgreSQL running
- [ ] Database created: `terrarust`
- [ ] User created: `terrarust_user`
- [ ] Migrations applied automatically on server start

---

## Troubleshooting

### PostgreSQL Connection Issues
```bash
# Test connection
psql -h localhost -U terrarust_user -d terrarust

# Reset PostgreSQL (if needed)
brew services restart postgresql@15  # macOS
sudo systemctl restart postgresql    # Linux
```

### Port Already in Use
```bash
# Find process using port
lsof -i :8080  # Backend
lsof -i :5432 # Database
lsof -i :5173 # Frontend

# Kill process
kill -9 <PID>
```

### Node Module Issues
```bash
# Clean install
rm -rf node_modules package-lock.json
npm install
```

### Go Module Issues
```bash
# Verify and tidy modules
go mod tidy
go mod verify
```

---

## First Run Walkthrough

1. **Start PostgreSQL**
   ```bash
   brew services start postgresql@15  # or equivalent for your OS
   ```

2. **Clone and setup backend**
   ```bash
   git clone <repo>
   cd TeraRush
   cp .env.example .env
   # Edit .env with DB credentials
   go run main.go
   ```

3. **In new terminal, setup frontend**
   ```bash
   cd TeraRush/web
   npm install
   npm run dev
   ```

4. **Test the application**
   - Open http://localhost:5173
   - Register a new account
   - Create a blog post
   - View blogs

---

## Development Workflow

### Making Changes

#### Backend
1. Make code changes in `/modules` or `/internal`
2. `air` will auto-reload the server
3. Test endpoints with curl or Postman

#### Frontend
1. Make code changes in `/web/src`
2. Vite HMR updates automatically
3. Check console for TypeScript errors

### Testing
```bash
# Backend
go test ./...
go test -v ./modules/...

# Frontend
npm test
npm run lint
```

### Building
```bash
# Backend
go build -o app
./app

# Frontend
npm run build
# Output in web/dist
```

---

## Environment Variables Reference

### Backend (.env)
```env
# Database Configuration
DB_HOST=localhost              # PostgreSQL host
DB_PORT=5432                   # PostgreSQL port
DB_NAME=terrarust              # Database name
DB_USER=terrarust_user         # Database user
DB_PASSWORD=password           # Database password
DB_SSL_MODE=disable           # SSL mode (disable/require)

# Server Configuration
PORT=8080                      # Server port
ENVIRONMENT=development        # development/production
ALLOWED_ORIGINS=*             # CORS origins

# JWT Configuration
JWT_SECRET=your-secret-key    # Min 32 characters
JWT_EXPIRY=24h                # Token expiration

# Logging
LOG_LEVEL=INFO                # DEBUG/INFO/WARN/ERROR
LOG_OUTPUT=stdout             # stdout/file

# Optional Services
REDIS_URL=                    # For caching (optional)
```

### Frontend (.env.local)
```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=TeraRush
```

---

## Next Steps

1. **Read API Documentation**: `docs/api/`
2. **Review Code Conventions**: `docs/conventions/`
3. **Understand Architecture**: `docs/structure/`
4. **Create First Feature**: Follow patterns in existing modules
5. **Set up IDE**: VSCode extensions for Go and React development

---

## Getting Help

- Check logs: `docker-compose logs`
- Review AGENT.md for guidelines
- Check GitHub Issues
- Report bugs with reproduction steps