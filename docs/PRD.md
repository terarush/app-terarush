# TeraRush - Product Requirements Document (PRD)

**Status**: Active Development  
**Version**: 2.0  
**Last Updated**: July 15, 2026

---

## 1. Executive Summary

**TeraRush** is a full-stack blogging platform that empowers users to create, publish, and discover quality content. The platform features a modern, responsive interface for content creators and readers, with a robust backend supporting scalability and performance.

### Key Goals
- Provide intuitive content creation experience with markdown support
- Build community through organized, discoverable content
- Establish scalable, maintainable architecture for future expansion
- Ensure excellent user experience across all devices

---

## 2. Project Overview

### 2.1 Product Vision
A modern, open-source blogging platform that combines the simplicity of traditional blogs with the power of markdown and contemporary web technologies.

### 2.2 Problem Statement
Existing blogging platforms are often heavyweight, hard to customize, or lack flexibility. TeraRush provides a lightweight, developer-friendly alternative with clean code and extensible architecture.

### 2.3 Target Users
- **Content Creators**: Writers, developers, technical bloggers
- **Readers**: Users discovering and reading quality content
- **Developers**: Those interested in full-stack development examples

### 2.4 Success Criteria (MVP)
- ✅ Users can register and authenticate
- ✅ Users can create, read, update, and delete blog posts
- ✅ Markdown support with live preview
- ✅ Responsive frontend on desktop and mobile
- ✅ Fast API responses (<200ms p95)
- ✅ Modular, maintainable codebase

---

## 3. Features & Functional Requirements

### 3.1 Core Features (Phase 1 - MVP)

#### 3.1.1 User Authentication
- **Register**: New users can create accounts with email
- **Login**: Existing users authenticate with email/password
- **JWT Tokens**: Secure, stateless authentication
- **Session**: Tokens expire after 24 hours
- **Password Security**: Bcrypt hashing, minimum 8 characters

**Requirements**:
- Email validation (RFC 5322 compliant)
- Strong password enforcement
- Error handling for duplicate accounts
- Secure token refresh mechanism

#### 3.1.2 Blog Management (CRUD)
- **Create**: Authenticated users create new blog posts
- **Read**: All users view published blogs (paginated)
- **Update**: Authors edit their own posts
- **Delete**: Authors delete their own posts
- **Slug Generation**: Auto-generate URL-friendly slugs from titles

**Fields**:
- Title (3-200 characters)
- Excerpt (10-500 characters)
- Content (Markdown format, 50+ characters)
- Tags (max 5, 2-30 chars each)
- Featured Image (optional)
- Author attribution
- Timestamps (created, updated)
- View counter

**Requirements**:
- Pagination: 10 items per page, customizable
- Search: Full-text search across titles and content
- Draft support: (Future phase)
- Revision history: (Future phase)

#### 3.1.3 Content Creation Editor
- **Markdown Editor**: @uiw/react-md-editor integration
- **Live Preview**: Real-time markdown rendering
- **Auto-save**: (Future: Local storage drafts)
- **Image Upload**: Upload and embed images
- **Formatting Tools**: Buttons for markdown syntax

**Supported Markdown**:
- Headings (H1-H6)
- Text formatting (bold, italic, strikethrough)
- Lists (ordered, unordered, nested)
- Code blocks with syntax highlighting
- Links and images
- Tables
- Blockquotes
- Horizontal rules

#### 3.1.4 Blog Discovery
- **Homepage**: Featured/recent blogs
- **Blog Listing**: Paginated view of all blogs
- **Search**: Filter by title, content, tags
- **Blog Detail**: Full blog view with author info
- **Related Posts**: (Future: Based on tags)

#### 3.1.5 User Profile
- **View Profile**: User information and published blogs
- **Profile Info**: Name, email, bio (future)
- **My Blogs**: Dashboard of user's posts
- **Statistics**: Post count, view metrics (future)

### 3.2 Non-Functional Requirements

#### 3.2.1 Performance
- API response time: <200ms (p95)
- Database query time: <100ms (p95)
- Frontend initial load: <3s (3G network)
- Time to Interactive: <5s
- Page load Lighthouse score: 90+

#### 3.2.2 Scalability
- Support 1000+ concurrent users initially
- Stateless backend for horizontal scaling
- Database connection pooling
- Future: CDN for static assets
- Future: Redis caching layer

#### 3.2.3 Security
- HTTPS in production
- CORS: Restricted to known origins
- XSS prevention: React escaping + CSP (future)
- CSRF: JWT-based (stateless)
- Rate limiting: (Recommended, future)
- Input validation: All user inputs validated
- Password hashing: Bcrypt with salt
- No sensitive data in logs

