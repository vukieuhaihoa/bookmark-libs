# bookmark-libs

A shared Go library providing common infrastructure and cross-cutting concerns for the bookmark management system, built following Domain-Driven Design (DDD) principles.

## Overview

`bookmark-libs` is designed to be reused across multiple services in a DDD architecture. It centralizes common concerns — authentication, rate limiting, database connectivity, structured responses, and utilities — so individual services stay focused on business logic.

## Packages

### `pkg/common`
Shared types for HTTP response formatting and query handling.
- `SuccessResponse[T]` — generic structured success response
- `Paging` — pagination parameters with bounds normalization (page ≥ 1, 1 ≤ limit ≤ 50)
- `ParseSortParams` — parses comma-separated sort strings (e.g. `-created_at,updated_at`) against an allowlist
- `QueryOptions` — combines `Paging` and `[]SortedField` into a single query descriptor

### `pkg/dbutils`
Database error categorization for consistent cross-service error handling. Wraps raw database errors into typed errors: duplication, not found, foreign key violation, invalid sort field.

### `pkg/encoding`
Base62 encoding/decoding (alphabet: `0–9`, `A–Z`, `a–z`) for compact integer representation — useful for short IDs in URLs.

### `pkg/jwtutils`
RSA-based JWT token generation and validation.
- `JWTGenerator` — signs tokens using an RSA private key (RS256)
- `JWTValidator` — validates tokens using an RSA public key

### `pkg/logger`
Global log level configuration via `zerolog`. Reads `LOG_LEVEL` from environment (defaults to `info`).

### `pkg/redis`
Redis client factory. Reads connection config from environment variables (`REDIS_ADDRESS`, `REDIS_PASSWORD`, `REDIS_DB`).

### `pkg/sqldb`
PostgreSQL connectivity and schema migrations using GORM and `golang-migrate`.
- Reads DB config from environment (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`)
- Supports `up` and `steps` migration modes

### `pkg/utils`
General-purpose utilities:
- `PasswordHashing` — bcrypt-based password hashing and comparison
- `CodeGenerator` — cryptographically secure random alphanumeric codes
- JWT claims extraction from a Gin context (extracts user ID from `sub` claim)

### `pkg/csv`
CSV parsing utilities for multipart file uploads.
- `ParseFromMultipartFile` — decodes a multipart uploaded CSV file into a target struct slice using struct field tags

### `pkg/validators`
Custom `go-playground/validator` rules:
- `PasswordStrength` — requires at least one uppercase, lowercase, digit, and special character (`@#$%!^&*()_+`)

### `middlewares`
Gin middleware:
- `JWTAuth` — extracts and validates Bearer tokens from the `Authorization` header, injects claims into context
- `RateLimit` — per-IP (100 req/min) and per-user-ID (20 req/s) rate limiting backed by Redis

### `ratelimit`
Redis-backed rate limit repository:
- `GetCurrentRateLimit` — reads current counter for a key
- `IncreaseRateLimit` — atomically increments counter and sets expiration

## Requirements

- Go 1.25+
- PostgreSQL
- Redis

## Installation

```bash
go get github.com/vukieuhaihoa/bookmark-libs
```

## Usage

Import the packages you need:

```go
import (
    "github.com/vukieuhaihoa/bookmark-libs/pkg/common"
    "github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils"
    "github.com/vukieuhaihoa/bookmark-libs/middlewares"
)
```

## Development

### Generate mocks

```bash
make mock-gen
```

### Run tests

```bash
make test
```

Runs the full test suite, generates an HTML coverage report at `coverage/coverage.html`, and enforces an **80% coverage threshold**.

### Run tests in Docker

```bash
make docker-test
```

### Clean test cache and coverage artifacts

```bash
make clean
```

## Key Dependencies

| Dependency | Purpose |
|---|---|
| `gin-gonic/gin` | HTTP framework |
| `golang-jwt/jwt/v5` | JWT handling |
| `go-playground/validator/v10` | Struct validation |
| `golang-migrate/migrate/v4` | Schema migrations |
| `gorm.io/gorm` | ORM |
| `redis/go-redis/v9` | Redis client |
| `rs/zerolog` | Structured logging |
| `golang.org/x/crypto` | bcrypt |
| `stretchr/testify` | Test assertions |
| `alicebob/miniredis/v2` | In-memory Redis for tests |
| `gocarina/gocsv` | CSV encoding/decoding |
