# Frontend Architecture (React)

## Overview
TeraRush frontend is a modern React application built with TypeScript, Vite for fast development, and Tailwind CSS v4 for styling. The architecture emphasizes component reusability, separation of concerns, and maintainability.

## Technology Stack
- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite (lightning-fast development)
- **Styling**: Tailwind CSS v4 with CSS variables
- **UI Components**: Shadcn/UI (headless components)
- **Icons**: Lucide React
- **Markdown Editor**: @uiw/react-md-editor
- **Markdown Renderer**: custom MarkdownRenderer (using marked)
- **HTTP Client**: Axios
- **Form Management**: React Hook Form
- **Animations**: Framer Motion, GSAP

## Directory Structure

```
web/
├── index.html              # HTML entry point
├── tsconfig.json           # TypeScript configuration
├── vite.config.ts          # Vite configuration
├── tailwind.config.ts      # Tailwind CSS configuration
├── eslint.config.js        # ESLint rules
├── package.json            # Dependencies
│
└── src/
    ├── main.tsx           # React app entry
    ├── App.tsx            # Root component
    ├── index.css          # Global styles (minimal)
    │
    ├── assets/
    │   ├── globals.css    # Tailwind directives, CSS variables
    │   ├── fonts/         # Custom fonts
    │   └── images/        # Static images
    │
    ├── components/
    │   ├── ui/            # Atomic Shadcn/UI primitives
    │   │   ├── button.tsx
    │   │   ├── input.tsx
    │   │   ├── dialog.tsx
    │   │   ├── dropdown-menu.tsx
    │   │   ├── card.tsx
    │   │   └── ...
    │   │
    │   └── fragments/     # Composite components
    │       ├── navbar.tsx
    │       ├── sidebar.tsx
    │       ├── blog-card.tsx
    │       ├── blog-form.tsx
    │       ├── user-profile.tsx
    │       ├── markdown-editor.tsx
    │       └── markdown-renderer.tsx
    │
    ├── pages/             # Route-level components
    │   ├── home.tsx
    │   ├── blog-list.tsx
    │   ├── blog-detail.tsx
    │   ├── blog-create.tsx
    │   ├── blog-edit.tsx
    │   ├── login.tsx
    │   ├── register.tsx
    │   ├── profile.tsx
    │   └── not-found.tsx
    │
    ├── lib/
    │   ├── api/          # API client and services
    │   │   ├── client.ts    # Axios instance configuration
    │   │   ├── auth-service.ts
    │   │   ├── blog-service.ts
    │   │   └── user-service.ts
    │   │
    │   └── utils/        # Utility functions
    │       ├── cn.ts        # Classname combiner
    │       ├── date.ts      # Date formatting
    │       ├── validation.ts
    │       └── constants.ts
    │
    ├── hooks/            # Custom React hooks
    │   ├── useFetchBlogs.ts
    │   ├── useFetchBlogDetail.ts
    │   ├── useCreateBlog.ts
    │   ├── useUpdateBlog.ts
    │   ├── useFetchUser.ts
    │   ├── useAuth.ts
    │   ├── useTheme.ts
    │   └── useLocalStorage.ts
    │
    ├── contexts/         # Global state & providers
    │   ├── AuthContext.tsx      # Authentication state
    │   ├── ThemeContext.tsx      # Dark/light theme
    │   ├── NotificationContext.tsx # Toast notifications
    │   └── providers.tsx         # Combined providers
    │
    └── types/
        ├── auth.ts
        ├── blog.ts
        ├── user.ts
        └── api.ts
```

## Component Patterns

### 1. UI Components (Atomic)
Located in `components/ui/`, these are reusable, framework-agnostic primitives from Shadcn.

```typescript
// components/ui/button.tsx
import React from "react"
import { cn } from "@/lib/utils/cn"

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "outline" | "ghost"
  size?: "sm" | "md" | "lg"
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "md", ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          "px-4 py-2 rounded-lg font-medium transition",
          variant === "default" && "bg-primary text-white hover:bg-primary/90",
          variant === "outline" && "border border-input hover:bg-muted",
          variant === "ghost" && "hover:bg-muted",
          className
        )}
        {...props}
      />
    )
  }
)

export { Button }
```

### 2. Fragment Components (Composite)
Composed of multiple UI components, specific to app domain.

```typescript
// components/fragments/blog-card.tsx
import React from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Blog } from "@/types/blog"

interface BlogCardProps {
  blog: Blog
  onRead?: (id: string) => void
  onEdit?: (id: string) => void
  onDelete?: (id: string) => void
}

export const BlogCard: React.FC<BlogCardProps> = ({ blog, onRead, onEdit, onDelete }) => {
  return (
    <Card className="p-6 hover:shadow-lg transition">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-xl font-semibold text-foreground">{blog.title}</h3>
        <span className="text-xs bg-muted px-2 py-1 rounded">{blog.category}</span>
      </div>
      <p className="text-sm text-muted-foreground mb-4">{blog.excerpt}</p>
      <div className="flex items-center justify-between">
        <span className="text-xs text-muted-foreground">By {blog.author.name}</span>
        <Button 
          size="sm" 
          onClick={() => onRead?.(blog.id)}
        >
          Read More
        </Button>
      </div>
    </Card>
  )
}
```

### 3. Page Components (Route-level)
Full-page components that combine multiple fragments.

