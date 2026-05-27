# AI Agent Guidelines & Project Docs

This file provides a concise reference for AI agents working on the TeraRush repository. The canonical and detailed documentation lives under the `docs/` directory — prefer that as the source of truth and keep this file synced with `docs/`.

Short summary (canonical docs are in `docs/`):
- Project: TeraRush — a modular full-stack blogging platform.
- Backend: Go (Echo) with GORM and PostgreSQL.
- Frontend: React + TypeScript, Vite, Tailwind CSS v4, Shadcn UI primitives.

Key conventions (high level):
- Styling: use Tailwind v4 with CSS variables defined in `web/src/assets/globals.css`; avoid hardcoded hex colors — prefer theme variables (e.g. `hsl(var(--primary))`).
- UI primitives: check `web/src/components/ui` before creating new primitives; use `lucide-react` for icons.
- Markdown: use `@uiw/react-md-editor` and the repository's `MarkdownRenderer.tsx` (uses `marked`). Ensure editor is wrapped with `data-color-mode` to follow theme.
- Backend: follow the modular-monolith layout under `/modules`; add new features as new module folders. Use DTOs for request/response and `go-playground/validator` struct tags for validation.

Where to look in `docs/`:
- `docs/PRD.md` — product requirements and roadmap (MVP, features, data models).
- `docs/structure/` — architecture and module layout.
- `docs/api/` — API endpoint definitions and contracts.
- `docs/conventions/` — coding conventions and style guides.

If you make changes to conventions or architecture, update `docs/` first and then sync this file.

Last sync: see `docs/PRD.md` and `docs/conventions/` for the authoritative content.
