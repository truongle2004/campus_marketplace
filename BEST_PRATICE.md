# Go Best Practices — Campus Marketplace

Guidelines every developer (and AI agent) must follow when writing code in this project.

---

## 1. Dependency Injection (DI)

**Never instantiate dependencies inside a struct. Receive them from outside.**

```go
// ❌ Bad — hard to test, hidden dependency
func NewOrderService() *OrderService {
    return &OrderService{
        repo: postgres.NewOrderRepository(), // hardcoded
        bus:  bus.Global,                    // global state
    }
}

// ✅ Good — dependencies are explicit and swappable
func NewOrderService(repo Repository, bus bus.Bus, listing ListingService) *OrderService {
    return &OrderService{repo: repo, bus: bus, listing: listing}
}
```

> Everything gets wired once in `cmd/server/main.go`. Nowhere else.

---

## 2. Program to Interfaces

**Depend on interfaces, not concrete types.** The interface is defined by the consumer — inside the package that uses it, not the package that implements it.

```go
// ❌ Bad — depends on a concrete struct
type OrderService struct {
    repo *postgres.OrderRepository
}

// ✅ Good — depends on a behaviour
type Repository interface {
    Create(ctx context.Context, o *Order) error
    GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
}

type OrderService struct {
    repo Repository // any implementation works
}
```

Keep interfaces **small**. One method is fine. The bigger an interface, the harder it is to mock.

---

## 3. DRY — Don't Repeat Yourself

Extract logic the moment it appears a second time.

```go
// ❌ Repeated in every handler
userID, ok := r.Context().Value(ctxKeyUserID).(uuid.UUID)
if !ok {
    http.Error(w, "unauthorized", 401)
    return
}

// ✅ Extracted once in pkg/auth
userID := auth.UserIDFromContext(r.Context()) // panics if middleware not applied
```

**Where to put shared code:**

| Scope | Location |
|-------|----------|
| Used by one module only | inside that module |
| Used by 2+ modules | `pkg/` |
| HTTP helpers | `pkg/respond` or `pkg/apierr` |
| Test helpers | `pkg/testutil` |

---

## 4. Single Responsibility

Each type does one thing. Split when a file grows beyond ~200 lines or handles more than one concern.

```
listing/
  handler.go      ← HTTP only: decode, validate, call service, respond
  service.go      ← business rules only
  repository.go   ← SQL only
  search.go       ← search logic extracted from service when it grows
```

A handler should never contain SQL. A repository should never contain business rules.

---

## 5. Error Handling

**Always wrap errors with context. Never swallow them.**

```go
// ❌ Bad
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    return nil, err  // no context
}

// ✅ Good
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    return nil, fmt.Errorf("UserService.GetByID %s: %w", id, err)
}
```

**Use sentinel errors for expected failures:**

```go
// internal/listing/errors.go
var (
    ErrNotFound    = errors.New("listing not found")
    ErrForbidden   = errors.New("not the listing owner")
    ErrAlreadySold = errors.New("listing already sold")
)
```

**Map to HTTP codes in handlers only — never in service or repository:**

```go
switch {
case errors.Is(err, listing.ErrNotFound):
    respond.NotFound(w)
case errors.Is(err, listing.ErrForbidden):
    respond.Forbidden(w)
default:
    respond.InternalError(w, err)
}
```

---

## 6. Keep Functions Small

A function should do **one thing** and fit on one screen (~30 lines max). If you need a comment to explain a block inside a function, that block should be its own function.

```go
// ❌ One giant function
func (s *service) PlaceOrder(ctx context.Context, req Request) (*Order, error) {
    // validate buyer
    // check listing exists
    // check listing is not sold
    // check buyer != seller
    // create order record
    // publish event
    // send notification
}

// ✅ Composed from small, named functions
func (s *service) PlaceOrder(ctx context.Context, req Request) (*Order, error) {
    if err := s.validateBuyer(ctx, req.BuyerID, req.ListingID); err != nil {
        return nil, err
    }
    order, err := s.repo.Create(ctx, newOrder(req))
    if err != nil {
        return nil, fmt.Errorf("PlaceOrder: %w", err)
    }
    s.bus.Publish(ctx, orderCreatedEvent(order))
    return order, nil
}
```

