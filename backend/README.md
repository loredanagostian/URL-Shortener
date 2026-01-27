# Go URL Shortener - Backend

A high-performance URL shortening service built with Go.

## Overview

This backend service provides a RESTful API for creating, managing, and redirecting shortened URLs. It's built using Go with PostgreSQL for data persistence.

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/
│   │   ├── router.go        # HTTP router configuration
│   │   └── handlers/        # HTTP request handlers
│   ├── config/
│   │   └── config.go        # Application configuration
│   ├── core/
│   │   └── shortener.go     # Core shortening logic
│   ├── db/
│   │   ├── models.go        # Data models
│   │   ├── postgres.go      # PostgreSQL connection
│   │   └── repository.go    # Data access layer
│   └── middleware/
│       └── cors.go          # CORS middleware
├── tests/
│   └── db_test.go           # Database tests
├── docs/
│   ├── architecture.md      # System architecture documentation
│   ├── decisions.md         # Design decisions log
│   └── screenshots/         # Project screenshots
├── init.sql/                # Database initialization scripts
├── docker-compose.yml       # Docker Compose configuration
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── .gitignore               # Git ignore rules
```

## Prerequisites

- Go 1.25.3 or higher
- PostgreSQL 14+
- Docker (optional, for containerized deployment)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone <repository-url>
cd url-shortener/backend
```

### 2. Configure Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials and other settings.

### 3. Start the Database

Using Docker Compose:

```bash
docker-compose up -d
```

Or connect to an existing PostgreSQL instance.

### 4. Run the Application

```bash
go run cmd/server/main.go
```

The server will start on the configured port (default: 8080).

## API Endpoints

| Method | Endpoint                     | Description              |
| ------ | ---------------------------- | ------------------------ |
| POST   | `/api/shorten`               | Create a shortened URL   |
| GET    | `/{shortCode}`               | Redirect to original URL |
| GET    | `/api/urls`                  | List all URLs            |
| GET    | `/api/urls/{shortCode}`      | Get URL details          |
| DELETE | `/api/urls/{shortCode}`      | Delete a URL             |
| GET    | `/api/analytics/{shortCode}` | Get URL analytics        |

## Running Tests

```bash
go test ./tests/...
```

## Building for Production

```bash
go build -o bin/server cmd/server/main.go
```

## Docker Deployment

```bash
docker-compose up --build
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver

## License

MIT License
