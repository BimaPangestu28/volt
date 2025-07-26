# Volt - API Client Project Documentation

## 🚀 Tech Stack

- **Frontend**: Svelte + SvelteKit + TypeScript
- **Backend**: Go (Gin framework)
- **Database**: MongoDB
- **Authentication**: JWT
- **Real-time**: WebSockets (Live Collaboration)
- **Deployment**: Docker + Docker Compose

## 🎨 Brand Identity

### Electric Blue Palette
- **Primary**: `#2563EB` (Electric Blue)
- **Secondary**: `#1E40AF` (Deep Blue)
- **Accent**: `#FBBF24` (Lightning Yellow)
- **Background**: `#F8FAFC` (Light Gray)
- **Dark**: `#0F172A` (Dark Slate)
- **Success**: `#10B981` (Emerald)
- **Warning**: `#F59E0B` (Amber)
- **Error**: `#EF4444` (Red)

## 📦 Volt Auto-Import Package System

### Overview
Volt provides automatic API discovery and import capabilities through framework-specific packages. Developers can install Volt packages in their projects to automatically sync API endpoints to the Volt dashboard in real-time.

### Supported Frameworks (Phase 1)
- **NestJS** - `@volt/nestjs` (Primary focus for MVP)

### Volt NestJS Package

#### Installation
```bash
npm install @volt/nestjs
# or
yarn add @volt/nestjs
```

#### Basic Setup
```typescript
// main.ts
import { VoltModule } from '@volt/nestjs';

@Module({
  imports: [
    VoltModule.forRoot({
      apiKey: process.env.VOLT_API_KEY,
      workspaceId: process.env.VOLT_WORKSPACE_ID,
      projectName: 'My NestJS API',
      autoSync: true, // Auto-sync on server start
      watchMode: process.env.NODE_ENV === 'development'
    }),
    // ... other modules
  ],
})
export class AppModule {}
```

#### Automatic Detection Features
```typescript
// Volt automatically detects and syncs:

@Controller('users')
export class UsersController {
  
  // @volt-description Get all users with pagination
  // @volt-tags users, admin
  @Get()
  @ApiOperation({ summary: 'Get all users' }) // Uses existing Swagger decorators
  async findAll(@Query() query: PaginationDto) {
    return this.usersService.findAll(query);
  }

  // @volt-environment dev:{{DEV_API_URL}}, prod:{{PROD_API_URL}}
  @Post()
  @ApiBody({ type: CreateUserDto })
  async create(@Body() createUserDto: CreateUserDto) {
    return this.usersService.create(createUserDto);
  }

  @Get(':id')
  async findOne(@Param('id') id: string) {
    return this.usersService.findOne(id);
  }
}
```

#### Configuration Options
```typescript
// volt.config.ts
export default {
  apiKey: process.env.VOLT_API_KEY,
  workspaceId: process.env.VOLT_WORKSPACE_ID,
  
  // Auto-discovery settings
  discovery: {
    includeRoutes: ['/api/**'], // Only sync specific routes
    excludeRoutes: ['/health', '/metrics'], // Exclude system routes
    includeControllers: ['UsersController', 'PostsController'],
    autoDetectSwagger: true, // Use existing Swagger metadata
  },
  
  // Sync settings
  sync: {
    onStartup: true,
    watchMode: true, // Watch for file changes in dev
    batchSize: 50, // Batch API calls for performance
    syncInterval: 30000, // 30 seconds for watch mode
  },
  
  // Collection organization
  organization: {
    groupBy: 'controller', // 'controller', 'tag', 'module'
    collectionPrefix: 'Auto-Sync',
    environmentMapping: {
      development: 'dev',
      staging: 'staging',
      production: 'prod'
    }
  },
  
  // Request generation
  requestGeneration: {
    generateSampleData: true, // Auto-generate sample request bodies
    inferAuthMethods: true, // Detect auth requirements
    generateTests: true, // Basic assertion generation
  }
};
```

#### CLI Commands
```bash
# Install Volt CLI globally
npm install -g @volt/cli

# Initialize Volt in existing NestJS project
volt init --framework nestjs

# Manual sync (alternative to auto-sync)
volt sync

# Watch mode for development
volt watch

# Generate Volt collection from code
volt generate --output collections/

# Validate configuration
volt validate

# Show sync status
volt status
```

#### Advanced Features

**Environment Detection:**
```typescript
// Automatically detects environments from:
// - .env files (.env.development, .env.staging, .env.production)
// - process.env.NODE_ENV
// - Custom environment mapping in config

// Auto-generates environment variables in Volt:
// {{API_BASE_URL}}, {{DB_HOST}}, {{JWT_SECRET}}, etc.
```