---

## 7. Context Is Not a Bag

Only pass values in `context.Context` that are **request-scoped and cross-cutting** — things like `userID`, `requestID`, `traceID`. Never use context to pass business logic dependencies.

```go
// ❌ Business data in context
ctx = context.WithValue(ctx, "listing", listing)

// ✅ Pass explicitly as arguments
func (s *service) CreateOrder(ctx context.Context, listing *Listing, buyerID uuid.UUID) ...
```

---

## 8. Testing

### Structure
```
internal/listing/
  service_test.go      ← unit tests with mocked repo
  repository_test.go   ← integration tests against real test DB
  handler_test.go      ← HTTP tests with httptest + mocked service
```

### Unit test with mock
```go
func TestService_Create_ForbidsUnverifiedBuyer(t *testing.T) {
    repo := &mockRepository{}
    userSvc := &mockUserService{verified: false}
    svc := NewService(repo, userSvc, bus.NewNoop())

    _, err := svc.Create(context.Background(), createReq)

    assert.ErrorIs(t, err, ErrBuyerNotVerified)
    assert.Empty(t, repo.created) // repo was never called
}
```

### Rules
- **Never test implementation** — test behaviour and outcomes.
- **One assertion per test** where possible. One `t.Run` per scenario.
- **No real DB in unit tests.** Use `pkg/testutil.NewTestDB(t)` only in `repository_test.go`.
- Table-driven tests for multiple input cases:

```go
tests := []struct {
    name    string
    input   Request
    wantErr error
}{
    {"missing title", Request{Title: ""}, ErrValidation},
    {"price negative", Request{Title: "x", Price: -1}, ErrValidation},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        _, err := svc.Create(ctx, tt.input)
        assert.ErrorIs(t, err, tt.wantErr)
    })
}
```

---

## 9. Avoid Global State

No package-level variables that hold mutable state or connections.

```go
// ❌ Bad
var DB *sql.DB  // package-level global

// ✅ Good — injected
type Repository struct {
    db *sql.DB
}
```

The only acceptable package-level values are **constants** and **sentinel error variables**.

---

## 10. Naming

| Thing | Convention | Example |
|-------|-----------|---------|
| Interface | noun (what it is) | `Repository`, `Notifier` |
| Constructor | `New` + type name | `NewOrderService` |
| Boolean | starts with `Is`/`Has`/`Can` | `IsVerified`, `HasExpired` |
| Error variable | `Err` prefix | `ErrNotFound` |
| Context key | unexported type | `type ctxKey string` |
| Test helper | `must` prefix | `mustCreateUser(t, ...)` |

Avoid generic names like `data`, `info`, `manager`, `handler2`. Name things by what they represent.

---

## 11. Validation

Validate at the **edge** — in the handler, before calling the service. The service trusts its inputs.

```go
// handler.go — validate here
type CreateListingRequest struct {
    Title    string  `json:"title"    validate:"required,min=3,max=100"`
    Price    float64 `json:"price"    validate:"required,gt=0"`
    Category string  `json:"category" validate:"required"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateListingRequest
    json.NewDecoder(r.Body).Decode(&req)
    if err := validate.Struct(req); err != nil {
        respond.ValidationError(w, err)
        return
    }
    // service receives clean, validated data
    h.svc.Create(r.Context(), req)
}
```

---

## Summary Cheatsheet

| Principle | Rule |
|-----------|------|
| DI | Inject dependencies via constructor, wire in `main.go` |
| Interfaces | Define at consumer side, keep small |
| DRY | Extract on second use; shared code goes in `pkg/` |
| SRP | One file = one concern; split when >200 lines |
| Errors | Wrap with context; sentinel errors; HTTP mapping in handler only |
| Functions | One thing, one screen; extract named helpers |
| Context | Request-scoped values only (`userID`, `traceID`) |
| Testing | Unit = mocks; Integration = real DB; test behaviour not internals |
| Globals | Constants and sentinel errors only |
| Naming | Be specific; no `manager`, `data`, `info` |
| Validation | At the edge (handler); service trusts its inputs |
