# Technology Stack Overview

## Architecture
**Modular Full-Stack Application** - Frontend and backend operate as separate services communicating via REST API.

---

## Backend

### Runtime & Framework
| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go (Golang) | 1.24.0+ |
| **Web Framework** | Echo | v4 |
| **Architecture** | Modular Monolith | DDD |

### Database & ORM
| Component | Technology | Version |
|-----------|-----------|---------|
| **Primary DB** | PostgreSQL | 15+ |
| **ORM** | GORM | 1.25+ |
| **Driver** | gorm.io/driver/postgres | 1.5+ |
| **Migration** | GORM Auto-Migration | Built-in |

### Authentication & Security
| Component | Technology | Version |
|-----------|-----------|---------|
| **Authentication** | JWT | golang-jwt/jwt v3.2.2 |
| **Password Hashing** | bcrypt | go-playground/validator v9.31.0 |
| **Input Validation** | go-playground/validator | v9.31.0 |

### Logging & Monitoring
| Component | Technology | Version |
|-----------|-----------|---------|
| **Structured Logging** | Uber Zap | 1.27.0 |
| **Log Rotation** | lumberjack | 2.0.0 |
| **Configuration** | Viper | 1.20.0 |
| **In-Memory Cache** | go-cache | 2.1.0 |

### DevOps & Deployment
| Component | Technology |
|-----------|-----------|
| **Containerization** | Docker |
| **Container Orchestration** | Docker Compose |
| **Development Server** | Air (hot reload) |
| **Build Tool** | Make |

### Key Dependencies
```
github.com/labstack/echo v3.3.10 - Web framework
github.com/go-playground/validator v9.31.0 - Validation
github.com/golang-jwt/jwt v3.2.2 - JWT tokens
go.uber.org/zap v1.27.0 - Logging
gorm.io/gorm v1.25.12 - ORM
gorm.io/driver/postgres v1.5.11 - PostgreSQL driver
```

---

## Frontend

### Runtime & Build
| Component | Technology | Version |
|-----------|-----------|---------|
| **Framework** | React | 18+ |
| **Language** | TypeScript | Latest |
| **Build Tool** | Vite | 5+ |
| **Node** | Node.js | 18+ |
| **Package Manager** | npm | 9+ |

### UI & Styling
| Component | Technology | Version |
|-----------|-----------|---------|
| **CSS Framework** | Tailwind CSS | v4 |
| **Component Library** | Shadcn/UI | Latest |
| **Icons** | Lucide React | 0.564+ |
| **CSS-in-JS** | CSS Variables + Tailwind | Built-in |

### Editor & Content
| Component | Technology | Version |
|-----------|-----------|---------|
| **Markdown Editor** | @uiw/react-md-editor | 4.1.0+ |
| **Markdown Renderer** | marked.js | Latest (custom wrapper) |
| **Markdown Preview** | @uiw/react-markdown-preview | 5.2.0+ |

### HTTP & State Management
| Component | Technology | Version |
|-----------|-----------|---------|
| **HTTP Client** | Axios | 1.13.5+ |
| **Form Management** | React Hook Form | 5.2.2+ |
| **Validation** | @hookform/resolvers | 5.2.2+ |
| **State Management** | React Context | Built-in |
| **Storage** | js-cookie | 3.0.5+ |

### Animations & Effects
| Component | Technology | Version |
|-----------|-----------|---------|
| **Animation Library** | Framer Motion | 12.34.0+ |
| **Animation Engine** | GSAP | 3.14.2+ |
| **Scroll Library** | Lenis | 1.3.23+ |

### Accessibility & UI
| Component | Technology | Version |
|-----------|-----------|---------|
| **UI Primitives** | Radix UI | 1.1+ |
| **Utilities** | clsx | 2.1.1+ |
| **CVA** | class-variance-authority | 0.7.1+ |
| **Vite Plugin** | @tailwindcss/vite | 4.1.18+ |

