# AI Agent Guidelines & Context

This document provides essential context and instructions for AI agents working on the TeraRush project.

## 🧠 Project Core Identity
TeraRush is a modular full-stack application.
- **Backend**: Go (Golang) with Echo framework.
- **Frontend**: React (TypeScript) with Vite and Tailwind CSS v4.
- **Database**: PostgreSQL with GORM.

## 🛠 Critical Conventions for Agents

### 1. Styling & UI
- **Tailwind v4**: Use modern CSS variables (e.g., `oklch` colors) found in `web/src/assets/globals.css`.
- **Shadcn UI**: Always check `web/src/components/ui` for existing primitives before building new UI.
- **Icons**: Exclusively use `lucide-react`.

### 2. Editor & Content
- **Markdown Editor**: Use `@uiw/react-md-editor`.
- **Markdown Renderer**: Use the custom `MarkdownRenderer.tsx` which uses `marked`.
- **Theme Sync**: Ensure any new UI component supports the `dark` class and syncs with `ThemeProvider`.

### 3. Backend Logic
- **Modular Monolith**: New features should be added as new folders in `/modules`.
- **DTOs**: Always use Data Transfer Objects for API requests/responses.
- **Validation**: Use struct tags with `go-playground/validator`.

## 📂 Documentation Navigation
When exploring the project, refer to these subdirectories:
- `docs/structure/`: Architecture details.
- `docs/api/`: Endpoint definitions.
- `docs/conventions/`: Coding standards.

## ⚠️ Common Pitfalls to Avoid
- **Duplicate Styles**: Do not use hardcoded hex colors; use Shadcn variables (e.g., `hsl(var(--primary))`).
- **Markdown Issues**: Ensure the `@uiw/react-md-editor` is wrapped with `data-color-mode` to match the current theme.
- **Package Management**: Use `npm` for frontend and `go mod` for backend.
