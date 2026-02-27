# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a monorepo containing multiple independent projects:

| Project | Description | Tech Stack |
|---------|-------------|------------|
| `hotgo-2.0/` | Full-stack admin system with Go backend and Vue frontend | GoFrame 2.9.4, Go 1.24+, Vue3, NaiveUI, TypeScript |
| `soybean-admin/` | Frontend-only admin template | Vue3, Vite7, TypeScript, Pinia, UnoCSS, NaiveUI |
| `gf2/` | Basic GoFrame 2 template project | GoFrame 2.7.1, Go 1.18+ |

---

## HotGo 2.0 (`hotgo-2.0/`)

### Commands

**Backend (from `hotgo-2.0/server/`):**
```bash
go mod tidy                          # Install dependencies
go run main.go                       # Start all services (HTTP, Queue, Cron)
go run main.go http                  # Start HTTP service only
go run main.go queue                 # Start message queue service only
go run main.go cron                  # Start cron job service only
go run main.go help                  # View CLI help
gf run main.go                       # Hot reload (requires gf CLI)
```

**Frontend (from `hotgo-2.0/web/`):**
```bash
pnpm install                         # Install dependencies
pnpm dev                             # Start dev server (port 8001)
pnpm build                           # Production build
```

### Environment Requirements

- Node.js >= 20.0.0
- Go >= 1.23
- GoFrame >= 2.9.4
- MySQL >= 5.7 or PostgreSQL >= 14

### Project Structure

**Backend (`server/`):**
```
server/
├── addons/              # Plugin modules (each plugin is self-contained)
│   └── hgexample/       # Example plugin with own api/controller/logic/router
├── api/                 # API input/output definitions
│   ├── admin/           # Admin API definitions
│   ├── api/             # Public API definitions
│   ├── home/            # Home page API definitions
│   └── websocket/       # WebSocket API definitions
├── internal/            # Private business logic (Go internal package)
│   ├── cmd/             # CLI commands (http, queue, cron, tools)
│   ├── controller/      # Request handlers/controllers
│   ├── dao/             # Data access objects (auto-generated)
│   ├── logic/           # Business logic implementation
│   ├── model/           # Data models (entity, do, input)
│   ├── router/          # Route registration
│   └── service/         # Service interfaces
├── manifest/config/     # Configuration files
└── utility/             # Utility functions
```

**Frontend (`web/`):**
```
web/
├── src/
│   ├── api/             # API request definitions
│   ├── components/      # Reusable Vue components
│   ├── views/           # Page components
│   ├── router/          # Vue Router configuration
│   ├── store/           # Pinia stores
│   └── locales/         # i18n translations
└── build/               # Build scripts
```

### Key Architecture Patterns

- **Multi-application entry**: Admin (backend), Home (frontend), Api (public API), WebSocket
- **Plugin system**: Each plugin in `addons/` has isolated api/controller/logic/router
- **Layered architecture**: API → Controller → Logic → DAO → Database
- **Service interface pattern**: Interfaces in `service/`, implementations in `logic/`
- **JWT + Casbin**: Authentication via JWT, authorization via Casbin

### Database Setup

1. Import `storage/data/hotgo.sql` into MySQL/PostgreSQL
2. Copy `manifest/config/config.example.yaml` to `manifest/config/config.yaml`
3. Update `database.default.link` with your database credentials
4. Update `hack/config.yaml` for code generation

---

## SoybeanAdmin (`soybean-admin/`)

### Commands

```bash
pnpm install                         # Install dependencies
pnpm dev                             # Start dev server (mock data mode)
pnpm dev:prod                        # Start with production API
pnpm build                           # Production build
pnpm build:test                      # Test build
pnpm lint                            # Run ESLint
pnpm typecheck                       # TypeScript type check
pnpm commit                          # Interactive git commit
```

### Environment Requirements

- Node.js >= 20.19.0
- pnpm >= 10.5.0

### Project Structure

```
soybean-admin/
├── src/
│   ├── layouts/         # Layout components
│   ├── views/           # Page views
│   ├── router/          # File-based routing (Elegant Router)
│   ├── store/           # Pinia stores
│   ├── service/         # API service layer
│   ├── locales/         # i18n translations
│   └── theme/           # Theme settings
└── packages/            # Monorepo packages (axios, hooks, utils, etc.)
```

### Key Architecture Patterns

- **pnpm monorepo**: Shared packages in `packages/` directory
- **Elegant Router**: File-based routing with auto-generated route declarations
- **Theme system**: Built-in UnoCSS with multiple theme presets
- **Permission routing**: Supports both static and dynamic (backend) routes

---

## gf2 (`gf2/`)

### Commands

```bash
go mod tidy                          # Install dependencies
go run main.go                       # Start server
gf run main.go                       # Hot reload (requires gf CLI)
```

### Project Structure

Standard GoFrame single-repo template with:
- `api/` - API definitions
- `internal/cmd/` - CLI commands
- `internal/controller/` - Controllers
- `manifest/config/` - Configuration

---

## Code Generation (HotGo)

HotGo includes a built-in code generator for CRUD operations:

1. Configure database in `hack/config.yaml`
2. Access code generator from admin panel
3. Select table and configure fields
4. Generate frontend views, API, and backend logic

---

## Git Commit Convention

- HotGo and SoybeanAdmin follow Conventional Commits
- Use `pnpm commit` in soybean-admin for interactive commit messages