### Development Tools
| Component | Technology |
|-----------|-----------|
| **Linting** | ESLint |
| **Code Formatting** | Prettier (via ESLint) |
| **Type Checking** | TypeScript |
| **Dev Server** | Vite HMR |

### Key Dependencies
```
react 18+ - UI library
typescript - Type safety
vite - Build tool
tailwindcss v4 - Styling
@shadcn/ui - Components
axios - HTTP client
react-hook-form - Form management
@uiw/react-md-editor - Markdown editor
framer-motion - Animations
lucide-react - Icons
```

---

## Communication & APIs

### API Protocol
- **Standard**: REST
- **Format**: JSON
- **Version**: /api/v1
- **Authentication**: JWT Bearer Token

### API Base URL
- **Development**: http://localhost:8080
- **Production**: (Configured via environment variables)

### Response Format
```json
{
  "data": {},
  "status": 200,
  "message": "Success"
}
```

---

## Infrastructure

### Local Development
```
Frontend: Vite Dev Server (http://localhost:5173)
Backend: Go Application (http://localhost:8080)
Database: PostgreSQL (localhost:5432)
```

### Docker Compose Services
- `api`: Go backend application
- `postgres`: PostgreSQL database
- `frontend`: React development server (optional)

### Environment Configuration
```
.env.example - Template for environment variables
.env - Local configuration (git-ignored)
```

---

## Performance Targets

### Backend
- Response time: <200ms (p95)
- Database queries: <100ms (p95)
- Concurrent connections: 1000+
- Memory: <500MB

### Frontend
- Initial load: <3s (3G network)
- Time to Interactive: <5s
- Lighthouse Score: 90+
- Bundle Size: <500KB (gzipped)

---

## Security & Compliance

### Backend
- HTTPS enforced in production
- JWT token expiration: 24 hours
- Password hashing: bcrypt
- Input validation: go-playground/validator
- CORS: Configurable by origin
- Rate limiting: Recommended (to implement)

### Frontend
- XSS prevention: React escaping
- CSRF protection: JWT-based (backend)
- Secure cookies: HttpOnly flags
- Content Security Policy: Recommended (to implement)

---

## Supported Browsers

| Browser | Version |
|---------|---------|
| Chrome | Latest 2 versions |
| Firefox | Latest 2 versions |
| Safari | Latest 2 versions |
| Edge | Latest 2 versions |

---

## Development & Deployment

### Development Commands
```bash
# Backend
go run main.go

# Frontend
npm run dev

# Docker Compose
docker-compose up

# Hot reload
air (backend)
npm run dev (frontend)
```

### Production Build
```bash
# Frontend
npm run build

# Backend
go build -o app

# Docker
docker build -t terrarust:latest .
```

---

## Version Management

- **Go Modules**: go.mod (backend)
- **npm Packages**: package.json & package-lock.json (frontend)
- **Docker**: Dockerfile & docker-compose.yml

---

## Monitoring & Logging

### Backend Logging
- **Level**: DEBUG, INFO, WARN, ERROR
- **Format**: JSON (structured)
- **Output**: stdout + file (logs/)

### Frontend Logging
- **Console logs**: Development only
- **Error tracking**: Recommended (to implement)

---

## Database Schema

### Tables
- `users` - User accounts
- `blogs` - Blog posts
- `(expandable)` - Additional feature tables

### Features
- Timestamps: created_at, updated_at
- Soft deletes: (optional)
- Indexes: On frequently queried columns
- Relationships: Foreign keys with ON CASCADE

---

## Future Considerations

### Planned Upgrades
- [ ] WebSocket support for real-time features
- [ ] GraphQL alternative API
- [ ] Redis caching layer
- [ ] Message queue (Bull/RabbitMQ)
- [ ] Full-text search (PostgreSQL or Elasticsearch)
- [ ] Image optimization (CDN)
- [ ] Monitoring (Prometheus, Grafana)
- [ ] Error tracking (Sentry)

### Scalability
- Horizontal scaling: Stateless Go backend
- Database scaling: PostgreSQL replication
- Frontend: CDN distribution
- Caching: Redis layer