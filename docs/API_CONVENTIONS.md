# TeraRush API & Conventions

## 🛣 API Routes (Backend Modules)

Each module follows a consistent pattern for its endpoints.

### 📝 Blogs Module
- `GET /api/v1/blogs` - List all blogs (supports pagination/filter).
- `GET /api/v1/blogs/:slug` - Get blog by slug.
- `POST /api/v1/blogs` - Create blog (Admin only).
- `PUT /api/v1/blogs/:id` - Update blog (Admin only).
- `DELETE /api/v1/blogs/:id` - Delete blog (Admin only).
- `POST /api/v1/blogs/upload` - Upload blog images.

### 🔐 Auth Module
- `POST /api/v1/auth/login` - User login.
- `POST /api/v1/auth/register` - User registration.
- `GET /api/v1/auth/me` - Get current user profile.

---

## 🏗 Coding Conventions

### Backend (Go)
- **Framework**: `Echo` for routing and handlers.
- **Validation**: Uses `go-playground/validator`.
- **Config**: Managed via `spf13/viper` (reads from `.env` or YAML).
- **Error Handling**: Custom error response structure returned via `handler`.

### Frontend (React)
- **Styling**: Tailwind CSS v4 using modern `@theme` syntax in `globals.css`.
- **State Management**: `React Hook Form` for forms, `Context API` for global state.
- **Components**: Primitive components from Shadcn/UI (in `web/src/components/ui`).
- **Icons**: Always use `lucide-react`.
- **Markdown**: 
  - Editor: `@uiw/react-md-editor`.
  - Renderer: `MarkdownRenderer.tsx` using `marked`.

---

## 🚀 Environment Variables
- `PORT`: Server port (Backend).
- `DATABASE_URL`: PostgreSQL connection string.
- `VITE_API_URL`: Backend API URL for frontend.
- `JWT_SECRET`: Secret key for token signing.
