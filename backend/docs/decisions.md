# Design Decisions Log

This document records key architectural and design decisions made during the development of the URL Shortener backend.

---

## Decision 001: Go as Backend Language

**Date:** 2025

**Status:** Accepted

**Context:**  
We needed to choose a backend language that offers high performance, simplicity, and excellent concurrency support for handling many simultaneous URL redirects.

**Decision:**  
Use Go (Golang) as the primary backend language.

**Consequences:**

- ✅ Excellent performance for HTTP handling
- ✅ Built-in concurrency with goroutines
- ✅ Simple deployment (single binary)
- ✅ Strong standard library
- ⚠️ Team needs Go expertise

---

## Decision 002: PostgreSQL for Data Storage

**Date:** 2025

**Status:** Accepted

**Context:**  
Need a reliable database for storing URL mappings with support for indexing, transactions, and data integrity.

**Decision:**  
Use PostgreSQL as the primary database.

**Consequences:**

- ✅ ACID compliance ensures data integrity
- ✅ Excellent indexing for fast short code lookups
- ✅ Rich feature set (JSON, full-text search if needed)
- ✅ Strong community and tooling
- ⚠️ Requires database server management

**Alternatives Considered:**

- Redis (considered for caching layer)
- SQLite (rejected - not suitable for concurrent access)
- MongoDB (rejected - no need for document flexibility)

---

## Decision 003: Gorilla Mux for HTTP Routing

**Date:** 2025

**Status:** Accepted

**Context:**  
Need a flexible HTTP router that supports URL parameters, middleware, and route matching patterns.

**Decision:**  
Use gorilla/mux as the HTTP router.

**Consequences:**

- ✅ Mature and well-tested library
- ✅ URL parameter extraction (for short codes)
- ✅ Middleware support
- ✅ Good documentation
- ⚠️ External dependency (minimal risk)

**Alternatives Considered:**

- Chi router (similar features)
- Standard library net/http (less flexible routing)
- Gin (more opinionated)

---

## Decision 004: Repository Pattern for Data Access

**Date:** 2025

**Status:** Accepted

**Context:**  
Need a clean separation between business logic and database operations for testability and maintainability.

**Decision:**  
Implement the Repository pattern for all database operations.

**Consequences:**

- ✅ Business logic is decoupled from database
- ✅ Easy to mock for unit testing
- ✅ Can swap database implementations
- ✅ Clear interfaces for data operations
- ⚠️ Additional abstraction layer

---

## Decision 005: Internal Package Structure

**Date:** 2025

**Status:** Accepted

**Context:**  
Go projects need clear package organization to prevent import cycles and maintain encapsulation.

**Decision:**  
Use the `internal/` directory for application-specific packages that should not be imported by external projects.

**Package Structure:**

- `internal/api/` - HTTP layer
- `internal/config/` - Configuration
- `internal/core/` - Business logic
- `internal/db/` - Data access
- `internal/middleware/` - HTTP middleware

**Consequences:**

- ✅ Clear separation of concerns
- ✅ Enforced encapsulation by Go compiler
- ✅ Easy to navigate codebase

---

## Decision 006: Short Code Generation Strategy

**Date:** 2025

**Status:** Accepted

**Context:**  
Need to generate unique, short, URL-safe codes for shortened URLs.

**Decision:**  
Use a random alphanumeric string generation approach with collision detection.

**Implementation:**

- Character set: a-z, A-Z, 0-9
- Default length: 6 characters
- Check database for existing codes before saving

**Consequences:**

- ✅ Simple to implement
- ✅ URL-safe characters
- ✅ ~56 billion combinations (62^6)
- ⚠️ Collision possible (handled with retry)

**Alternatives Considered:**

- Base62 encoding of auto-increment ID (predictable)
- UUID truncation (less readable)
- Hash-based (collision handling complex)

---

## Decision 007: URL Expiration Support

**Date:** 2025

**Status:** Accepted

**Context:**  
Some use cases require temporary shortened URLs that expire after a certain time.

**Decision:**  
Add optional expiration timestamp to URL model with default expiration of 60 minutes.

**Consequences:**

- ✅ Supports temporary links
- ✅ Automatic cleanup possible
- ✅ Optional - can create permanent links
- ⚠️ Need background job for cleanup (future)

---

## Decision 008: CORS Configuration

**Date:** 2025

**Status:** Accepted

**Context:**  
Frontend application needs to communicate with the backend API from a different origin during development.

**Decision:**  
Implement CORS middleware with configurable origins.

**Consequences:**

- ✅ Frontend can access API
- ✅ Configurable for different environments
- ⚠️ Must be properly restricted in production

---

## Future Decisions to Consider

1. **Rate Limiting Strategy** - Prevent API abuse
2. **Caching Layer** - Redis for frequently accessed URLs
3. **Authentication** - User accounts for URL management
4. **Analytics Enhancement** - Detailed click tracking
5. **Custom Short Codes** - User-defined codes
