# Quick Setup Guide

## Prerequisites Check

Before starting, ensure you have:
- ✅ Go 1.21+ installed
- ✅ Node.js 18+ and npm installed
- ✅ PostgreSQL 15+ installed (or Docker)

## Quick Start (5 minutes)

### 1. Database Setup

**Option A: Using Docker (Easiest)**
```bash
docker-compose up -d postgres
```

Wait 10 seconds for PostgreSQL to start, then run migrations:
```bash
docker exec -i umess_postgres psql -U postgres -d umess < backend/migrations/001_initial_schema.sql
```

**Option B: Manual PostgreSQL**
```bash
createdb umess
psql -U postgres -d umess -f backend/migrations/001_initial_schema.sql
```

### 2. Backend Setup

```bash
cd backend

# Create .env file
cat > .env << EOF
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/umess?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production-$(date +%s)
ENVIRONMENT=development
MEDIA_STORAGE=./uploads
EOF

# Install dependencies
go mod download

# Create uploads directory
mkdir -p uploads

# Run server
go run cmd/server/main.go
```

Backend should be running on `http://localhost:8080`

### 3. Frontend Setup (New Terminal)

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

Frontend should be running on `http://localhost:3000`

## Testing the Application

1. Open browser to `http://localhost:3000`
2. Register a new account
3. Open another browser/incognito window
4. Register another account
5. Start chatting!

## Troubleshooting

### Database Connection Error
- Check PostgreSQL is running: `pg_isready` or `docker ps`
- Verify DATABASE_URL in `.env` matches your PostgreSQL setup
- Check PostgreSQL logs: `docker logs umess_postgres`

### Port Already in Use
- Change PORT in backend/.env
- Update VITE_API_URL in frontend/.env if needed

### WebSocket Connection Failed
- Ensure backend is running
- Check browser console for errors
- Verify token is being sent correctly

## Next Steps

- Read [README.md](README.md) for full documentation
- Check [backend/README.md](backend/README.md) for backend details
- Check [frontend/README.md](frontend/README.md) for frontend details


