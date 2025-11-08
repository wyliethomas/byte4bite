# Byte4Bite - Quick Start Guide

## What We've Built (Phase 1 Complete!)

Phase 1 of the Byte4Bite project is now complete! Here's what's been set up:

### ✅ Backend (Go)
- Complete project structure with clean architecture
- Database models for all core entities:
  - Users (with role-based access)
  - Pantries
  - Categories
  - Items (inventory)
  - Carts and Cart Items
  - Orders
  - Donations
  - Notifications
- Configuration management with environment variables
- Database connection with GORM
- Auto-migration system
- RESTful API structure with Gin framework
- Health check endpoint

### ✅ Frontend (React + TypeScript)
- Vite-powered React application
- TypeScript type definitions for all models
- API service layer with Axios
- JWT authentication interceptors
- Organized project structure
- Ready for routing with React Router

### ✅ Infrastructure
- Docker setup for production deployment
- Docker Compose for local development
- PostgreSQL database configuration
- Nginx configuration for frontend
- Environment variable management
- .gitignore for security

### ✅ Documentation
- Comprehensive README
- Implementation plan
- Environment configuration examples

## Running the Application

### Option 1: Docker Compose (Easiest)

Start everything with one command:

```bash
# Start database only (for local development)
docker-compose -f docker-compose.dev.yml up -d

# Or start everything (full production-like setup)
docker-compose up -d
```

### Option 2: Local Development (Recommended for Active Development)

**Terminal 1 - Database:**
```bash
docker-compose -f docker-compose.dev.yml up
```

**Terminal 2 - Backend:**
```bash
go run cmd/server/main.go
```

**Terminal 3 - Frontend:**
```bash
cd frontend
npm run dev
```

### Access Points

Once running:
- **Frontend**: http://localhost:5173 (dev) or http://localhost:3000 (docker)
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Database**: localhost:5432

## Testing the Setup

### 1. Check Backend Health
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Byte4Bite API is running"
}
```

### 2. Check Database Connection
The backend will automatically:
- Connect to PostgreSQL
- Create all tables via auto-migration
- Display migration status in logs

### 3. Verify Frontend
Open http://localhost:5173 in your browser to see the React app.

## Project Structure Overview

```
byte4bite/
├── cmd/server/           # Backend entry point
├── internal/
│   ├── api/             # API handlers, middleware, routes
│   ├── config/          # Configuration management
│   ├── database/        # DB connection & migrations
│   ├── models/          # Database models (User, Item, Cart, etc.)
│   └── ...
├── frontend/
│   └── src/
│       ├── components/  # React components (to be built)
│       ├── pages/       # Page components (to be built)
│       ├── services/    # API client setup
│       └── types/       # TypeScript definitions
└── ...
```

## Next Steps (Phase 2)

Phase 2 will implement **Authentication & User Management**:

1. User registration endpoint
2. Login with JWT tokens
3. Password hashing with bcrypt
4. Protected routes with middleware
5. User profile management
6. Password reset functionality

## Troubleshooting

### Backend won't start
- Check if PostgreSQL is running: `docker-compose -f docker-compose.dev.yml ps`
- Verify environment variables in `.env`
- Check database credentials

### Frontend build errors
- Delete `node_modules` and reinstall: `cd frontend && rm -rf node_modules && npm install`
- Clear Vite cache: `cd frontend && rm -rf node_modules/.vite`

### Database connection errors
- Ensure PostgreSQL container is running
- Check `DB_HOST` in `.env` (should be `localhost` for local dev)
- Verify port 5432 is not in use by another process

### Port already in use
- Backend (8080): Change `SERVER_PORT` in `.env`
- Frontend (5173): Vite will auto-increment to 5174
- Database (5432): Change port mapping in `docker-compose.dev.yml`

## Available Commands

### Backend
```bash
# Run server
go run cmd/server/main.go

# Build binary
go build -o bin/byte4bite ./cmd/server

# Run tests
go test ./...

# Install dependencies
go mod download
```

### Frontend
```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm test
```

### Docker
```bash
# Start development database
docker-compose -f docker-compose.dev.yml up -d

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild images
docker-compose build
```

## Database Management

### Access PostgreSQL CLI
```bash
docker exec -it byte4bite-db-dev psql -U postgres -d byte4bite
```

### Useful SQL Commands
```sql
-- List all tables
\dt

-- Describe a table
\d users

-- View all users
SELECT * FROM users;

-- Check migrations
SELECT * FROM schema_migrations;
```

## Environment Variables Reference

### Backend (.env)
| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | 8080 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_NAME` | Database name | byte4bite |
| `JWT_SECRET` | JWT signing key | (required in prod) |

### Frontend (.env)
| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | http://localhost:8080/api/v1 |

## Current API Endpoints

Currently implemented (placeholder responses):

- `GET /health` - Health check
- `GET /api/v1/pantries` - List pantries
- `POST /api/v1/auth/register` - Register (coming in Phase 2)
- `POST /api/v1/auth/login` - Login (coming in Phase 2)
- `GET /api/v1/users/me` - Current user (coming in Phase 2)
- `GET /api/v1/carts/current` - Current cart (coming in Phase 4)
- `GET /api/v1/admin/dashboard` - Admin dashboard (coming in Phase 3)

## Support

For issues or questions:
1. Check this guide's Troubleshooting section
2. Review the main [README.md](README.md)
3. See the full [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

---

**Phase 1 Status**: ✅ Complete
**Next Phase**: Phase 2 - Authentication & User Management

Happy coding! 🚀
