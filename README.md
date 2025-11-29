# Golang gRPC Microservice

[![Go Version](https://img.shields.io/badge/Go-1.24.2-00ADD8?style=flat&logo=go)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-1.77.0-244c5a?style=flat&logo=grpc)](https://grpc.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Electronic Commerce microservice built with Go, gRPC, and REST API. Features include user authentication, product management, shopping cart, order processing with Xendit payment integration, and newsletter subscription.

## ✨ Features

### Core Features
- 🔐 **User Authentication & Authorization** - JWT-based authentication with role management
- 🛍️ **Product Management** - CRUD operations for products with image upload
- 🛒 **Shopping Cart** - Add, update, and remove items from cart
- 📦 **Order Processing** - Complete order workflow with status tracking
- 💳 **Payment Integration** - Xendit payment gateway integration
- 📧 **Newsletter Subscription** - Email newsletter management
- 🖼️ **Image Storage** - Product image upload and retrieval

### Technical Features
- ⚡ **High Performance** - Built with gRPC for efficient inter-service communication
- 🔒 **Secure** - JWT authentication, password hashing with bcrypt
- 🐳 **Containerized** - Docker and Docker Compose ready
- 📊 **Database** - PostgreSQL with optimized queries
- 🔄 **Middleware** - Authentication and error handling middleware
- 📝 **Validation** - Request validation using Protocol Buffers
- 🌐 **CORS Enabled** - Cross-Origin Resource Sharing support
- 📦 **Caching** - In-memory caching with go-cache

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

| Layer | Description | Components |
|-------|-------------|------------|
| **Presentation** | API Interfaces | gRPC Server (50052), REST Server (3000) |
| **Handler** | Request Handling | Auth, Product, Cart, Order, Newsletter, Webhook |
| **Service** | Business Logic | Use Cases & Domain Logic |
| **Repository** | Data Access | Database Operations & Queries |
| **Database** | Data Storage | PostgreSQL (Supabase) |



## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.24.2
- **Framework**: gRPC, Fiber v2 (REST)
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Caching**: go-cache
- **Payment**: Xendit SDK
- **Protocol Buffers**: protobuf v1.36.10
- **Validation**: protovalidate

### DevOps
- **Containerization**: Docker, Docker Compose
- **CI/CD**: GitHub Actions
- **Database Hosting**: Supabase

## 📁 Project Structure

```
Golang_GRPC/
├── .github/
│   └── workflows/
│       └── main.yml              # CI/CD pipeline configuration
├── cmd/
│   ├── grpc/
│   │   └── main.go              # gRPC server entry point
│   └── rest/
│       └── main.go              # REST server entry point
├── internal/
│   ├── dto/                     # Data Transfer Objects
│   ├── entity/                  # Domain entities
│   │   ├── jwt/                 # JWT utilities
│   │   ├── cart.go
│   │   ├── newsletter.go
│   │   ├── numbering.go
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── grpcmiddlerware/         # gRPC middleware
│   │   ├── auth_middleware.go
│   │   └── error_middleware.go
│   ├── handler/                 # Request handlers
│   │   ├── auth.go
│   │   ├── cart.go
│   │   ├── newsletter.go
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── product_upload_image.go
│   │   ├── service.go
│   │   └── webhook_handler.go
│   ├── repository/              # Data access layer
│   │   ├── auth_repository.go
│   │   ├── cart_repository.go
│   │   ├── newsletter_repository.go
│   │   ├── order_repository.go
│   │   └── product_repository.go
│   ├── service/                 # Business logic layer
│   │   ├── auth_service.go
│   │   ├── cart_service.go
│   │   ├── newsletter_service.go
│   │   ├── order_service.go
│   │   ├── product_service.go
│   │   └── webhook_service.go
│   └── utils/                   # Utility functions
│       ├── response.go
│       └── validator.go
├── pb/                          # Generated protobuf files
│   ├── auth/
│   ├── cart/
│   ├── common/
│   ├── newsletter/
│   ├── order/
│   ├── product/
│   └── service/
├── pkg/
│   └── database/                # Database connection & queries
│       ├── connection.go
│       └── database_query.go
├── proto/                       # Protocol buffer definitions
│   ├── auth/
│   ├── cart/
│   ├── common/
│   ├── newsletter/
│   ├── order/
│   ├── product/
│   └── service/
├── storage/
│   └── product/                 # Product images storage
├── .dockerignore
├── .env.example                 # Environment variables template
├── .gitignore
├── command_generate_protobuf.txt # Protobuf generation commands
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Multi-stage Docker build
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── LICENSE                      # MIT License
├── README.md                    # This file
└── sql_migration.sql            # Database schema
```

## 📋 Prerequisites

- **Go**: 1.24.2 or higher
- **Docker**: 20.10 or higher
- **Docker Compose**: 2.0 or higher
- **PostgreSQL**: 15 or higher (or Supabase account)
- **Protocol Buffers Compiler**: protoc 3.x

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/arthurhzna/Golang_gRPC.git
cd Golang_gRPC
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Database

#### Option A: Using Supabase (Recommended)

1. Create a Supabase project at [supabase.com](https://supabase.com)
2. Run the SQL migration script in Supabase SQL Editor:
   ```bash
   # Copy content from sql_migration.sql and execute in Supabase
   ```
3. Get your connection string from Supabase Dashboard

#### Option B: Using Local PostgreSQL

```bash
# Create database
createdb golang_grpc

# Run migration
psql -d golang_grpc -f sql_migration.sql
```

### 4. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` file with your credentials:

```env
ENVIRONMENT=DEV

# Supabase Database URL
DB_URL=postgres://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:5432/postgres

# JWT Secret (generate a strong random string)
JWT_SECRET=your_super_secret_jwt_key_here

# Storage Service URL
STORAGE_SERVICE_URL=http://localhost:3000

# Xendit API Key (get from xendit.co)
XENDIT_SECRET_KEY=xnd_development_xxx

# Frontend Base URL
FE_BASE_URL=http://localhost:5173
```

### 5. Generate Protocol Buffers (Optional)

If you modify `.proto` files, regenerate the Go code:

```bash
# Install protoc-gen-go and protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf files (see command_generate_protobuf.txt)
protoc --go_out=./pb --go-grpc_out=./pb --proto_path=./proto \
  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
  auth/auth.proto

# Repeat for other proto files
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `DEV` or `PROD` |
| `DB_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/db` |
| `JWT_SECRET` | Secret key for JWT signing | `your_secret_key` |
| `STORAGE_SERVICE_URL` | Base URL for image storage | `http://localhost:3000` |
| `XENDIT_SECRET_KEY` | Xendit API secret key | `xnd_development_xxx` |
| `FE_BASE_URL` | Frontend application URL | `http://localhost:5173` |

## 🏃 Running the Application

### Development Mode

#### Run gRPC Server

```bash
go run cmd/grpc/main.go
```

Server will start on `localhost:50052`

#### Run REST Server

```bash
go run cmd/rest/main.go
```

Server will start on `localhost:3000`

### Production Mode (Docker)

#### Using Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Using Docker Manually

```bash
# Build image
docker build -t golang-grpc-app .

# Run gRPC server
docker run -d \
  --name grpc-server \
  -p 50052:50052 \
  --env-file .env \
  golang-grpc-app

# Run REST server
docker run -d \
  --name rest-server \
  -p 3000:3000 \
  -v $(pwd)/storage:/app/storage \
  --env-file .env \
  golang-grpc-app /app/rest-server
```

## 📚 API Documentation

### gRPC Services

The application exposes the following gRPC services on port `50052`:

#### Authentication Service
- `Register` - Register new user
- `Login` - User login
- `GetProfile` - Get user profile (requires auth)

#### Product Service
- `CreateProduct` - Create new product (requires auth)
- `GetProduct` - Get product by ID
- `ListProducts` - List all products with pagination
- `UpdateProduct` - Update product (requires auth)
- `DeleteProduct` - Delete product (requires auth)

#### Cart Service
- `AddToCart` - Add item to cart (requires auth)
- `GetCart` - Get user's cart (requires auth)
- `UpdateCartItem` - Update cart item quantity (requires auth)
- `RemoveFromCart` - Remove item from cart (requires auth)

#### Order Service
- `CreateOrder` - Create new order (requires auth)
- `GetOrder` - Get order by ID (requires auth)
- `ListOrders` - List user's orders (requires auth)
- `UpdateOrderStatus` - Update order status (requires auth)

#### Newsletter Service
- `Subscribe` - Subscribe to newsletter
- `Unsubscribe` - Unsubscribe from newsletter

### REST Endpoints

The REST API runs on port `3000`:

| Method | Endpoint | Description | 
|--------|----------|-------------|
| `GET` | `/storage/product/:filename` | Get product image |
| `POST` | `/product/upload` | Upload product image |
| `POST` | `/webhook/xendit/invoice` | Xendit payment webhook |


## 🐳 Docker Deployment

### Multi-Stage Build

The Dockerfile uses multi-stage builds for optimal image size:

1. **Builder Stage**: Compiles Go binaries
2. **Runtime Stage**: Minimal Alpine Linux with only necessary dependencies

### Image Size Optimization

- Base image: `alpine:latest` (~5MB)
- Final image size: ~30MB (including both binaries)
- No source code in final image
- Only runtime dependencies included

### Docker Compose Services

```yaml
services:
  grpc-server:  # gRPC service on port 50052
  rest-server:  # REST API on port 3000
```

Both services share the same Docker image but run different binaries.

## 💻 Development

### Code Structure Guidelines

This project follows **Clean Architecture** principles:

1. **Handler Layer**: Handles HTTP/gRPC requests
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Database operations
4. **Entity Layer**: Domain models

### Adding New Features

1. Define Protocol Buffer in `proto/`
2. Generate Go code: `protoc --go_out=...`
3. Create entity in `internal/entity/`
4. Implement repository in `internal/repository/`
5. Implement service in `internal/service/`
6. Create handler in `internal/handler/`
7. Register service in `cmd/grpc/main.go`

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Commit Message Convention

```
feat: Add new feature
fix: Fix bug
docs: Update documentation
style: Format code
refactor: Refactor code
test: Add tests
chore: Update dependencies
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**blackjokie**

- GitHub: [@arthurhzna](https://github.com/arthurhzna)

## 🙏 Acknowledgments

- [gRPC](https://grpc.io/) - High-performance RPC framework
- [Fiber](https://gofiber.io/) - Express-inspired web framework
- [Supabase](https://supabase.com/) - Open source Firebase alternative
- [Xendit](https://xendit.co/) - Payment gateway for Southeast Asia
- [Protocol Buffers](https://protobuf.dev/) - Language-neutral data serialization

## 📞 Support

For support, email your-email@example.com or open an issue in the GitHub repository.

---

**⭐ If you find this project useful, please consider giving it a star!**

Made with ❤️ using Go and gRPC

