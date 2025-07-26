# ⚡ Volt - API Client for Teams

Collaborative API client with real-time features, auto-import from your codebase, and team synchronization. Built for modern development teams.

## 🚀 Features

- **Real-time Collaboration**: Work together with your team in real-time with live cursors and shared executions
- **Auto-Import APIs**: Automatically discover and sync API endpoints from NestJS, Express, and other frameworks
- **Electric UI**: Beautiful electric blue interface designed for developer productivity
- **Team Workspaces**: Organize your APIs in shared workspaces with proper access control
- **Request Builder**: Intuitive interface for building and testing API requests
- **Environment Management**: Multiple environments with variable synchronization

## 🏗️ Tech Stack

- **Frontend**: SvelteKit + TypeScript
- **Backend**: Go + Gin Framework
- **Database**: MongoDB
- **Real-time**: WebSockets
- **Deployment**: Docker + Docker Compose

## 🚦 Quick Start

### Prerequisites

- Node.js 18-24 (for frontend) - **Node 23.7.0 supported with engine bypass**
- pnpm 9+ (package manager)
- Go 1.21+ (for backend)
- MongoDB 6.0+ (or use Docker)
- Docker & Docker Compose (recommended)

### Option 1: Docker Development (Recommended)

```bash
# Clone the repository
git clone <repo-url>
cd volt

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

Services will be available at:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- MongoDB: localhost:27017

### Option 2: Local Development

**Backend Setup:**
```bash
cd backend
cp .env.example .env
# Edit .env with your MongoDB connection
go mod tidy
go run cmd/server/main.go
```

**Frontend Setup:**
```bash
cd frontend
npm install
npm run dev
```

## 🔧 Development Commands

### Backend
```bash
cd backend
go mod tidy                    # Install dependencies
go run cmd/server/main.go      # Start server
go test ./...                  # Run tests
go test -race ./...            # Run tests with race detection
```

### Frontend
```bash
cd frontend
pnpm install                   # Install dependencies
pnpm run dev                   # Start development server
pnpm run build                 # Build for production
pnpm run check                 # Type checking
```

### Docker
```bash
docker-compose up -d           # Start all services
docker-compose down            # Stop all services
docker-compose logs -f backend # View backend logs
docker-compose up --build      # Rebuild and start
```

## 📡 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `DELETE /api/auth/logout` - User logout

### Workspaces
- `GET /api/workspaces` - List workspaces
- `POST /api/workspaces` - Create workspace
- `GET /api/workspaces/:id` - Get workspace
- `PUT /api/workspaces/:id` - Update workspace
- `DELETE /api/workspaces/:id` - Delete workspace

### Collections & Requests
- `GET /api/workspaces/:id/collections` - List collections
- `POST /api/workspaces/:id/collections` - Create collection
- `GET /api/collections/:id/requests` - List requests
- `POST /api/collections/:id/requests` - Create request
- `POST /api/requests/:id/execute` - Execute request

## 🎨 Brand Colors

- **Primary**: `#2563EB` (Electric Blue)
- **Secondary**: `#1E40AF` (Deep Blue)
- **Accent**: `#FBBF24` (Lightning Yellow)
- **Background**: `#F8FAFC` (Light Gray)
- **Success**: `#10B981` (Emerald)
- **Error**: `#EF4444` (Red)

## 📊 Current Status

**Phase 1 MVP Implementation - IN PROGRESS**

✅ **Completed:**
- Project structure setup with pnpm + Makefile
- Backend API with authentication & all handlers
- Frontend with electric blue theme
- Authentication pages (login/register)
- Dashboard with workspace management
- Full request builder interface
- Response viewer with formatting
- Request execution functionality
- Collection management UI
- Docker development environment
- MongoDB integration with proper indexing

🚧 **In Progress:**
- Environment variables management
- WebSocket foundation for real-time features

📋 **Next Steps:**
- Environment variables UI
- Basic WebSocket for collaboration
- Request history tracking
- Auto-import framework packages

## 📁 Project Structure

```
volt/
├── backend/                 # Go API server
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── auth/            # JWT authentication
│   │   ├── db/              # Database connection
│   │   └── models/          # Data models
│   └── Dockerfile
├── frontend/                # SvelteKit application
│   ├── src/
│   │   ├── routes/          # Pages and API routes
│   │   └── lib/             # Components and utilities
│   └── Dockerfile
├── docs/                    # Documentation
├── scripts/                 # Database initialization
└── docker-compose.yml       # Development environment
```

## 🔐 Environment Variables

**Backend (.env):**
```
MONGO_URI=mongodb://62.171.174.152:27018/volt?username=98xXlLdyb66ekZwLUi4wgs80Gm1B6Bln&password=Ztkth0pbhC7yroHGYC5IRZVzXCsZTC6Z
REDIS_URI=redis://62.171.174.152:6387
REDIS_PASSWORD=Tv6gBhP8fXq7WJsNrD9y4YkZmLc3VuKx
JWT_SECRET=volt-super-secret-jwt-key-2024-electric-blue-api
PORT=8080
CORS_ORIGIN=http://localhost:5173
```

**Frontend (.env):**
```
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=Volt
VITE_BRAND_COLOR=#2563EB
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

## 📋 Roadmap

### Phase 1: Core Features (Current)
- ✅ Authentication & workspace creation
- 🚧 Request builder with electric blue UI
- ⏳ Response viewer and execution
- ⏳ Basic collection management

### Phase 2: Real-time Collaboration
- ⏳ WebSocket implementation
- ⏳ Live cursor tracking
- ⏳ Shared request execution
- ⏳ Team presence indicators

### Phase 3: Auto-Import & Advanced Features
- ⏳ NestJS auto-import package
- ⏳ CLI tools for sync
- ⏳ Multi-framework support
- ⏳ Advanced collaboration features

## 📄 License

This project is licensed under the MIT License.

---

**Built with ⚡ for developers by developers**