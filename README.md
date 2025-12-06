# UMess - WhatsApp-like Messaging System

A real-time messaging application built with Go backend and React frontend, designed to handle 700,000 users efficiently.

## Features

- ✅ User Authentication (Register/Login with JWT)
- ✅ Real-time Messaging via WebSocket
- ✅ One-to-one Chat
- ✅ Image Upload & Sharing
- ✅ Message Status (Sent, Delivered, Read)
- ✅ Conversation List
- ✅ Message History with Pagination
- ✅ Typing Indicators (ready for implementation)
- ✅ Modern, Responsive UI

## Tech Stack

### Backend
- **Go (Golang)** - High-performance backend
- **PostgreSQL** - Relational database
- **WebSocket** (gorilla/websocket) - Real-time communication
- **JWT** - Authentication
- **Gorilla Mux** - HTTP router

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **React Router** - Routing
- **Axios** - HTTP client
- **TanStack Query** - Data fetching

## Project Structure

```
umess/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── handler/         # HTTP/WebSocket handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Database layer
│   │   ├── model/           # Data models
│   │   ├── middleware/      # Auth, CORS middleware
│   │   └── websocket/       # WebSocket hub & clients
│   ├── pkg/
│   │   ├── auth/            # JWT authentication
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   └── utils/           # Utilities
│   ├── migrations/          # Database migrations
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/           # Page components
│   │   ├── services/        # API services
│   │   └── hooks/           # Custom hooks
│   └── package.json
├── docker-compose.yml       # Local development setup
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm/yarn
- PostgreSQL 15+
- Docker & Docker Compose (optional, for easy setup)

## Installation & Setup

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd umess
   ```

2. **Start PostgreSQL with Docker**
   ```bash
   docker-compose up -d postgres
   ```

3. **Run database migrations**
   ```bash
   # Connect to PostgreSQL and run migrations
   docker exec -i umess_postgres psql -U postgres -d umess < backend/migrations/001_initial_schema.sql
   ```

4. **Backend Setup**
   ```bash
   cd backend
   
   # Copy environment file
   cp .env.example .env
   # Edit .env with your settings if needed
   
   # Install dependencies
   go mod download
   
   # Run the server
   go run cmd/server/main.go
   ```

5. **Frontend Setup** (in a new terminal)
   ```bash
   cd frontend
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

6. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - WebSocket: ws://localhost:8080/ws

### Option 2: Manual Setup

1. **Database Setup**
   ```bash
   # Create database
   createdb umess
   
   # Run migrations
   psql -U postgres -d umess -f backend/migrations/001_initial_schema.sql
   ```

2. **Backend Setup**
   ```bash
   cd backend
   
   # Set environment variables
   export DATABASE_URL="postgres://user:password@localhost:5432/umess?sslmode=disable"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   
   # Or create .env file
   cp .env.example .env
   # Edit .env with your settings
   
   # Install dependencies
   go mod download
   
   # Run the server
   go run cmd/server/main.go
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Environment Variables

### Backend (.env)
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/umess?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
ENVIRONMENT=development
MEDIA_STORAGE=./uploads
```

### Frontend
Create `.env` file in frontend directory:
```env
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/me` - Get current user (protected)

### Conversations
- `GET /api/conversations` - Get user conversations (protected)
- `POST /api/conversations` - Create direct conversation (protected)
- `GET /api/conversations/{id}` - Get conversation details (protected)

### Messages
- `GET /api/conversations/{id}/messages` - Get messages (protected)
- `POST /api/conversations/{id}/messages` - Send message (protected)
- `POST /api/conversations/{id}/read` - Mark as read (protected)

### Media
- `POST /api/media/upload` - Upload image (protected)
- `GET /api/media/{filename}` - Get media file (protected)

### WebSocket
- `ws://localhost:8080/ws?token={jwt_token}` - WebSocket connection

## Database Schema

- **users** - User accounts
- **conversations** - Chat conversations
- **participants** - Conversation participants
- **messages** - Message content
- **media** - File/media metadata

## Scalability Features

- **Connection Pooling** - Database connection optimization
- **Indexed Queries** - Fast database lookups
- **Pagination** - Efficient message loading
- **Horizontal Scaling Ready** - Stateless backend architecture
- **WebSocket Hub** - Efficient real-time message delivery

## Development

### Running Tests
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests (when implemented)
cd frontend
npm test
```

### Building for Production

**Backend:**
```bash
cd backend
go build -o bin/server cmd/server/main.go
```

**Frontend:**
```bash
cd frontend
npm run build
```

## Production Deployment Considerations

For handling 700,000 users, consider:

1. **Infrastructure:**
   - Load balancer (nginx, AWS ALB)
   - Multiple backend instances
   - Redis for session/pub-sub across servers
   - CDN for media files
   - Database replication

2. **Database:**
   - Connection pooling (already implemented)
   - Read replicas
   - Partitioning for messages table
   - Regular backups

3. **Monitoring:**
   - Application metrics
   - Database performance
   - WebSocket connection count
   - Error tracking

4. **Security:**
   - HTTPS/WSS
   - CORS configuration
   - Rate limiting
   - Input validation

## Budget Feasibility (50 Juta Rupiah)

**Development Phase:** ✅ Complete (Open source stack)

**Infrastructure Costs (Estimated):**
- VPS/Cloud: ~2-5 juta/bulan untuk 700k users
- Database: Managed PostgreSQL atau self-hosted
- Storage: Object storage untuk media
- CDN: Untuk media delivery

**Recommendation:** Start with MVP, then scale infrastructure as user base grows.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.


