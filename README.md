# Banking System Simulator

A simplified backend banking system simulator built with Go, demonstrating clean
architecture, REST API design, and safe concurrent money transfers.

> **Note:** This is a simplified version of the original task. RabbitMQ, Redis,
> and gRPC (originally specified for a microservices setup) were intentionally
> left out to focus on core backend logic within the given timeframe. All
> operations that were meant to be asynchronous (via RabbitMQ) or use
> service-to-service gRPC calls are instead handled synchronously within a
> single monolithic service.

## Tech Stack

- **Language:** Go 1.26
- **HTTP Framework:** Gin
- **Database:** PostgreSQL (via `pgx` / `pgxpool`)
- **Migrations:** golang-migrate
- **Architecture:** Clean Architecture (Handler → Usecase → Repository)
- **Containerization:** Docker & Docker Compose

## Architecture

```
Handler  →  Usecase  →  Repository  →  PostgreSQL
```

- **Handler** — parses/validates HTTP requests, formats responses
- **Usecase** — business rules (locking checks, currency validation, etc.)
- **Repository** — database access, transactions, row-level locking

## Project Structure

```
.
├── cmd/                    # application entrypoint
├── internal/
│   ├── app/                # wiring / Run() function
│   ├── config/              # config loading
│   ├── db/                  # DB connection + migrations runner
│   ├── errors/               # centralized app errors
│   ├── handlers/             # HTTP handlers + request/response DTOs
│   ├── models/                # domain models
│   ├── repository/             # data access layer
│   ├── router/                  # route registration
│   ├── usecase/                  # business logic
│   └── validations/                # input validation rules
├── migrations/               # SQL migration files
├── config.yaml                # local config (gitignored recommended)
├── config.docker.yaml          # config used inside Docker containers
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Running the Project

### Option A — with Docker (recommended)

```bash
docker-compose up -d --build
```

This will:
1. Start a PostgreSQL container
2. Build and start the app container
3. Automatically run database migrations on startup

The API will be available at `http://localhost:8080`.

Check logs:
```bash
docker compose logs app
```

Stop everything:
```bash
docker-compose down
```

### Option B — running locally (without Docker)

1. Make sure PostgreSQL is running locally and create the database:
   ```bash
   psql -U postgres -p 5433 -c "CREATE DATABASE banking;"
   ```
2. Update `config.yaml` with your local PostgreSQL credentials.
3. Run the app:
   ```bash
   go run cmd/main.go
   ```

Migrations run automatically on startup in both cases.

## API Endpoints

### Accounts

| Method | Endpoint              | Description                              |
|--------|------------------------|-------------------------------------------|
| POST   | `/api/accounts`         | Create a new account                     |
| GET    | `/api/accounts/:id`     | Get account by ID                        |
| GET    | `/api/accounts`         | List accounts (pagination, currency filter) |
| DELETE | `/api/accounts/:id`     | Soft-delete an account (must have zero balance and not be locked) |

### Transactions

| Method | Endpoint                              | Description                              |
|--------|-----------------------------------------|-------------------------------------------|
| POST   | `/api/accounts/:id/deposit`             | Deposit funds into an account            |
| POST   | `/api/accounts/:id/withdraw`            | Withdraw funds from an account           |
| POST   | `/api/accounts/:id/transfer`            | Transfer funds to another account (same currency only) |
| GET    | `/api/accounts/:id/transactions`         | List an account's transaction history (pagination, type/date-range filter) |

### Example requests

**Create account**
```json
POST /api/accounts
{
  "balance": 0,
  "currency": "USD"
}
```

**Deposit**
```json
POST /api/accounts/1/deposit
{
  "amount": 500
}
```

**Transfer**
```json
POST /api/accounts/1/transfer
{
  "amount": 100,
  "to_account_id": 2
}
```

**List transactions with filters**
```
GET /api/accounts/1/transactions?type=deposit&start_date=2026-01-01&end_date=2026-07-01&page=1&limit=20
```

## Error Handling

Errors are returned as JSON with a consistent shape:
```json
{
  "message": "not enough balance"
}
```

The HTTP status code is derived from an internal error code
(e.g. bad request → 400, not found → 404, internal errors → 500).

## Definition of Done Checklist

- [x] `docker-compose up` runs, all containers healthy
- [x] Database migrations applied automatically
- [x] Account creation via Postman/cURL
- [x] Deposit verified — balance increases in DB
- [x] Withdraw verified — balance decreases, rejected if insufficient funds
- [x] Transfer verified — sender debited, receiver credited, currency mismatch rejected
- [x] Transaction history retrievable with pagination and filters
