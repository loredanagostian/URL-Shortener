# System Architecture

## Overview

The URL Shortener backend is built following Go best practices with a clean architecture pattern that separates concerns into distinct layers.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client (Frontend)                        │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HTTP Layer (API)                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Router    │  │  Handlers   │  │      Middleware         │  │
│  │ (gorilla/   │  │             │  │  (CORS, Logging, etc.)  │  │
│  │   mux)      │  │             │  │                         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Core Layer                                │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Shortener Service                     │    │
│  │         (URL generation, validation, business logic)     │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Data Access Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Models    │  │ Repository  │  │    PostgreSQL Driver    │  │
│  │             │  │  (CRUD)     │  │       (lib/pq)          │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PostgreSQL Database                        │
└─────────────────────────────────────────────────────────────────┘
```

## Layer Descriptions

### 1. HTTP Layer (`internal/api/`)

The HTTP layer handles all incoming HTTP requests and responses.

- **Router** (`router.go`): Configures routes using gorilla/mux
- **Handlers** (`handlers/`): Request handlers for each endpoint
    - `shorten.go`: URL shortening endpoint
    - `redirect.go`: Redirect to original URL
    - `url.go`: URL management (CRUD)
    - `analytics.go`: Usage analytics
- **Middleware** (`middleware/`): Cross-cutting concerns
    - `cors.go`: CORS configuration for frontend communication

### 2. Core Layer (`internal/core/`)

Contains the business logic for URL shortening.

- **Shortener** (`shortener.go`):
    - Generates unique short codes
    - Validates URLs
    - Manages URL lifecycle

### 3. Configuration (`internal/config/`)

Handles application configuration from environment variables.

- **Config** (`config.go`): Loads and validates configuration

### 4. Data Access Layer (`internal/db/`)

Manages all database interactions.

- **Models** (`models.go`): Data structures (URL entity)
- **Repository** (`repository.go`): CRUD operations interface
- **PostgreSQL** (`postgres.go`): Database connection management

## Data Flow

### URL Shortening Flow

```
1. Client POST /api/shorten with original URL
2. Router → Handler (shorten.go)
3. Handler validates request
4. Core Shortener generates unique code
5. Repository saves to PostgreSQL
6. Response with shortened URL
```

### Redirect Flow

```
1. Client GET /{shortCode}
2. Router → Handler (redirect.go)
3. Repository fetches original URL
4. Analytics updated (click count)
5. HTTP 301/302 redirect to original URL
```

## Database Schema

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    click_count INTEGER DEFAULT 0
);

CREATE INDEX idx_short_code ON urls(short_code);
```

## Security Considerations

1. **Input Validation**: All URLs are validated before processing
2. **SQL Injection Prevention**: Using parameterized queries
3. **CORS**: Configured to allow only trusted origins
4. **Rate Limiting**: (Future) Prevent abuse

## Scalability

The architecture supports horizontal scaling:

1. **Stateless Application**: No server-side sessions
2. **Database Connection Pooling**: Efficient connection management
3. **Containerization**: Docker support for easy deployment
