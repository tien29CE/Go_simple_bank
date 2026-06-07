# Simple Bank - Go Backend API

A comprehensive banking system backend built with Go, featuring REST API (Gin), gRPC, PostgreSQL database, JWT/PASETO authentication, and comprehensive testing.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Complete Setup Instructions](#complete-setup-instructions)
4. [Project Structure](#project-structure)
5. [Running the Application](#running-the-application)
6. [Development Workflow](#development-workflow)
7. [Testing](#testing)
8. [Docker Deployment](#docker-deployment)

---

## Project Overview

This project is a complete banking system backend with the following features:
- **REST API** built with Gin framework
- **gRPC Services** for high-performance communication
- **PostgreSQL** database with migrations
- **Authentication** using JWT and PASETO tokens
- **Mock Database** for unit testing
- **Swagger Documentation** for API endpoints
- **Proto Buffer** for service definitions

---

## Prerequisites

You need to have the following installed on your system:

### Required Tools
- **Go 1.26.1 or higher** - [Download Go](https://go.dev/doc/install)
- **PostgreSQL** (or Docker to run PostgreSQL in container)
- **Docker** & **Docker Compose** (recommended for database)
- **Git** - For version control

### System Check

```bash
# Verify Go installation
go version

# Verify Docker installation
docker --version
docker-compose --version
```

---

## Complete Setup Instructions

### Step 1: Clone the Repository

```bash
# Clone the project
git clone https://github.com/tien29CE/Go_simple_bank.git
cd Go_simple_bank
```

### Step 2: Install Required Go Tools

This step installs all the necessary tools for code generation, database migration, and testing.

#### 2.1 Install golang-migrate CLI (for database migrations)

```bash
# On Linux/Mac
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

# Verify installation
migrate -version
```

#### 2.2 Create migrate up or down schema (for database migrations)

```bash
# On Linux/Mac
# migrate create -ext <extension_file> -dir <relative/path/to/sql/folder> -seq <name_of_the_sequence>
migrate create -ext sql -dir db/migration -seq init_schema
```

#### 2.2 Install SQLc (for generating Go code from SQL)

SQLc generates type-safe Go code from SQL queries. This eliminates the need for manual SQL parsing.

```bash
# On Linux/Mac
# Visit: https://docs.sqlc.dev/en/v1.18.0/overview/install.html
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

# Verify installation
sqlc version
```

#### 2.3 Install Go Mock (mockgen)

mockgen generates mock implementations of interfaces for testing.

```bash
# Install mockgen CLI
go install github.com/golang/mock/cmd/mockgen@latest


# This installs mockgen to $GOPATH/bin
# Make sure $GOPATH/bin is in your PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" || echo "Add $(go env GOPATH)/bin to your PATH"

# Verify installation
mockgen -version
```

#### 2.4 Install Protocol Buffer Compiler

The protobuf compiler is required to generate Go code from .proto files.

```bash
# On Ubuntu/Debian
sudo apt update
sudo apt install -y protobuf-compiler

# On Mac (using Homebrew)
brew install protobuf

# Verify installation
protoc --version
```

#### 2.5 Install gRPC and gRPC Gateway Plugins

These plugins generate gRPC and REST gateway code.

```bash
# Install protoc-gen-go (generates Go code from .proto)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install protoc-gen-go-grpc (generates gRPC service code)
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install protoc-gen-grpc-gateway (generates REST gateway code)
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Install protoc-gen-openapiv2 (generates Swagger/OpenAPI documentation)
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Verify all are installed
which protoc-gen-go protoc-gen-go-grpc protoc-gen-grpc-gateway protoc-gen-openapiv2
```

#### 2.6 (Optional) Install Evans for gRPC Testing

Evans is a CLI tool for testing gRPC services.

```bash
# Download and install Evans
go install github.com/ktr0731/evans@latest

# Verify installation
evans --version
```

### Step 3: Download Google Protocol Buffer Files

The proto definitions use Google's standard library for HTTP annotations and RPC definitions. You need to copy these files to your project.

These files are from: https://github.com/googleapis/googleapis

#### 3.1 Create the proto/google directory structure

```bash
# Navigate to your project
cd /path/to/Go_simple_bank

# Create the google/api directory if it doesn't exist
mkdir -p proto/google/api
```

#### 3.2 Download the 4 required Google proto files

You have two options:

**Option A: Manual Download (Recommended for understanding)**

1. Visit https://github.com/googleapis/googleapis/tree/master/google/api
2. Download these 4 files and place them in `proto/google/api/`:
   - `annotations.proto` - Defines HTTP binding annotations for gRPC services
   - `field_behavior.proto` - Defines field behavior annotations
   - `http.proto` - Defines HTTP request/response mappings
   - `httpbody.proto` - Defines the HTTP request/response body types

**Option B: Using curl (Quick)**

```bash
cd proto/google/api

# Download each file
curl -O https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
curl -O https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto
curl -O https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
curl -O https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto

# Verify files are downloaded
ls -la
```

**Option C: Using a script**

```bash
#!/bin/bash
PROTO_DIR="proto/google/api"
GITHUB_URL="https://raw.githubusercontent.com/googleapis/googleapis/master/google/api"

mkdir -p "$PROTO_DIR"

files=("annotations.proto" "field_behavior.proto" "http.proto" "httpbody.proto")

for file in "${files[@]}"; do
  echo "Downloading $file..."
  curl -o "$PROTO_DIR/$file" "$GITHUB_URL/$file"
done

echo "Done! Files downloaded to $PROTO_DIR"
```

#### 3.3 Verify the files

```bash
# Check that all 4 files exist
ls -la proto/google/api/
# You should see:
# - annotations.proto
# - field_behavior.proto
# - http.proto
# - httpbody.proto
```

### Step 4: Set Up Environment Variables

The application reads configuration from `app.env` file.

```bash
# Copy and review the example configuration
cat > app.env << 'EOF'
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:root@localhost:5432/simple_bank?sslmode=disable
MIGRATION_URL=file://db/migration
HTTP_SERVER_ADDRESS=0.0.0.0:8080
GRPC_SERVER_ADDRESS=0.0.0.0:9090
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
EOF
```

**Configuration Explanation:**
- `DB_DRIVER`: PostgreSQL driver
- `DB_SOURCE`: PostgreSQL connection string (user:password@host:port/dbname)
- `MIGRATION_URL`: Path to database migrations
- `HTTP_SERVER_ADDRESS`: REST API server address and port
- `GRPC_SERVER_ADDRESS`: gRPC server address and port
- `TOKEN_SYMMETRIC_KEY`: 32-character key for JWT/PASETO token generation (keep this secret!)
- `ACCESS_TOKEN_DURATION`: How long access tokens are valid (15 minutes)
- `REFRESH_TOKEN_DURATION`: How long refresh tokens are valid (24 hours)

### Step 5: Generate SQLc Code

SQLc generates type-safe Go code from your SQL queries.

#### 5.1 Generate SQL code

```bash
# From project root directory
make sqlc

# Or manually run:
sqlc generate
```

#### 5.2 What happens:
- SQLc reads SQL queries from `db/query/*.sql` files
- SQLc reads database schema from `db/migration/*.sql` files
- Generates Go code in `db/sqlc/` directory with:
  - Type-safe query functions
  - Model structs for database rows
  - Query interfaces

#### 5.3 Check generated files:
```bash
# Generated files should appear in db/sqlc/
ls -la db/sqlc/
# You should see:
# - account.sql.go
# - entry.sql.go
# - transfer.sql.go
# - user.sql.go
# - session.sql.go
# - models.go (database models)
# - querier.go (query interface)
# - store.go (store interface)
```

### Step 6: Set Up PostgreSQL Database

You can run PostgreSQL using Docker (recommended) or use a local installation.

#### 6.1 Start PostgreSQL Container

```bash
# Create a Docker network for better communication between containers
sudo docker network create bank-network

# Start PostgreSQL container
make postgres

# Or manually:
sudo docker run --name postgres18.3 \
  --network bank-network \
  -p 5432:5432 \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -d postgres:18.3-alpine

# Verify container is running
sudo docker ps | grep postgres18.3
```

#### 6.2 Create the Database

```bash
# Create the simple_bank database
make createdb

# Or manually:
sudo docker exec -it postgres18.3 createdb --username=root --owner=root simple_bank
```

#### 6.3 Verify Database Connection

```bash
# Connect to the database to verify it works
sudo docker exec -it postgres18.3 psql -U root -d simple_bank

# In psql shell, you can run:
\dt              # List tables (should be empty at first)
\l               # List databases
\q               # Quit psql
```

### Step 7: Run Database Migrations

Migrations create the database schema.

#### 7.1 Run all migrations

```bash
# Run all pending migrations
make migrateup

# Or manually:
migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

#### 7.2 Verify migration was successful

```bash
# Connect to database and check tables
sudo docker exec -it postgres18.3 psql -U root -d simple_bank

# In psql shell:
\dt                    # Should show: accounts, entries, transfers, users, sessions
\d accounts            # Describe the accounts table
\q                     # Quit
```

#### 7.3 Understanding Migrations

```bash
# View migration files
ls -la db/migration/

# Migrations are applied in order:
# 1. 000001_init_schema.up.sql - Creates accounts, entries, transfers tables
# 2. 000002_add_users.up.sql - Creates users table
# 3. 000003_add_sessions.up.sql - Creates sessions table

# Useful migration commands:
make migratedown      # Revert all migrations
make migrateup1       # Apply only 1 migration
make migratedown1     # Revert 1 migration
```

### Step 8: Download Go Dependencies

```bash
# Download all Go module dependencies
go mod tidy

# Verify dependencies are downloaded
go mod verify
```

### Step 9: Generate Mock Database

Create mock database interfaces for unit testing.

#### 9.1 Generate mock

```bash
# Generate mock database implementation
make mock

# Or manually:
mockgen -package mockdb \
  -destination db/mock/store.go \
  github.com/tien29CE/Go_simple_bank.git/db/sqlc Store
```

#### 9.2 Verify mock was generated

```bash
# Check generated mock file
ls -la db/mock/store.go

# The file should contain mock implementations of all database operations
```

### Step 10: Generate gRPC Code from Proto Files

Generate Go code from Protocol Buffer definitions.

#### 10.1 Generate gRPC code

```bash
# Generate all protobuf code (Go, gRPC, gateway, OpenAPI)
make proto

# Or manually:
rm -f pb/*.go
rm -f doc/swagger/*.swagger.json
protoc \
  --experimental_allow_proto3_optional \
  --proto_path=proto \
  --go_out=pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=pb \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pb \
  --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=doc/swagger \
  --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
  proto/*.proto
```

#### 10.2 What gets generated:

```bash
# Generated Go files for gRPC
ls -la pb/
# Should contain:
# - service_simple_bank.pb.go (service message definitions)
# - service_simple_bank_grpc.pb.go (gRPC service interfaces)
# - service_simple_bank.pb.gw.go (REST gateway)
# - user.pb.go, rpc_*.pb.go (message definitions)

# Generated Swagger/OpenAPI documentation
ls -la doc/swagger/
# Should contain: simple_bank.swagger.json
```

#### 10.3 Proto Explanation

```
proto/
  ├── service_simple_bank.proto     # Main service definition
  ├── user.proto                    # User message definition
  ├── rpc_create_user.proto         # CreateUser RPC definition
  ├── rpc_login_user.proto          # LoginUser RPC definition
  ├── rpc_update_user.proto         # UpdateUser RPC definition
  └── google/api/                   # Google standard protos
      ├── annotations.proto         # HTTP binding annotations
      ├── field_behavior.proto      # Field behavior annotations
      ├── http.proto                # HTTP request/response mappings
      └── httpbody.proto            # HTTP body definitions
```

---

## Project Structure

```
Go_simple_bank/
├── api/                           # REST API handlers (Gin)
│   ├── account.go                 # Account endpoints
│   ├── account_test.go            # Account tests
│   ├── user.go                    # User endpoints
│   ├── user_test.go               # User tests
│   ├── middleware.go              # JWT/PASETO auth middleware
│   ├── token.go                   # Token endpoints
│   ├── server.go                  # Gin server setup
│   └── validator.go               # Request validation
│
├── db/                            # Database layer
│   ├── migration/                 # SQL migrations
│   │   ├── 000001_init_schema.*   # Schema creation
│   │   ├── 000002_add_users.*     # Add users table
│   │   └── 000003_add_sessions.*  # Add sessions table
│   ├── query/                     # SQL query files
│   │   ├── account.sql            # Account queries
│   │   ├── user.sql               # User queries
│   │   ├── entry.sql              # Entry queries
│   │   ├── transfer.sql           # Transfer queries
│   │   └── session.sql            # Session queries
│   ├── sqlc/                      # Generated SQLc code
│   │   ├── *.sql.go               # Generated query functions
│   │   ├── models.go              # Database models
│   │   └── store.go               # Store interface
│   └── mock/                      # Mock database (generated)
│       └── store.go               # Mock implementation
│
├── gapi/                          # gRPC API handlers
│   ├── server.go                  # gRPC server setup
│   ├── rpc_create_user.go         # CreateUser RPC handler
│   ├── rpc_login_user.go          # LoginUser RPC handler
│   ├── rpc_update_user.go         # UpdateUser RPC handler
│   ├── authorization.go           # gRPC auth logic
│   ├── error.go                   # gRPC error handling
│   ├── converter.go               # Proto ↔ SQL model conversion
│   └── metadata.go                # gRPC metadata extraction
│
├── pb/                            # Generated gRPC code
│   ├── *.pb.go                    # Generated message types
│   ├── *_grpc.pb.go               # Generated service interfaces
│   └── *.pb.gw.go                 # Generated REST gateway
│
├── proto/                         # Protocol Buffer definitions
│   ├── service_simple_bank.proto  # Service definition
│   ├── user.proto                 # User messages
│   ├── rpc_*.proto                # RPC definitions
│   └── google/api/                # Google standard library
│
├── token/                         # Token generation
│   ├── jwt_maker.go               # JWT token implementation
│   ├── jwt_maker_test.go          # JWT tests
│   ├── paseto_maker.go            # PASETO token implementation
│   ├── paseto_maker_test.go       # PASETO tests
│   ├── maker.go                   # Token interface
│   └── payload.go                 # Token payload
│
├── util/                          # Utilities
│   ├── config.go                  # Config loader
│   ├── currency.go                # Currency validation
│   ├── password.go                # Password hashing
│   ├── password_test.go           # Password tests
│   └── random.go                  # Random data generation
│
├── val/                           # Custom validators
│   └── validator.go               # Currency validator
│
├── tools/                         # Build tools
│   └── tools.go                   # Tool dependencies
│
├── doc/                           # Documentation
│   └── swagger/                   # Auto-generated Swagger/OpenAPI
│       └── simple_bank.swagger.json
│
├── main.go                        # Application entry point
├── app.env                        # Environment configuration
├── Makefile                       # Build targets
├── Dockerfile                     # Container image definition
├── docker-compose.yaml            # Multi-container setup
├── sqlc.yaml                      # SQLc configuration
├── go.mod                         # Go module definition
└── go.sum                         # Go module checksums
```

---

## Running the Application

### Option 1: Run Locally (Development)

#### Prerequisites
- PostgreSQL must be running (see Step 6)
- All migrations must be run (see Step 7)
- All code generation must be complete (see Steps 5, 9, 10)

#### Start the server

```bash
# Run the application
make server

# Or manually:
go run main.go

# You should see:
# - HTTP server listening on 0.0.0.0:8080
# - gRPC server listening on 0.0.0.0:9090
```

#### Test the server

```bash
# In another terminal, test the REST API
curl http://localhost:8080/health

# Or test gRPC using evans CLI
evans --host localhost --port 9090 -r repl

# In evans shell:
call SimpleBank.CreateUser
```

### Option 2: Run with Docker Compose (Recommended)

Docker Compose automatically handles PostgreSQL, migrations, and the application server.

```bash
# Start all services
docker-compose up --build

# Or in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

The application will automatically:
1. Start PostgreSQL database
2. Wait for database to be ready
3. Run database migrations
4. Start the gRPC and HTTP servers

---

## Development Workflow

### When you modify code:

#### 1. If you modify SQL queries

```bash
# Regenerate Go code from SQL
make sqlc

# The generated code in db/sqlc/ will update
```

#### 2. If you modify .proto files

```bash
# Regenerate gRPC code
make proto

# The generated code in pb/ will update
# Swagger docs in doc/swagger/ will update
```

#### 3. If you modify the database interface (Store interface)

```bash
# Regenerate mock database
make mock

# Updated mock in db/mock/store.go
```

#### 4. Complete regeneration

```bash
# Regenerate everything
make sqlc && make mock && make proto
```

### Adding new API endpoints

#### REST API (Gin)

1. Create handler in `api/` directory
2. Add route in `api/server.go`
3. Create tests in `api/*_test.go`
4. Run tests: `make test`

#### gRPC API

1. Define RPC in `.proto` file (proto/rpc_*.proto)
2. Add proto to `proto/service_simple_bank.proto`
3. Run `make proto` to generate code
4. Implement handler in `gapi/rpc_*.go`
5. Create tests
6. Run tests: `make test`

### Adding database migrations

When you need to change the database schema:

```bash
# 1. Create migration files
touch db/migration/000004_your_migration_name.up.sql
touch db/migration/000004_your_migration_name.down.sql

# 2. Write SQL in the .up.sql file (schema changes)
# 3. Write SQL in the .down.sql file (rollback instructions)

# 4. Apply the migration
make migrateup

# 5. Regenerate SQLc code
make sqlc

# 6. Create SQL queries for your new tables in db/query/

# 7. Regenerate SQLc code again
make sqlc

# 8. Use the new generated functions in your handlers
```

---

## Testing

### Run all tests

```bash
# Run all tests with coverage
make test

# Or manually:
go test -v -cover ./...

# Run tests for specific package
go test -v -cover ./api
go test -v -cover ./gapi
go test -v -cover ./db/sqlc
```

### Run specific test

```bash
# Run a specific test function
go test -run TestFunctionName -v ./api
```

### Testing with database

Tests that interact with the database require:
1. PostgreSQL running with test database
2. Migrations applied to test database
3. Or use mock database (db/mock/store.go)

---

## Docker Deployment

### Build Docker Image

```bash
# Build image
docker build -t simplebank:latest .

# Or using docker-compose
docker-compose build
```

### Run with Docker

```bash
# Run just the application
docker run -p 8080:8080 -p 9090:9090 \
  -e DB_SOURCE="postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" \
  simplebank:latest

# Or use docker-compose (recommended)
docker-compose up -d
```

### Clean Up

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (deletes database)
docker-compose down -v

# Remove database container
sudo docker stop postgres18.3
sudo docker rm postgres18.3

# Remove network
sudo docker network rm bank-network
```

---

## Troubleshooting

### "cannot load config"
- Verify `app.env` exists in project root
- Check environment variables are correctly set

### "cannot connect to db"
- Verify PostgreSQL is running: `sudo docker ps | grep postgres`
- Verify connection string in `app.env`
- Test connection: `sudo docker exec -it postgres18.3 psql -U root -d simple_bank`

### Migration errors
- Ensure `migrate` CLI is installed
- Verify migration files exist in `db/migration/`
- Check connection string is correct
- Rollback: `make migratedown`

### Code generation errors (sqlc/proto)

```bash
# sqlc errors
# - Ensure database schema matches queries
# - Run migrations first
# - Check sqlc.yaml configuration

# Proto errors
# - Verify protobuf compiler is installed: protoc --version
# - Check all proto files are valid
# - Ensure google proto files are in proto/google/api/
```

### Port already in use
```bash
# Kill process on port 8080
lsof -i :8080 | grep LISTEN | awk '{print $2}' | xargs kill -9

# Kill process on port 9090
lsof -i :9090 | grep LISTEN | awk '{print $2}' | xargs kill -9

# Or use docker-compose
docker-compose down
```

### Mock generation errors
- Ensure the Store interface is exported (starts with capital letter)
- Verify the package path is correct in mockgen command

---

## Key Concepts

### SQLc (SQL Code Generation)
- Generates type-safe Go code from SQL queries
- Eliminates runtime SQL parsing errors
- Reduces boilerplate
- SQL queries must match database schema exactly

### gRPC
- High-performance RPC framework using Protocol Buffers
- Supports streaming and bidirectional communication
- Runs on HTTP/2
- REST gateway translates HTTP requests to gRPC calls

### Migrations
- Version control for database schema
- Applied sequentially with timestamps
- Can be rolled back
- Keeps database schema in sync across environments

### Protocol Buffers (Proto)
- Language and platform agnostic serialization format
- Defines service contracts and messages
- Compiled to language-specific code
- More compact and faster than JSON

### Mock Testing
- mockgen generates mock implementations
- Allows testing without real database
- Speeds up unit tests
- Enables testing error conditions

---

## Useful Commands Reference

```bash
# Database
make postgres          # Start PostgreSQL container
make createdb         # Create database
make dropdb           # Drop database
make migrateup        # Run all migrations
make migratedown      # Rollback all migrations
make migrateup1       # Run 1 migration
make migratedown1     # Rollback 1 migration

# Code generation
make sqlc             # Generate SQLc code
make mock             # Generate mock database
make proto            # Generate gRPC code

# Testing and running
make test             # Run all tests
make server           # Run application locally

# gRPC testing
make evans            # Start Evans CLI for gRPC testing

# Cleaning
docker-compose down   # Stop and remove containers
docker system prune   # Clean up unused Docker resources
```

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Documentation](https://grpc.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [SQLc Documentation](https://sqlc.dev/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Gin Web Framework](https://gin-gonic.com/)

---

## License

This project is open source and available under the MIT License.

---

**Last Updated:** June 2026  
**Project Version:** 1.0.0
