# TeraRush Project Structure

TeraRush is a full-stack application with a Go backend and a React (Vite) frontend. It follows a modular architecture designed for scalability and maintainability.

## 🏗 Backend Structure (Go)

The backend is located in the root directory and follows standard Go project layouts.

### `/internal`
Core application logic and shared packages.
- **/app**: Entry point for application initialization and module registration.
- **/pkg**: Shared infrastructure components:
  - `config`: Environment and configuration management.
  - `database`: DB connections and migrations (PostgreSQL).
  - `server`: HTTP server setup (using Fiber/Gin).
  - `jwt`, `middleware`, `logger`: Cross-cutting concerns.

### `/modules`
Domain-driven modules. Each module is self-contained.
- **/[module-name]**: (e.g., `blogs`, `auth`, `users`)
  - `domain`: Entities and Repository interfaces.
  - `handler`: HTTP controllers/handlers.
  - `dto`: Request and response data transfer objects.
  - `module.go`: Wire up dependencies for the module.

---

## 💻 Frontend Structure (React + TypeScript)

The frontend is located in the `/web` directory.

### `/web/src`
- **/components**:
  - `ui`: Shadcn/UI primitive components.
  - `fragments`: Complex UI parts (e.g., `MarkdownEditor`, `BlogForm`).
  - `elements`: Smaller reusable UI elements.
- **/pages**: Route-based components (Blog, Dashboard, Admin).
- **/lib/api**: API client and service functions using Axios.
- **/contexts**: React Contexts for global state (Theme, Auth).
- **/assets**: Global CSS (`globals.css`) and static files.

---

## 🔌 API & Integration

- **Client**: Axios instance configured in `web/src/lib/api/client.ts`.
- **Base URL**: Managed via `.env` (usually `VITE_API_URL`).
- **Auth**: JWT-based authentication stored in `localStorage` or Cookies.

---

## 🛠 Tech Stack

- **Backend**: Go, Fiber/Gin, PostgreSQL, GORM.
- **Frontend**: React 18, Vite, TypeScript, Tailwind CSS v4, Shadcn/UI.
- **Markdown**: `@uiw/react-md-editor` for editing, `marked` for rendering.