```typescript
// pages/blog-list.tsx
import React, { useState, useEffect } from "react"
import { useFetchBlogs } from "@/hooks/useFetchBlogs"
import { BlogCard } from "@/components/fragments/blog-card"
import { Button } from "@/components/ui/button"
import { useNavigate } from "react-router-dom"

const BlogListPage: React.FC = () => {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const { data: blogs, loading, error } = useFetchBlogs({ page })

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>

  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Blogs</h1>
        <Button onClick={() => navigate("/blogs/create")}>New Blog</Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {blogs.map(blog => (
          <BlogCard
            key={blog.id}
            blog={blog}
            onRead={(id) => navigate(`/blogs/${id}`)}
            onEdit={(id) => navigate(`/blogs/${id}/edit`)}
          />
        ))}
      </div>
    </div>
  )
}

export default BlogListPage
```

## API Integration Pattern

### 1. Typed API Client
```typescript
// lib/api/client.ts
import axios from "axios"

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    "Content-Type": "application/json",
  },
})

// Add token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("auth_token")
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle errors globally
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("auth_token")
      window.location.href = "/login"
    }
    return Promise.reject(error)
  }
)

export default apiClient
```

### 2. Service Layer
```typescript
// lib/api/blog-service.ts
import apiClient from "./client"
import { Blog, CreateBlogRequest } from "@/types/blog"

export const blogService = {
  async getBlogs(page: number = 1, limit: number = 10) {
    const response = await apiClient.get<{ data: Blog[] }>("/blogs", {
      params: { page, limit },
    })
    return response.data.data
  },

  async getBlogBySlug(slug: string) {
    const response = await apiClient.get<Blog>(`/blogs/${slug}`)
    return response.data
  },

  async createBlog(data: CreateBlogRequest) {
    const response = await apiClient.post<Blog>("/blogs", data)
    return response.data
  },

  async updateBlog(id: string, data: Partial<CreateBlogRequest>) {
    const response = await apiClient.put<Blog>(`/blogs/${id}`, data)
    return response.data
  },

  async deleteBlog(id: string) {
    await apiClient.delete(`/blogs/${id}`)
  },
}
```

### 3. Custom Hook Pattern
```typescript
// hooks/useFetchBlogs.ts
import { useState, useEffect } from "react"
import { blogService } from "@/lib/api/blog-service"
import { Blog } from "@/types/blog"

interface UseFetchBlogsOptions {
  page?: number
  limit?: number
}

export const useFetchBlogs = (options: UseFetchBlogsOptions = {}) => {
  const { page = 1, limit = 10 } = options
  const [data, setData] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchBlogs = async () => {
      try {
        setLoading(true)
        const blogs = await blogService.getBlogs(page, limit)
        setData(blogs)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch blogs")
      } finally {
        setLoading(false)
      }
    }

    fetchBlogs()
  }, [page, limit])

  return { data, loading, error }
}
```

## State Management

### Authentication Context
```typescript
// contexts/AuthContext.tsx
import React, { createContext, useContext, useState } from "react"

interface User {
  id: string
  email: string
  name: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(localStorage.getItem("auth_token"))

  const login = async (email: string, password: string) => {
    // Authentication logic
  }

  const logout = () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem("auth_token")
  }

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used within AuthProvider")
  return context
}
```

## Styling Strategy

### Tailwind CSS v4 with CSS Variables
```css
/* assets/globals.css */
@import "tailwindcss";

@theme {
  --color-primary: oklch(48.6% 0.253 265.75);
  --color-primary-foreground: oklch(98.2% 0.025 270);
  --color-secondary: oklch(78.5% 0.155 180);
  --color-muted: oklch(92% 0.04 270);
  --color-foreground: oklch(20% 0.1 265);
  --color-background: oklch(100% 0 0);
  --color-border: oklch(92% 0.04 270);
}
```

### Component Styling
```typescript
// Always use Tailwind classes, never hardcode colors
className="bg-primary text-primary-foreground border border-border hover:bg-primary/90"

// Use CSS variables for custom values
className="text-[hsl(var(--foreground))]"
```

## Theme Management

```typescript
// Ensure dark mode class syncs with theme
<div className={isDarkMode ? "dark" : ""} data-color-mode={isDarkMode ? "dark" : "light"}>
  <MDEditor />
</div>
```

## Type Safety

### Global Types
```typescript
// types/api.ts
export interface ApiResponse<T> {
  data: T
  status: number
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total_items: number
    total_pages: number
  }
}

// types/blog.ts
export interface Blog {
  id: string
  title: string
  slug: string
  excerpt: string
  content: string
  author: User
  tags: string[]
  views_count: number
  created_at: string
  updated_at: string
}

export interface CreateBlogRequest {
  title: string
  excerpt: string
  content: string
  tags?: string[]
}
```

## Performance Best Practices

### Code Splitting
```typescript
// routes.tsx
const BlogCreatePage = lazy(() => import("@/pages/blog-create"))
const BlogEditPage = lazy(() => import("@/pages/blog-edit"))

export const routes = [
  { path: "/blogs/create", element: <Suspense fallback={<Loading />}><BlogCreatePage /></Suspense> },
]
```

### Memoization
```typescript
// Prevent unnecessary re-renders
const BlogCard = React.memo(({ blog, onRead }: BlogCardProps) => {
  return <Card>{/* ... */}</Card>
})
```

### Image Optimization
```typescript
// Use Next.js-like image optimization or native lazy loading
<img src={imageUrl} alt={title} loading="lazy" />
```

## Environment Variables
```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=TeraRush
```

## Development Workflow

1. **Component Development**: Start with UI components in Storybook (optional)
2. **Fragment Assembly**: Combine UI components into features
3. **Page Integration**: Build page components using fragments
4. **API Integration**: Connect to backend via services
5. **State Management**: Add context when needed
6. **Styling**: Use Tailwind, follow design system
7. **Testing**: Write component and integration tests

## Build & Deployment

```bash
# Development
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Linting
npm run lint
```