**Request Body Generation:**
```typescript
// Uses NestJS DTOs to generate sample requests
export class CreateUserDto {
  @IsString()
  @ApiProperty({ example: 'john.doe@example.com' })
  email: string;

  @IsString()
  @MinLength(8)
  @ApiProperty({ example: 'strongPassword123' })
  password: string;

  @IsOptional()
  @ApiProperty({ example: 'John Doe' })
  name?: string;
}

// Volt auto-generates:
// {
//   "email": "john.doe@example.com",
//   "password": "strongPassword123", 
//   "name": "John Doe"
// }
```

**Authentication Auto-Detection:**
```typescript
@UseGuards(JwtAuthGuard)
@Controller('protected')
export class ProtectedController {
  // Volt automatically detects JWT requirement
  // and sets up Bearer token authentication
}

@ApiSecurity('api-key')
@Controller('api-key-protected')
export class ApiKeyController {
  // Volt detects API key requirement
}
```

**Real-time Sync Flow:**
1. Developer starts NestJS server with Volt package
2. Volt scans controllers, routes, and decorators
3. Generates Volt collections automatically
4. Syncs to dashboard via WebSocket
5. Team members see new APIs instantly
6. File changes trigger incremental sync

#### Dashboard Integration

**Auto-Generated Collections:**
- Collection per controller/module
- Organized folder structure
- Environment variables auto-detected
- Sample requests with realistic data
- Basic test assertions

**Sync Status Indicators:**
- Green: Successfully synced
- Yellow: Sync in progress  
- Red: Sync failed
- Blue: Manual sync required

**Conflict Resolution:**
- Manual edits in Volt preserved
- Auto-sync creates new versions
- Merge conflicts highlighted
- Option to accept/reject auto-changes

### Future Framework Support (Roadmap)
- **Express.js** - `@volt/express`
- **Fastify** - `@volt/fastify`  
- **Gin (Go)** - `github.com/volt/gin`
- **FastAPI (Python)** - `volt-fastapi`
- **Spring Boot (Java)** - `volt-spring-boot`

### Live Collaboration
- **Multi-user editing** - Multiple developers can edit collections simultaneously
- **Live cursor tracking** - See where team members are working in real-time
- **Shared request execution** - Watch API test results live across the team
- **Real-time comments** - Discuss API endpoints directly in the interface
- **Presence indicators** - Know who's online and working on what

### Live Sync
- **Auto-sync collections** - Changes propagate instantly across all devices
- **Environment synchronization** - Variable updates in real-time
- **Request history streaming** - New execution history appears live
- **Notification system** - Team activity and updates
- **Conflict resolution** - Smart merge for simultaneous edits

### WebSocket Message Types
```go
type WSMessage struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Room      string      `json:"room"`
    UserID    string      `json:"user_id"`
    Timestamp time.Time   `json:"timestamp"`
}

// Message types:
// - "collection_update" - Collection changes
// - "request_execute" - Live request execution results
// - "user_join/leave" - User presence in workspace
// - "cursor_move" - Live cursor position tracking
// - "environment_sync" - Environment variable updates
// - "comment_add" - New comments on requests
// - "typing_indicator" - Show who's typing
```

### WebSocket Rooms/Namespaces
- `workspace:{id}` - Workspace-level updates and presence
- `collection:{id}` - Collection-specific changes and collaboration
- `request:{id}` - Individual request editing and comments

## 📁 Project Structure

```
volt/
├── frontend/                 # Svelte frontend
│   ├── src/
│   │   ├── routes/          # SvelteKit routes
│   │   ├── lib/
│   │   │   ├── components/  # Reusable components
│   │   │   ├── stores/      # Svelte stores (state)
│   │   │   ├── services/    # API calls
│   │   │   └── utils/       # Helper functions
│   │   └── app.html
│   ├── static/              # Static assets
│   └── package.json
├── backend/                  # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── auth/            # Authentication
│   │   ├── db/              # Database connection
│   │   ├── models/          # Data models
│   │   ├── services/        # Business logic
│   │   │   ├── websocket/    # WebSocket handlers & real-time
│   │   │   ├── collaboration/ # Live collaboration logic
│   ├── pkg/                 # Shared packages
│   ├── go.mod
│   └── go.sum
├── docker-compose.yml
└── README.md
```

## 🗃️ Database Schema (MongoDB)