#### 3.2.4 Reliability
- API uptime: 99.5%
- Graceful error handling
- Database backups: (Future)
- Health check endpoint: `/health`
- Error tracking: (Future via Sentry)

#### 3.2.5 Usability
- Mobile-responsive design
- Accessible components (WCAG 2.1 AA target)
- Keyboard navigation support
- Clear error messages
- Loading states and feedback
- Dark mode support (future)

---

## 4. Technical Architecture

### 4.1 Technology Stack

**Backend**
- Runtime: Go 1.24.0+
- Framework: Echo v4 (REST API)
- Database: PostgreSQL 15+
- ORM: GORM 1.25+
- Authentication: JWT
- Logging: Uber Zap

**Frontend**
- Framework: React 18+ with TypeScript
- Build: Vite 5+
- Styling: Tailwind CSS v4
- Components: Shadcn/UI
- HTTP: Axios
- Editor: @uiw/react-md-editor

**Infrastructure**
- Containerization: Docker
- Orchestration: Docker Compose
- Version Control: Git
- Package Management: go mod, npm

### 4.2 Architecture Pattern
- **Backend**: Modular Monolith with DDD
- **Frontend**: Component-based with React Context
- **Communication**: RESTful JSON API
- **Authentication**: JWT Bearer tokens

### 4.3 API Design
- Base URL: `/api/v1`
- Method: REST (GET, POST, PUT, DELETE)
- Format: JSON
- Versioning: URL path-based
- Error format: Consistent error responses

---

## 5. User Stories & Acceptance Criteria

### 5.1 User Registration
```
As a new user,
I want to create an account with my email and password,
So that I can start publishing blogs.

Acceptance Criteria:
- Form validates email format
- Password minimum 8 characters required
- Duplicate emails rejected with clear message
- Success redirects to login
- Error messages are helpful
```

### 5.2 Login
```
As a registered user,
I want to log in with my credentials,
So that I can access my account and publish.

Acceptance Criteria:
- Form accepts email and password
- Invalid credentials show error
- Successful login stores JWT token
- Redirects to dashboard
- Token persists across page refreshes
```

### 5.3 Create Blog Post
```
As an authenticated user,
I want to create a new blog post,
So that I can share my ideas with others.

Acceptance Criteria:
- Title field is required (3-200 chars)
- Markdown editor supports all markdown syntax
- Live preview updates as I type
- Can upload and embed images
- Can add tags (max 5)
- Can set featured image
- Publish button creates post
- Slug auto-generates from title
- Success shows confirmation
```

### 5.4 View Blog Post
```
As a reader,
I want to view a published blog post,
So that I can read interesting content.

Acceptance Criteria:
- Markdown renders correctly
- Author information displayed
- Publication date shown
- Post is responsive on all devices
- Images load efficiently
- Related posts suggested (future)
```

### 5.5 Update Blog Post
```
As a blog author,
I want to edit my published posts,
So that I can fix errors or update information.

Acceptance Criteria:
- Only author can edit their post
- All fields are editable
- Markdown editor has full functionality
- Changes save without losing current content
- Updated timestamp reflects change
- Slug can be updated manually (future)
```

### 5.6 Delete Blog Post
```
As a blog author,
I want to delete my posts,
So that I can remove unwanted content.

Acceptance Criteria:
- Only author can delete
- Confirmation dialog prevents accidents
- Post removed from public listing
- Author's blog count updates
- Deletion logs for audit trail
```

### 5.7 Search Blogs
```
As a reader,
I want to search blogs by keywords,
So that I can find relevant content.

Acceptance Criteria:
- Search works on titles and content
- Results update as I type
- No results shows helpful message
- Results are paginated
- Search terms highlighted in results (future)
```

### 5.8 View User Profile
```
As a reader,
I want to see a user's profile,
So that I can discover their content.

Acceptance Criteria:
- Shows user's published blogs
- Displays post count and statistics
- Lists recent blogs first
- Author bio displayed (future)
- Can follow author (future)
```

---

## 6. Data Model

### 6.1 User Entity
```
User
├── id (UUID, primary key)
├── email (string, unique, indexed)
├── password (string, hashed)
├── name (string)
├── created_at (timestamp)
├── updated_at (timestamp)
└── deleted_at (timestamp, soft delete - future)
```

### 6.2 Blog Entity
```
Blog
├── id (UUID, primary key)
├── title (string, indexed)
├── slug (string, unique, indexed)
├── excerpt (string)
├── content (text, markdown)
├── author_id (FK -> User)
├── image_url (string, optional)
├── tags (array/string)
├── views_count (integer, default 0)
├── published_at (timestamp, nullable)
├── created_at (timestamp, indexed)
├── updated_at (timestamp)
└── deleted_at (timestamp, soft delete - future)
```

