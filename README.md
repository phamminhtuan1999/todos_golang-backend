# todos_golang-backend

A RESTful Todo API backend built with Go and PostgreSQL. This is a learning project demonstrating modern Go application development with features like JWT authentication, database migrations, and RESTful API design patterns.

## Purpose

A RESTful Todo API built with Go, PostgreSQL, and JWT authentication. This project demonstrates:

- **CRUD Operations**: Create, read, update, and delete todos
- **User Authentication**: JWT-based user authentication and authorization
- **Database Management**: PostgreSQL with migration support
- **REST API**: Gin framework for HTTP routing and handling

## Features

- User authentication with JWT tokens
- Create, retrieve, update, and delete todos
- User-specific todo isolation
- Database migrations with golang-migrate
- Hot reload development with Air

## Installation Guide

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 13 or higher (or Docker)
- golang-migrate CLI
- Air CLI (optional, for hot reload)

### Setup Steps

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd todos
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up PostgreSQL (Docker)**

   ```bash
   docker run -d \
     --name postgres-db \
     -e POSTGRES_USER=admin \
     -e POSTGRES_PASSWORD=your-secure-password \
     -e POSTGRES_DB=your-db-name \
     -p 5432:5432 \
     postgres:16
   ```

4. **Configure environment variables**
   Create a `.env` file in the project root:

   ```env
   DATABASE_URL=postgresql://admin:your-secure-password@127.0.0.1:5432/your-db-name?sslmode=disable
   PORT=8080
   JWT_SECRET=your-secret-jwt-key
   ```

5. **Run database migrations**

   ```bash
   .\scripts\migrate.ps1 up
   ```

   Or manually:

   ```bash
   migrate -path migrations -database "postgresql://admin:your-secure-password@127.0.0.1:5432/your-db-name?sslmode=disable" up
   ```

6. **Start the server**

   ```bash
   go run cmd/api/main.go
   ```

   Or with hot reload (Air):

   ```bash
   air
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

- `GET /ping` - Health check
- `POST /todos` - Create a todo
- `GET /todos` - Get all todos for the authenticated user
- `GET /todos/:id` - Get a specific todo
- `PUT /todos/:id` - Update a todo
- `DELETE /todos/:id` - Delete a todo
- `POST /users/register` - Register a new user
- `POST /users/login` - Login and get JWT token

## Project Structure

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models
│   └── repository/      # Database operations
├── migrations/          # Database migrations
├── scripts/             # Helper scripts (migrate.ps1)
└── go.mod              # Go module dependencies
```