### Users Collection
```json
{
  "_id": "ObjectId",
  "username": "string",
  "email": "string",
  "password_hash": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Workspaces Collection
```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "owner_id": "ObjectId",
  "members": ["ObjectId"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Auto-Import Collections
```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "workspace_id": "ObjectId",
  "source": {
    "type": "auto-import",
    "framework": "nestjs|express|gin",
    "project_name": "string",
    "sync_config": {},
    "last_sync": "timestamp"
  },
  "auto_generated": "boolean",
  "manual_edits": ["request_id"],
  "sync_status": "synced|pending|failed",
  "created_by": "ObjectId",
  "requests": ["ObjectId"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```
```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "workspace_id": "ObjectId",
  "created_by": "ObjectId",
  "requests": ["ObjectId"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Collections Collection
```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "workspace_id": "ObjectId",
  "created_by": "ObjectId",
  "requests": ["ObjectId"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```
```json
{
  "_id": "ObjectId",
  "name": "string",
  "method": "GET|POST|PUT|DELETE|PATCH",
  "url": "string",
  "headers": {
    "key": "value"
  },
  "body": {
    "type": "json|form|raw|binary",
    "content": "string"
  },
  "auth": {
    "type": "bearer|basic|apikey",
    "credentials": {}
  },
  "tests": ["string"],
  "collection_id": "ObjectId",
  "created_by": "ObjectId",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Environments Collection
```json
{
  "_id": "ObjectId",
  "name": "string",
  "workspace_id": "ObjectId",
  "variables": {
    "key": "value"
  },
  "created_by": "ObjectId",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### WebSocket Connections Collection
```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId",
  "workspace_id": "ObjectId",
  "connection_id": "string",
  "last_seen": "timestamp",
  "cursor_position": {
    "request_id": "ObjectId",
    "field": "string",
    "position": "number"
  },
  "status": "online|away|offline"
}
```

### Comments Collection
```json
{
  "_id": "ObjectId",
  "request_id": "ObjectId",
  "user_id": "ObjectId",
  "content": "string",
  "mentions": ["ObjectId"],
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "resolved": "boolean"
}
```
```json
{
  "_id": "ObjectId",
  "request_id": "ObjectId",
  "user_id": "ObjectId",
  "request_data": {},
  "response_data": {
    "status": "number",
    "headers": {},
    "body": "string",
    "time": "number"
  },
  "executed_at": "timestamp"
}
```

## 🔧 Backend API Endpoints

### Authentication
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh
DELETE /api/auth/logout
```

### Workspaces
```
GET    /api/workspaces
POST   /api/workspaces
GET    /api/workspaces/:id
PUT    /api/workspaces/:id
DELETE /api/workspaces/:id
POST   /api/workspaces/:id/members
DELETE /api/workspaces/:id/members/:userId
```

### Collections
```
GET    /api/workspaces/:workspaceId/collections
POST   /api/workspaces/:workspaceId/collections
GET    /api/collections/:id
PUT    /api/collections/:id
DELETE /api/collections/:id
```

### Requests
```
GET    /api/collections/:collectionId/requests
POST   /api/collections/:collectionId/requests
GET    /api/requests/:id
PUT    /api/requests/:id
DELETE /api/requests/:id
POST   /api/requests/:id/execute
```

### Environments
```
GET    /api/workspaces/:workspaceId/environments
POST   /api/workspaces/:workspaceId/environments
GET    /api/environments/:id
PUT    /api/environments/:id
DELETE /api/environments/:id
```

### WebSocket & Real-time
```
GET    /ws/connect                    # WebSocket connection endpoint
POST   /api/presence/update           # Update user presence
GET    /api/workspaces/:id/online     # Get online users
POST   /api/requests/:id/comments     # Add comment to request
GET    /api/requests/:id/comments     # Get request comments
PUT    /api/comments/:id              # Update comment
DELETE /api/comments/:id              # Delete comment
```
```
GET    /api/history
GET    /api/requests/:id/history
DELETE /api/history/:id
```

### History
```
GET    /api/history
GET    /api/requests/:id/history
DELETE /api/history/:id
```

### Core Components
- `RequestBuilder.svelte` - Main request form with live collaboration
- `ResponseViewer.svelte` - Display API responses with live updates
- `CollectionTree.svelte` - Sidebar collection navigator with presence
- `EnvironmentSelector.svelte` - Environment switcher with live sync
- `HistoryPanel.svelte` - Request history with real-time updates
- `CodeGenerator.svelte` - Export to various languages
- `CollaborationPanel.svelte` - Live collaboration sidebar
- `CommentThread.svelte` - Request comments and discussions
- `AutoSyncPanel.svelte` - Auto-import project management
- `SyncStatus.svelte` - Show sync status and conflicts
- `ProjectSelector.svelte` - Select synced projects

### Layout Components
- `Sidebar.svelte` - Main navigation with user presence
- `Header.svelte` - Top bar with online users and notifications
- `Modal.svelte` - Reusable modal component
- `Tabs.svelte` - Tab interface with live indicators
- `SplitPane.svelte` - Resizable panels
- `PresenceIndicator.svelte` - Show online users
- `LiveCursor.svelte` - Other users' cursor positions

### Utility Components
- `JsonEditor.svelte` - JSON input with syntax highlighting and live editing
- `HeaderEditor.svelte` - Key-value pair editor with collaboration
- `AuthConfig.svelte` - Authentication configuration
- `TestEditor.svelte` - Write API tests with live feedback
- `NotificationToast.svelte` - Real-time notifications
- `UserAvatar.svelte` - User avatars and status
- `ConflictResolver.svelte` - Resolve auto-sync conflicts
- `FrameworkBadge.svelte` - Show source framework (NestJS, etc.)

## 📦 Dependencies

### Frontend (package.json)
```json
{
  "devDependencies": {
    "@sveltejs/adapter-auto": "^2.0.0",
    "@sveltejs/kit": "^1.20.4",
    "svelte": "^4.0.5",
    "typescript": "^5.0.0",
    "vite": "^4.4.2"
  },
  "dependencies": {
    "axios": "^1.5.0",
    "codemirror": "^6.0.1",
    "monaco-editor": "^0.44.0",
    "socket.io-client": "^4.7.2",
    "lucide-svelte": "^0.285.0",
    "@tailwindcss/typography": "^0.5.10",
    "@volt/nestjs": "^1.0.0"
  }
}
```

### Backend (go.mod)
```go
module api-client-backend

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/gorilla/websocket v1.5.0
    go.mongodb.org/mongo-driver v1.12.1
    golang.org/x/crypto v0.13.0
    github.com/go-redis/redis/v8 v8.11.5
)
```

## 🚀 Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- MongoDB 6.0+
- Docker (optional)

### Development Setup

1. **Clone repository**
```bash
git clone <repo-url>
cd api-client
```

2. **Backend setup**
```bash
cd backend
go mod tidy
cp .env.example .env
# Edit .env with your MongoDB connection string
go run cmd/server/main.go
```

3. **Frontend setup**
```bash
cd frontend
npm install
npm run dev
```

4. **Using Docker**
```bash
docker-compose up -d
```

### Environment Variables

**Backend (.env)**
```
MONGO_URI=mongodb://localhost:27017/volt
REDIS_URI=redis://localhost:6379
JWT_SECRET=your-super-secret-key
PORT=8080
CORS_ORIGIN=http://localhost:5173
WS_ORIGIN=ws://localhost:8080
```

**Frontend (.env)**
```
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_APP_NAME=Volt
VITE_BRAND_COLOR=#2563EB
```

## 🔄 Development Workflow

### Phase 1: Core Features (MVP)
- [ ] User authentication & workspace creation
- [ ] Basic request builder (GET, POST) with Electric Blue UI
- [ ] Response viewer with JSON formatting
- [ ] Save requests to collections
- [ ] Environment variables
- [ ] Basic WebSocket connection setup
- [ ] **Volt NestJS Package** - Auto-import NestJS APIs

### Phase 2: Real-time Collaboration & Auto-Import
- [ ] Live collection sync via WebSocket
- [ ] Multi-user presence indicators
- [ ] Real-time request execution sharing
- [ ] Basic commenting system
- [ ] Live cursor tracking
- [ ] Conflict resolution for simultaneous edits
- [ ] **Auto-sync dashboard integration**
- [ ] **CLI tools for manual sync**

### Phase 3: Advanced Features & Framework Expansion
- [ ] Full HTTP method support
- [ ] Request history with live updates
- [ ] Import/export collections
- [ ] Advanced team collaboration features
- [ ] API testing & assertions
- [ ] Code generation with live preview
- [ ] **Express.js & Fastify packages**
- [ ] **Go Gin framework support**

## 🧪 Testing Strategy

### Backend Testing
```bash
go test ./...
go test -race ./...
go test -cover ./...
```

### Frontend Testing
```bash
npm run test
npm run test:e2e
```

### Integration Testing
- API endpoint testing with Postman/Newman
- End-to-end testing with Playwright

## 📈 Performance Optimizations

### Frontend
- Lazy loading for large collections
- Virtual scrolling for history
- Code splitting by routes
- Response caching

### Backend
- MongoDB indexing strategy
- Request/response compression
- Rate limiting
- Connection pooling

## 🔐 Security Considerations

- JWT token validation
- Input sanitization
- Rate limiting per user
- CORS configuration
- Password hashing with bcrypt
- Request size limits

## 📚 Resources

- [Svelte Documentation](https://svelte.dev/docs)
- [SvelteKit Documentation](https://kit.svelte.dev/docs)
- [Gin Framework](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [WebSocket in Go](https://pkg.go.dev/github.com/gorilla/websocket)

---

**Happy coding! 🎉**