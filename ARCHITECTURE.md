# Campus Marketplace — Architecture Summary

A Go **modular monolith** — one deployable binary, modules separated by clear boundaries.

---

## Stack

| | |
|---|---|
| Language | Go 1.26.2+ |
| Router | `gin` |
| Database | PostgreSQL (schema-per-module) |
| Cache | Redis |
| Auth | JWT + refresh tokens |
| Storage | S3 / MinIO |
| Async | In-process event bus |

---

## Project Structure

```
campus-marketplace/
├── internal/
├── ├── server/server.go     # Entry point, wires everything together
│   ├── user/                # Auth, profiles
│   ├── listing/             # Items for sale
│   ├── order/               # Buyer-seller transactions
│   ├── chat/                # Real-time messaging
│   ├── notification/        # Email & push alerts
│   ├── media/               # File uploads
│   ├── review/              # Ratings
│   └── report/              # Content moderation
├── pkg/                     # Shared utilities (auth, bus, errors, pagination)
├── migrations/              # SQL migration files
└── config/config.go         # Env-based config
```

Each module contains: `handler.go` · `service.go` · `repository.go` · `model.go`

---

## The One Core Rule

**Modules never touch each other's database.** Cross-module calls go through a public `Service` interface only.

```go
// ❌ order module reaching into listing's DB — forbidden
listing.NewRepository(db).GetByID(...)

// ✅ order module using listing's public interface
type ListingService interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Listing, error)
    MarkAsSold(ctx context.Context, id uuid.UUID) error

```

---

## How Modules Talk to Each Other

**1. Direct call** — for reads that need an immediate response.
```
order.Service → calls → listing.Service.GetByID()
```

**2. Event bus** — for side-effects that don't need a response.
```
order.Service  publishes → "order.Created"
notification   subscribes → sends push to seller
```

Key events: `user.Registered`, `listing.Created`, `order.Created`, `order.Accepted`, `order.Completed`, `chat.MessageSent`

---

## Order Lifecycle

```
PENDING → ACCEPTED → COMPLETED
   │           │
REJECTED    CANCELLED
```

---

## API Routes (summary)

```
POST   /auth/register|login|refresh
GET    /listings          search & filter
POST   /listings          create (auth)
POST   /orders            place an offer (auth)
POST   /orders/:id/accept|reject|complete
GET    /chat/conversations/:id/ws   WebSocket
POST   /reviews
POST   /reports
```

---

## Key Conventions

- **Errors** — sentinel errors per module (`ErrNotFound`, `ErrForbidden`), mapped to HTTP codes in handlers.
- **Auth** — JWT claims carry `userID` + `verified`. Ownership checks happen in the service layer, not handlers.
- **No ORM** — raw `sqlx`. SQL lives only in `repository.go`.
- **Config** — environment variables only. `.env` for local, never committed.
- **Tests** — repositories tested against a real test DB; services tested with mocked repositories.
