# CineBook

CineBook is a lightweight movie seat-booking REST API written in Go that lets users browse movies, temporarily hold seats with TTL-based sessions, and confirm or release bookings, backed by Redis for concurrency-safe holds and SQLite for movie catalog storage.  
The ultimate goal for this project is for me to learn Golang and Backend, that's why you'll see the frontend is done on the bare minimum level.

## Features

- Movie catalog CRUD management of movies (title, poster, seat grid dimensions) stored in SQLite, with seed data on startup.
- Seat hold with TTL: seats are held for 1 minute using Redis SET NX semantics, guaranteeing a seat can only be held by one user at a time.
- Confirm and release sessions: confirming a hold persists the booking (removes TTL); releasing deletes the session.
- Concurrent-safe stores: pluggable BookingStore interface with Redis, in-memory, and concurrent map implementations.
- Admin panel: basic-auth-protected admin page for movie management.
- Dev tooling: Air live-reload config and Docker Compose for Redis and RedisInsight.

## Architecture

```
┌─────────────┐     ┌──────────────────────────────────────────┐
│ static/     │     │            cmd/main.go                   │
│ (HTML/JS/CSS│────▶│  HTTP mux · routing · middleware         │
│  frontend)  │     └───────────────┬──────────────────────────┘
└─────────────┘                     │
                                    ▼
                    ┌───────────────────────────────┐
                    │   internal/booking            │
                    │  Handler → Service → Store    │
                    │  (BookingStore interface)     │
                    └──────┬───────────────┬────────┘
                           ▼               ▼
                    ┌────────────┐  ┌──────────────┐
                    │   SQLite   │  │    Redis     │
                    │ (movies)   │  │ (holds/TTL)  │
                    └────────────┘  └──────────────┘
```

The project uses a layered design. HTTP handlers (handler.go, DTOs) call the service layer (service.go, business logic), which delegates to a persistence layer decoupled via the BookingStore interface.

The Redis data model uses a `session:<id>` to seat key mapping plus a `seat:<movieID>:<seatID>` to booking JSON structure with SET NX and TTL for atomic, expiring holds.

Middleware provides basic auth for /admin and trailing-slash stripping for consistent routing.  
Better auth method will be implemented in the future

## Prerequisites
- Go 1.25+
- Redis (locally installed, or via Docker)
- Docker and Docker Compose (optional, for Redis and RedisInsight)
- Air (optional, for live reload)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/kkkfasya/CineBook.git
   cd CineBook
   ```

2. Copy the environment example and adjust credentials:

   ```bash
   cp .env.example .env
   ```

   | Variable          | Description            |
   |-------------------|------------------------|
   | ADMIN_USERNAME    | Admin panel username   |
   | ADMIN_PASSWORD    | Admin panel password   |

3. Start Redis (and RedisInsight, optional):

   ```bash
   docker compose up -d
   ```

4. Run the server:

   ```bash
   go run cmd/main.go
   ```

   Or with live reload:

   ```bash
   air
   ```

The server listens on http://localhost:8080.

## Usage

### Web UI

- Booking page: http://localhost:8080/
- Admin panel (basic auth): http://localhost:8080/admin

### API Examples

List movies:

```bash
curl http://localhost:8080/api/v1/movies
```

List seats for a movie:

```bash
curl http://localhost:8080/api/v1/movies/{movieID}/seats
```

Hold a seat (1-minute TTL):

```bash
curl -X POST http://localhost:8080/api/v1/movies/{movieID}/seats/{seatID}/hold \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123"}'
```

Confirm a held session:

```bash
curl -X PUT http://localhost:8080/api/v1/sessions/{sessionID}/confirm \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123"}'
```

Release a session:

```bash
curl -X DELETE http://localhost:8080/api/v1/sessions/{sessionID} \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123"}'
```

Add a movie (admin CRUD also available at /api/v1/admin/movies):

```bash
curl -X POST http://localhost:8080/api/v1/admin/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Inception", "poster": "https://example.com/poster.jpg", "rows": 5, "seats_per_row": 8}'
```

### Endpoint Summary

| Method   | Path                                              | Description              |
|----------|---------------------------------------------------|--------------------------|
| GET      | /health                                           | Health check             |
| GET      | /api/v1/movies                                    | List movies              |
| GET      | /api/v1/movies/{movieID}/seats                    | List booked seats        |
| POST     | /api/v1/movies/{movieID}/seats/{seatID}/hold      | Hold a seat              |
| PUT      | /api/v1/sessions/{sessionID}/confirm              | Confirm a hold           |
| DELETE   | /api/v1/sessions/{sessionID}                      | Release a hold           |
| GET/POST/PUT/DELETE | /api/v1/admin/movies[/{movieID}]       | Admin movie CRUD         |

## Project Structure

```
CineBook/
├── cmd/
│   ├── main.go          # Server bootstrap, routing, static file serving
│   └── db.go            # SQLite schema creation and seed data
├── internal/
│   ├── booking/
│   │   ├── domain.go          # Booking model + BookingStore interface
│   │   ├── service.go         # Business logic layer
│   │   ├── handler.go         # HTTP handlers and DTOs
│   │   ├── redis_store.go     # Redis-backed store (TTL holds, NX locking)
│   │   ├── memory_store.go    # In-memory map store
│   │   ├── concurrent_store.go# Mutex-protected in-memory store
│   │   └── service_test.go    # Tests
│   ├── adapters/redis/
│   │   └── client.go          # Redis client adapter
│   ├── middleware/
│   │   ├── basic_auth.go      # Basic auth middleware
│   │   └── strip_path.go      # Trailing-slash stripping middleware
│   └── utils/
│       └── utils.go           # JSON/error helpers, env loading, Redis ping
├── static/              # Vanilla HTML/JS/CSS frontend (booking + admin)
├── docker-compose.yaml  # Redis + RedisInsight services
├── .air.toml            # Air live-reload config
├── Makefile             # Test targets (test, test_no_cache, test_race)
└── .env.example         # Environment variable template
```

## Testing

```bash
make test        # run all tests verbosely
make test_race   # run tests with race detector
```

## Contributing

1. Fork the repository and create a feature branch:

   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes and ensure tests pass:

   ```bash
   make test_race
   ```

3. Commit with a clear message and open a pull request against main.

Please keep changes focused, follow existing code style (idiomatic Go, layered architecture), and update tests for new behavior.  
AI usage is encouraged but please read the code your AI produces, do not throw slop grenade at me.

## License
MIT