### 6.3 Database Indexes
```sql
-- Performance optimization
CREATE INDEX idx_blogs_author_id ON blogs(author_id);
CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_blogs_created_at ON blogs(created_at DESC);
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

---

## 7. API Endpoints (MVP)

### Authentication
- `POST /api/v1/auth/register` - Create account
- `POST /api/v1/auth/login` - Authenticate
- `GET /api/v1/auth/me` - Current user info

### Blogs
- `GET /api/v1/blogs` - List blogs (paginated)
- `GET /api/v1/blogs/:slug` - Get single blog
- `POST /api/v1/blogs` - Create blog (protected)
- `PUT /api/v1/blogs/:id` - Update blog (protected)
- `DELETE /api/v1/blogs/:id` - Delete blog (protected)
- `POST /api/v1/blogs/upload` - Upload image (protected)

### Users
- `GET /api/v1/users/:id` - Get user profile
- `GET /api/v1/users/:id/blogs` - List user's blogs

---

## 8. Development Roadmap

### Phase 1: MVP (Weeks 1-2)
✅ User authentication (register, login)
✅ Blog CRUD operations
✅ Markdown editor with live preview
✅ Responsive frontend
✅ Basic styling with Tailwind
✅ Backend API implementation
✅ Database schema and migrations

### Phase 2: Enhancement (Weeks 3-4)
✅ Draft blog support
✅ Revision history
✅ Comments system
✅ Pagination refinements
✅ Search optimization
✅ Performance optimization
✅ Dark mode support

### Phase 3: Advanced (Weeks 5-6)
✅ Follow system (activity module)
✅ Social sharing (activity module)
✅ Email notifications (notifications module)
✅ Analytics dashboard (analytics module)
✅ Admin panel
✅ Agent system — task orchestration, scheduling, automation

### Phase 4: Production (Week 7+)
- [ ] Error tracking (Sentry)
- [ ] Monitoring (Prometheus, Grafana)
- [ ] CDN integration
- [ ] Cache layer (Redis)
- [ ] Rate limiting
- [ ] API documentation (Swagger)
- [ ] Load testing

---

## 9. Testing Strategy

### 9.1 Backend Testing
- **Unit Tests**: Services, repositories (target >80% coverage)
- **Integration Tests**: Database, API endpoints
- **End-to-End**: Full user workflows
- **Performance Tests**: Load testing, benchmarks

### 9.2 Frontend Testing
- **Unit Tests**: Utility functions, custom hooks
- **Component Tests**: UI components with React Testing Library
- **Integration Tests**: Page components with API mocks
- **E2E Tests**: Critical user flows

### 9.3 Quality Assurance
- Code review: All PRs require review
- Linting: ESLint (frontend), Go fmt (backend)
- Type checking: TypeScript, Go static analysis
- Security: OWASP top 10 compliance

---

## 10. Deployment & DevOps

### 10.1 Development Environment
- Local: Go backend + React frontend + PostgreSQL
- Docker Compose: Containerized setup
- Environment variables: .env files

### 10.2 Production Environment
- Docker containers: API + Database
- PostgreSQL database: Cloud-hosted (AWS RDS)
- Frontend: CDN (CloudFront/CloudFlare)
- Monitoring: Uptime checks, error tracking

### 10.3 CI/CD Pipeline
- GitHub Actions: Automated testing
- Linting: Pre-commit hooks
- Build: Docker image building
- Deployment: Blue-green deployment

---

## 11. Success Metrics

### 11.1 Technical Metrics
- API response time: <200ms (p95)
- Database queries: <100ms (p95)
- Error rate: <0.1%
- Uptime: 99.5%
- Test coverage: >80%

### 11.2 User Metrics
- User registrations: (Target: 100 by month 1)
- Blog posts created: (Target: 500 by month 1)
- Daily active users: (Target: 50 by month 1)
- Average session duration: >5 minutes
- Return user rate: >40%

### 11.3 Code Quality
- Lighthouse score: 90+
- Code review completion: <24 hours
- Bug regression: 0
- Technical debt: Tracked, <10% of backlog

---

## 12. Constraints & Assumptions

### Constraints
- Single backend server initially (horizontal scaling later)
- No real-time features in MVP (websockets future)
- Email verification: Optional (future)
- File storage: Local/cloud storage (TBD)

### Assumptions
- Users have basic markdown knowledge
- UTF-8 text encoding throughout
- PostgreSQL as primary database
- OAuth integration not required for MVP
- Single language (English) for MVP

---

## 13. Glossary

| Term | Definition |
|------|-----------|
| **JWT** | JSON Web Token - Stateless authentication |
| **Markdown** | Lightweight markup language for formatting |
| **DDD** | Domain-Driven Design - Architecture pattern |
| **ORM** | Object-Relational Mapping (GORM) |
| **Slug** | URL-friendly identifier (e.g., "my-blog-post") |
| **DTO** | Data Transfer Object - API data format |
| **CORS** | Cross-Origin Resource Sharing |
| **XSS** | Cross-Site Scripting - Security vulnerability |

---

## 14. Backend Module Catalog

TeraRush uses a modular monolith architecture. Each module is self-contained with domain entities, repositories, services, DTOs, and HTTP handlers.

| Module | Directory | Entities | Purpose |
|--------|-----------|----------|---------|
| **auth** | `modules/auth/` | — | JWT auth, login, register, token refresh, OAuth |
| **users** | `modules/users/` | User | User CRUD, profile, avatar/banner upload, roles |
| **blogs** | `modules/blogs/` | Blog | Blog CRUD, publish/draft, slug, tags, upload |
| **comments** | `modules/comments/` | Comment | Blog comments, nested replies, moderation flags |
| **favorites** | `modules/favorites/` | Favorite | Bookmark/like blogs |
| **assets** | `modules/assets/` | Asset | File upload management, media library |
| **agent** | `modules/agent/` | Agent, AgentTask, AgentSession, AgentCapability, AgentLog, AgentTemplate, AgentSchedule | AI agent orchestration, task lifecycle, scheduling, capability registry, sessions, templates |
| **notifications** | `modules/notifications/` | Notification, NotificationTemplate, NotificationPreference | In-app + email/push notifications, templates, quiet hours, digest |
| **analytics** | `modules/analytics/` | AnalyticsEvent, AnalyticsAggregate, BlogStats, UserAnalytics | Page view & event tracking, blog stats, user analytics, dashboard aggregation |
| **moderation** | `modules/moderation/` | Report, ModerationAction, BannedUser, ContentFilter | Content reporting, user bans, keyword auto-filters, moderation queue |
| **subscriptions** | `modules/subscriptions/` | Subscription, Newsletter, NewsletterSubscriber | Follow authors/blogs, newsletter management, subscriber lists |
| **bookmarks** | `modules/bookmarks/` | Bookmark, Collection, ReadingGoal | Reading list, collections/folders, reading progress, yearly goals |
| **activity** | `modules/activity/` | Activity, ActivityFeed, Follow, EngagementMetric | Activity feed with follower fan-out, follow/unfollow, XP/level, leaderboard |

### Architecture Pattern

```
modules/<name>/
├── module.go                          # Module interface (Name, Initialize, RegisterRoutes, Migrations, Logger)
├── domain/
│   ├── entity/<name>.go               # GORM model with TableName()
│   ├── repository/
│   │   ├── <name>_repository.go       # Interface
│   │   └── <name>_repository_impl.go  # GORM implementation
│   └── service/<name>_service.go      # Business logic
├── handler/<name>_handler.go          # HTTP handlers + route registration
└── dto/
    ├── request/<name>_request.go      # Request DTOs with validate tags
    └── response/<name>_response.go    # Response DTOs + FromEntity converters
