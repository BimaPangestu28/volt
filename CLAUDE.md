# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Volt is an API client project designed as a collaborative alternative to Postman. The project features:

- **Frontend**: Svelte + SvelteKit + TypeScript
- **Backend**: Go (Gin framework) 
- **Database**: MongoDB
- **Real-time Features**: WebSockets for live collaboration
- **Auto-Import System**: Framework-specific packages (@volt/nestjs) for automatic API discovery

## Current Development Status

**Phase 1 MVP Implementation - IN PROGRESS**

✅ **Completed:**
- Project structure (frontend & backend directories)
- Backend Go project with Gin framework setup
- Frontend SvelteKit project with TypeScript
- JWT authentication system
- MongoDB models and database connection
- Basic workspace and collection management API
- Electric blue themed landing page
- Docker development environment

🚧 **In Progress:**
- WebSocket foundation for real-time features
- Authentication pages (login/register)
- Dashboard and request builder UI

📋 **Next Steps:**
- Complete authentication flow
- Build request builder interface
- Implement request execution
- Add basic collaboration features

## Architecture Overview

### Frontend Structure (Planned)
```
frontend/
├── src/
│   ├── routes/          # SvelteKit routes
│   ├── lib/
│   │   ├── components/  # Reusable components
│   │   ├── stores/      # Svelte stores (state)
│   │   ├── services/    # API calls
│   │   └── utils/       # Helper functions
│   └── app.html
├── static/              # Static assets
└── package.json
```

### Backend Structure (Planned)
```
backend/
├── cmd/server/main.go   # Entry point
├── internal/
│   ├── api/             # HTTP handlers
│   ├── auth/            # Authentication
│   ├── db/              # Database connection
│   ├── models/          # Data models
│   └── services/        # Business logic
│       ├── websocket/   # WebSocket handlers
│       └── collaboration/ # Live collaboration logic
├── pkg/                 # Shared packages
├── go.mod
└── go.sum
```

## Key Features to Implement

### Core API Client Features
- Request builder with full HTTP method support
- Response viewer with JSON formatting
- Collection management and organization
- Environment variables and configuration
- Request history and execution tracking

### Real-time Collaboration
- Multi-user editing with live sync via WebSockets
- Presence indicators showing online users
- Live cursor tracking and shared editing
- Real-time comment system on requests
- Conflict resolution for simultaneous edits

### Auto-Import System
- Framework-specific packages starting with NestJS (@volt/nestjs)
- Automatic API endpoint discovery and synchronization
- CLI tools for manual sync operations
- Dashboard integration for sync status and management
- Support planned for Express.js, Fastify, Go Gin, FastAPI, Spring Boot

## Development Commands

### Frontend Development
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Start development server (http://localhost:5173)
npm run build        # Build for production
npm run check        # Type checking
```

### Backend Development
```bash
cd backend
go mod tidy          # Install dependencies
cp .env.example .env # Copy environment variables
go run cmd/server/main.go  # Start development server (http://localhost:8080)
go test ./...        # Run all tests
go test -race ./...  # Run tests with race detection
go test -cover ./... # Run tests with coverage
```

### Docker Development Environment
```bash
# Start all services (MongoDB, Redis, Backend, Frontend)
docker-compose up -d

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop all services
docker-compose down

# Rebuild services
docker-compose up --build
```

### Testing the Implementation
```bash
# Test backend health
curl http://localhost:8080/health

# Test frontend
open http://localhost:5173
```

## Environment Variables

### Backend (.env)
```
MONGO_URI=mongodb://localhost:27017/volt
REDIS_URI=redis://localhost:6379
JWT_SECRET=your-super-secret-key
PORT=8080
CORS_ORIGIN=http://localhost:5173
WS_ORIGIN=ws://localhost:8080
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_APP_NAME=Volt
VITE_BRAND_COLOR=#2563EB
```

## Database Schema

MongoDB collections are planned for:
- Users (authentication and profiles)
- Workspaces (team organization)
- Collections (request groupings)
- Requests (individual API calls)
- Environments (variable sets)
- WebSocket Connections (real-time collaboration)
- Comments (collaborative discussions)
- Execution History (request/response logs)

## API Endpoints (Planned)

### Authentication
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `DELETE /api/auth/logout`

### Core Resources
- Workspaces: `/api/workspaces/*`
- Collections: `/api/collections/*`
- Requests: `/api/requests/*`
- Environments: `/api/environments/*`

### Real-time Features
- WebSocket: `/ws/connect`
- Presence: `/api/presence/*`
- Comments: `/api/comments/*`
- History: `/api/history/*`

## Technology Stack

### Dependencies (Planned)

**Frontend**: Svelte, SvelteKit, TypeScript, Axios, CodeMirror/Monaco Editor, Socket.io-client, Lucide icons, TailwindCSS

**Backend**: Gin (Go), JWT, Gorilla WebSocket, MongoDB Driver, Redis, bcrypt

## Brand Identity

Electric blue theme with the following color palette:
- Primary: `#2563EB` (Electric Blue)
- Secondary: `#1E40AF` (Deep Blue)  
- Accent: `#FBBF24` (Lightning Yellow)
- Background: `#F8FAFC` (Light Gray)
- Dark: `#0F172A` (Dark Slate)

## Development Workflow

The project is planned in phases:
1. **Phase 1 (MVP)**: Core API client features + NestJS auto-import
2. **Phase 2**: Real-time collaboration + auto-sync dashboard
3. **Phase 3**: Advanced features + multi-framework support

## Project Structure (Implemented)

```
volt/
├── CLAUDE.md                 # This file - guidance for AI agents
├── docker-compose.yml        # Development environment setup
├── docs/
│   └── project.md           # Complete project specification
├── scripts/
│   └── init-mongo.js        # MongoDB initialization
├── backend/                 # Go API server
│   ├── cmd/server/main.go   # Application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers and routing
│   │   ├── auth/            # JWT authentication
│   │   ├── db/              # Database connection
│   │   └── models/          # Data models
│   ├── .env.example         # Environment variables template
│   ├── Dockerfile           # Backend container
│   ├── go.mod               # Go dependencies
│   └── go.sum
└── frontend/                # SvelteKit application
    ├── src/
    │   ├── routes/          # SvelteKit routes
    │   ├── lib/
    │   │   ├── stores/      # Svelte stores (auth, etc.)
    │   │   └── api/         # API client
    │   └── app.html         # HTML template with electric blue theme
    ├── .env.example         # Environment variables template
    ├── Dockerfile           # Frontend container
    ├── package.json         # Node.js dependencies
    └── vite.config.ts       # Vite configuration
```

## Implementation Notes

- **Backend**: Fully functional Go/Gin API with JWT auth, MongoDB integration, CORS setup
- **Frontend**: SvelteKit with TypeScript, electric blue theme, responsive design
- **Database**: MongoDB with proper indexing and validation
- **Development**: Docker Compose for easy local development
- **Authentication**: JWT-based with proper middleware and token management
- **API Design**: RESTful endpoints following the planned specification

## Current API Endpoints (Implemented)

- `GET /health` - Health check
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `DELETE /api/auth/logout` - User logout
- `GET /api/workspaces` - List user workspaces
- `POST /api/workspaces` - Create workspace
- `GET/PUT/DELETE /api/workspaces/:id` - Workspace management
- Collection and Request endpoints are scaffolded but need handlers implementation