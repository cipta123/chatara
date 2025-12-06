# UMess Backend

Go backend for UMess messaging application.

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Setup environment:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. Run database migrations:
   ```bash
   psql -U postgres -d umess -f migrations/001_initial_schema.sql
   ```

4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## API Documentation

See main README.md for API endpoints.

## Project Structure

- `cmd/server/` - Application entry point
- `internal/` - Private application code
  - `handler/` - HTTP/WebSocket handlers
  - `service/` - Business logic
  - `repository/` - Database access
  - `model/` - Data models
  - `middleware/` - HTTP middleware
  - `websocket/` - WebSocket implementation
- `pkg/` - Public packages
  - `auth/` - JWT authentication
  - `config/` - Configuration management
  - `database/` - Database connection
  - `utils/` - Utility functions
- `migrations/` - Database migrations