```

### Route Prefixes

- **`/api/v1/<resource>`** — Public or user routes (auth-protected)
- **`/api/v1/admin/<resource>`** — Admin-only routes (auth + admin middleware)
- All admin routes use `middleware.Auth` + `middleware.AdminOnly`

### Key Design Decisions

- **Repository**: Uses global `database.DB` (`*gorm.DB`) directly. Empty struct impl.
- **Service**: Business logic only, maps repo errors to domain errors.
- **Handler**: Binds/validates requests, calls service, maps errors to HTTP responses.
- **EventBus**: In-memory pub/sub for cross-module communication.
- **Validation**: `go-playground/validator` via struct tags on request DTOs.

---

## 15. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | May 25, 2025 | TeraRush Team | Initial PRD document |
| 2.0 | July 15, 2026 | TeraRush Team | Added 7 backend modules: agent, notifications, analytics, moderation, subscriptions, bookmarks, activity |

---

## 16. Approval & Sign-Off

- [ ] Project Manager: ___________________
- [ ] Tech Lead (Backend): ___________________
- [ ] Tech Lead (Frontend): ___________________
- [ ] Product Owner: ___________________

---

## Document Version Control

**File**: `docs/PRD.md`  
**Last Updated**: May 25, 2025  
**Status**: Active  
**Next Review**: June 8, 2025