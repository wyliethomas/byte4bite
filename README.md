# Byte4Bite

A free and open platform for community pantries of all shapes and sizes. Byte4Bite is designed to allow community pantries to manage and maintain inventory in the most easiest and frictionless way possible.

## Features

- **Admin Dashboard** - Manage inventory, fulfill orders, and track donations
- **User Portal** - Browse items, create carts, and track orders
- **Multi-Pantry Support** - Support multiple community pantries on one platform
- **Order Fulfillment** - Complete workflow from cart to pickup
- **Notifications** - Email and SMS notifications for order updates
- **Analytics & Reporting** - Track usage, inventory trends, and impact metrics
- **Donation Tracking** - Manage donors and track contributions

## Technology Stack

### Backend
- **Go** - High-performance backend with Gin framework
- **PostgreSQL** - Robust relational database
- **GORM** - Go ORM for database operations
- **JWT** - Secure authentication

### Frontend
- **React** - Modern UI framework
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool
- **React Router** - Client-side routing
- **Axios** - HTTP client

### Deployment
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Standalone Binaries** - Easy deployment on Windows/Linux

## Getting Started

### Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **PostgreSQL 16+** - [Download](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional) - [Download](https://www.docker.com/)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone https://github.com/byte4bite/byte4bite.git
   cd byte4bite
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health

### Local Development Setup

#### Backend Setup

1. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

2. **Start PostgreSQL** (using Docker)
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run the backend**
   ```bash
   go run cmd/server/main.go
   ```

   The backend will be available at http://localhost:8080

#### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Install dependencies**
   ```bash
   npm install
   ```

4. **Start development server**
   ```bash
   npm run dev
   ```

   The frontend will be available at http://localhost:5173

## Project Structure

```
byte4bite/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── api/              # API handlers and routes
│   ├── auth/             # Authentication logic
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and migrations
│   ├── models/           # Database models
│   ├── repositories/     # Data access layer
│   ├── services/         # Business logic
│   └── utils/            # Utility functions
├── frontend/
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API services
│   │   ├── context/      # React context
│   │   ├── hooks/        # Custom hooks
│   │   ├── types/        # TypeScript types
│   │   └── utils/        # Utility functions
│   └── public/           # Static assets
├── migrations/           # Database migrations
├── scripts/              # Utility scripts
├── docker-compose.yml    # Production Docker setup
└── docker-compose.dev.yml # Development Docker setup
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout user

### Public Routes
- `GET /api/v1/pantries` - List all pantries
- `GET /api/v1/items` - List available items

### User Routes (Authentication Required)
- `GET /api/v1/users/me` - Get current user
- `GET /api/v1/carts/current` - Get current cart
- `POST /api/v1/carts/items` - Add item to cart
- `GET /api/v1/users/orders` - Get user orders

### Admin Routes (Admin Role Required)
- `GET /api/v1/admin/dashboard` - Admin dashboard
- `POST /api/v1/admin/items` - Create item
- `GET /api/v1/admin/orders` - Manage orders
- `POST /api/v1/admin/categories` - Create category

See [API Documentation](docs/api.md) for complete endpoint list.

## Database Schema

The application uses PostgreSQL with the following main tables:
- **users** - User accounts and authentication
- **pantries** - Community pantry information
- **categories** - Item categories
- **items** - Inventory items
- **carts** - User shopping carts
- **cart_items** - Items in carts
- **orders** - Submitted orders
- **donations** - Donation tracking
- **notifications** - Email/SMS notifications

See [Database Schema](docs/schema.md) for detailed information.

## Development

### Running Tests
```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm test
```

### Building for Production

#### Backend Binary
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o byte4bite-linux-amd64 ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o byte4bite-windows-amd64.exe ./cmd/server
```

#### Frontend
```bash
cd frontend
npm run build
```

### Docker Production Build
```bash
docker-compose build
docker-compose up -d
```

## Configuration

### Environment Variables

#### Backend
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/production)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT signing secret (REQUIRED in production)
- `JWT_EXPIRY_HOURS` - Token expiry time in hours

#### Frontend
- `VITE_API_URL` - Backend API URL

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/byte4bite/byte4bite/issues)
- **Discussions**: [GitHub Discussions](https://github.com/byte4bite/byte4bite/discussions)

## Roadmap

See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for the complete development roadmap.

### Current Status
✅ Phase 1: Project Foundation (Completed)
- Go backend structure
- React frontend with TypeScript
- PostgreSQL database setup
- Docker configuration
- Environment setup

### Next Steps
🚧 Phase 2: Authentication & User Management
- User registration and login
- JWT authentication
- Role-based access control
- User profile management

## Acknowledgments

Built with ❤️ for community pantries everywhere.
