# Coding Conventions & Standards

## General Principles
- Use clear, descriptive naming for variables, functions, and files
- Follow DRY (Don't Repeat Yourself) principle - avoid code duplication
- Keep functions/components small and focused (single responsibility)
- Write self-documenting code with meaningful comments only when logic is non-obvious
- Maintain consistency within the codebase

## Go (Backend)

### File Organization
- One package per directory
- Keep files under 500 lines when possible
- Group related functionality together

### Naming Conventions
```go
// Constants: UPPER_SNAKE_CASE (package-level)
const MaxRetries = 3

// Functions: camelCase, descriptive
func getUserByID(id string) (*User, error)
func validateEmailFormat(email string) bool

// Types: PascalCase
type User struct {}
type BlogService struct {}

// Interface methods: descriptive verbs
type Repository interface {
    GetByID(id string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id string) error
}
```

### Code Style
- Use `gofmt` for automatic formatting
- Keep line length under 100 characters
- Error handling: check errors immediately
```go
if err != nil {
    return fmt.Errorf("failed to fetch user: %w", err)
}
```

### Echo Handlers
- Always use `echo.Context` as first parameter
- Return errors using echo error handler
```go
func (h *Handler) GetUser(c echo.Context) error {
    userID := c.Param("id")
    user, err := h.service.GetUser(userID)
    if err != nil {
        return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, user)
}
```

### DTOs (Data Transfer Objects)
- Define request/response structs with validation tags
```go
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2"`
}
```

### Database (GORM)
- Use GORM hooks for timestamps (created_at, updated_at)
- Define models with appropriate indexes
- Use transactions for multi-step operations

### Logging
- Use structured logging (zap) at appropriate levels
```go
logger.Info("User created", zap.String("userID", user.ID), zap.String("email", user.Email))
logger.Error("Failed to save blog", zap.Error(err), zap.String("title", blog.Title))
```

---

## React/TypeScript (Frontend)

### File Organization
```
src/
├── components/
│   ├── ui/           # Shadcn/UI primitives
│   └── fragments/    # Composite components
├── pages/           # Route-level components
├── hooks/           # Custom React hooks
├── lib/
│   └── api/         # API clients & services
├── contexts/        # Global state providers
└── assets/          # Styles, icons, images
```

### Naming Conventions
```typescript
// Components: PascalCase
const UserProfile: React.FC = () => {}
const BlogCard: React.FC<BlogCardProps> = () => {}

// Hooks: camelCase, prefix with 'use'
const useFetchBlogs = () => {}
const useThemeMode = () => {}

// Functions: camelCase
const formatDate = (date: Date) => {}
const generateSlug = (title: string) => {}

// Constants: UPPER_SNAKE_CASE or camelCase
const API_BASE_URL = "http://localhost:8080"
const blogsPerPage = 10
```

### Component Structure
```typescript
import React from "react"
import { Button } from "@/components/ui/button"
import { useFetchBlogs } from "@/hooks/useFetchBlogs"

interface BlogCardProps {
  title: string
  excerpt: string
  author: string
}

const BlogCard: React.FC<BlogCardProps> = ({ title, excerpt, author }) => {
  return (
    <div className="p-4 border rounded-lg">
      <h3 className="text-lg font-semibold">{title}</h3>
      <p className="text-sm text-gray-600">{excerpt}</p>
      <span className="text-xs text-gray-500">By {author}</span>
    </div>
  )
}

export default BlogCard
```

### Styling
- Use **Tailwind CSS v4** with CSS variables
- Use `oklch()` color format from globals.css
- Leverage shadcn/UI components for consistency
- No hardcoded hex colors
```typescript
// ✅ Good
className="bg-primary text-primary-foreground"
className="text-[hsl(var(--foreground))]"

// ❌ Avoid
className="bg-#3498db text-white"
```

### TypeScript
- Always define types for props
- Use interfaces for objects
- Avoid `any` - use proper types
```typescript
// ✅ Good
interface User {
  id: string
  email: string
  name: string
}

const User: React.FC<{ user: User }> = ({ user }) => {
  return <div>{user.name}</div>
}

// ❌ Avoid
const User = ({ user: any }) => {}
```

### Hooks & State
- Use functional components with hooks
- Keep hooks custom and reusable
- Use proper dependency arrays in useEffect
```typescript
useEffect(() => {
  fetchBlogs()
}, [search, page]) // Include all dependencies
```

### API Integration
- Use axios client in lib/api
- Handle errors consistently
- Include proper loading/error states
```typescript
const useFetchBlogs = () => {
  const [data, setData] = React.useState<Blog[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    apiClient.get("/blogs")
      .then(res => setData(res.data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  return { data, loading, error }
}
```

### Markdown Editor
- Wrap `@uiw/react-md-editor` with `data-color-mode`
- Always sync with ThemeProvider
```typescript
import MDEditor from "@uiw/react-md-editor"

const BlogEditor = () => {
  return (
    <div data-color-mode="light">
      <MDEditor height={400} />
    </div>
  )
}
```

---

## Git & Version Control

### Commit Messages
Format: `type(scope): subject`

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Examples:**
```
feat(blog): add markdown editor to blog creation
fix(auth): resolve JWT token expiration bug
docs(api): update blog endpoints documentation
refactor(frontend): simplify blog list component
test(backend): add unit tests for user service
```

### Branch Naming
- Feature: `feature/description` (e.g., `feature/blog-comments`)
- Bug fix: `fix/description` (e.g., `fix/auth-jwt-issue`)
- Hotfix: `hotfix/description`
- Documentation: `docs/description`

### Pull Request Guidelines
- Keep PRs focused and reasonably sized (under 400 lines when possible)
- Include clear PR description with context
- Link related issues
- Ensure CI/CD passes before merging
- At least one code review before merge

---

## Testing

### Backend (Go)
- Write tests for services and repositories
- Use table-driven tests for multiple cases
- Aim for >80% code coverage
```go
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        want    *User
        wantErr bool
    }{
        {"valid user", "123", &User{ID: "123"}, false},
        {"invalid user", "invalid", nil, true},
    }
}
```

### Frontend (React)
- Test critical user interactions
- Mock API calls in tests
- Use React Testing Library for component tests

---

## Performance Guidelines

### Backend
- Use database indexes on frequently queried columns
- Implement caching for expensive operations
- Paginate large result sets
- Use connection pooling for database

### Frontend
- Use React.memo for expensive components
- Implement lazy loading for routes
- Optimize images and assets
- Monitor bundle size

---

## Documentation

### Code Comments
- Comment the "why", not the "what"
- Keep comments updated with code changes
```go
// ✅ Good - explains reasoning
// We cache user permissions for 5 minutes to reduce DB queries
// during active sessions
cache.Set(userID, permissions, 5*time.Minute)

// ❌ Avoid - obvious from code
// Set cache
cache.Set(userID, permissions, 5*time.Minute)
```

### Function/Component Documentation
- Include JSDoc/GoDoc comments for public functions
- Document parameters and return values
- Add examples for complex functionality

---

## Security

### Passwords
- Hash passwords using bcrypt
- Never log passwords
- Enforce minimum requirements

### API Tokens
- Use JWT with reasonable expiration (24 hours)
- Implement refresh token mechanism
- Validate tokens on every protected request

### Data Validation
- Validate all user inputs
- Use struct tags with go-playground/validator
- Sanitize data before storing

### CORS & HTTPS
- Enable CORS only for trusted origins
- Use HTTPS in production
- Implement CSRF protection if needed
