# Campus Marketplace — Use Cases

## Tech Stack

| Layer | Technology |
|---|---|
| Mobile client | Flutter (Android) |
| Backend API | Go (Golang) |
| Auth | Clerk (hosted; free tier for development / modest traffic) |
| Cloud emulation | LocalStack (Community edition, local dev) |
| Observability | Grafana + Prometheus (self-hosted) |
| Storage | LocalStack S3 (emulated) |
| Database | PostgreSQL (local) |
| Message queue | LocalStack SQS (emulated) |

---

## Constraints

- **Runtime split:** Flutter runs on device/emulator; **Go, PostgreSQL, LocalStack, Grafana, Prometheus** run on developer machine / server reachable from Android (LAN or emulator bridge). Clerk is **cloud-hosted**.
- **Auth external dependency:** **Clerk** provides sign-up, sessions, JWTs, and built-in transactional email for verification/password reset—no separate self-hosted identity server or SMTP product required for MVP.
- **No payment gateway** → cash-only transactions, agreed in person.
- **No FCM push** → in-app notifications only via polling or WebSocket.
- **Image storage** → LocalStack S3 on dev host in local dev; prod uses real S3-compatible storage when needed (not Clerk).
- **Background jobs** → Go goroutines / cron scheduler (no cloud functions needed).

---

## Progress Tracker

- Total use cases: 58
- Completed: 0 / 58

---

## 1. Authentication & User Management

> Handled by **Clerk** (sign-up, sessions, user profile fields). Go validates **JWTs** (JWKS from Clerk). Flutter uses **Clerk Flutter SDK** (OIDC-style session tokens for API calls).

- [ ] Register with campus email — `Clerk sign-up + optional allowed-email domain rules`
- [ ] Auto-verify email domain (regex check on campus domain) — `Go middleware` *(and/or Clerk-restricted sign-up domains)*
- [ ] Verify email with confirmation link — `Clerk email verification`
- [ ] Login with email and password — `Clerk session (Flutter SDK)`
- [ ] Logout / invalidate token — `Clerk signOut / session revocation`
- [ ] Refresh access token — `Clerk token refresh handled by SDK / session lifecycle`
- [ ] Forgot password — `Clerk password reset`
- [ ] Reset password via token — `Clerk password reset completion`
- [ ] Change password (authenticated) — `Clerk user profile / account settings`
- [ ] View own profile — `Go API → PostgreSQL`
- [ ] Edit profile (name, avatar) — `Go API → PostgreSQL + LocalStack S3`
- [ ] Deactivate own account — `Go API + Clerk user delete/disable via Clerk Backend API`

---

## 2. Campus & Categories

> Static-ish data served by **Go API**, cached in Flutter.

- [ ] List all active campuses — `Go API → PostgreSQL`
- [ ] List top-level categories — `Go API → PostgreSQL`
- [ ] List subcategories by parent — `Go API → PostgreSQL`

---

## 3. Listings

> Core CRUD on **Go API + PostgreSQL**. Images stored on **LocalStack S3**.

- [ ] Create a listing (title, description, price, condition, category, images, tags) — `Go API → PostgreSQL + S3`
- [ ] Edit a listing — `Go API → PostgreSQL`
- [ ] Delete a listing (soft delete) — `Go API → PostgreSQL`
- [ ] View listing detail (increment view count) — `Go API → PostgreSQL`
- [ ] Browse listings (filter by campus, category, condition, price range, tags) — `Go API → PostgreSQL`
- [ ] Search listings by keyword (full-text) — `Go API → PostgreSQL full-text (tsvector)`
- [ ] List own listings by status — `Go API → PostgreSQL`
- [ ] Mark listing as sold — `Go API → PostgreSQL`
- [ ] Archive a listing manually — `Go API → PostgreSQL`
- [ ] Auto-expire listings past `expires_at` — `Go cron job (robfig/cron)`
- [ ] Upload listing images — `Go API → LocalStack S3`
- [ ] Reorder / delete listing images — `Go API → PostgreSQL + S3`

---

## 4. Messaging & Offers

> Real-time via **Go WebSocket** (gorilla/websocket). Conversations and messages persisted in **PostgreSQL**.

- [ ] Start a conversation on a listing — `Go API → PostgreSQL`
- [ ] Send a message in a conversation — `Go WebSocket → PostgreSQL`
- [ ] View conversation thread — `Go API → PostgreSQL`
- [ ] List inbox (sorted by latest message) — `Go API → PostgreSQL`
- [ ] Mark messages as read — `Go API → PostgreSQL`
- [ ] Make an offer inside a conversation — `Go API → PostgreSQL`
- [ ] Accept an offer — `Go API → PostgreSQL`
- [ ] Reject an offer — `Go API → PostgreSQL`
- [ ] Withdraw an offer (buyer) — `Go API → PostgreSQL`
- [ ] Auto-expire offers past `expires_at` — `Go cron job`

---

## 5. Commerce

> Cash-only. No payment gateway. Transaction is a manual record created after both parties agree in person.

- [ ] Create a transaction when an offer is accepted — `Go API → PostgreSQL`
- [ ] Record payment method as `cash` — `Go API → PostgreSQL`
- [ ] Mark transaction as completed (manual confirmation) — `Go API → PostgreSQL`
- [ ] View transaction history (buyer & seller) — `Go API → PostgreSQL`

---

## 6. Reviews

> Post-transaction ratings stored in **PostgreSQL**, served by **Go API**.

- [ ] Leave a review after a completed transaction (buyer reviews seller) — `Go API → PostgreSQL`
- [ ] Leave a review after a completed transaction (seller reviews buyer) — `Go API → PostgreSQL`
- [ ] View reviews on a user profile — `Go API → PostgreSQL`
- [ ] View average rating on a user profile — `Go API → PostgreSQL (avg aggregate)`

---

## 7. Discovery

> Wishlists and saved searches in **PostgreSQL**. Saved search matching runs as a **Go background job** triggered on new listing creation.

- [ ] Save a listing to a wishlist — `Go API → PostgreSQL`
- [ ] Remove a listing from a wishlist — `Go API → PostgreSQL`
- [ ] Create a wishlist — `Go API → PostgreSQL`
- [ ] Rename a wishlist — `Go API → PostgreSQL`
- [ ] Delete a wishlist — `Go API → PostgreSQL`
- [ ] View wishlist items — `Go API → PostgreSQL`
- [ ] Save a search with filters — `Go API → PostgreSQL`
- [ ] Toggle new-listing notifications on a saved search — `Go API → PostgreSQL`
- [ ] Delete a saved search — `Go API → PostgreSQL`
- [ ] Match new listings against saved searches and notify — `Go background job → LocalStack SQS → notifications`

---

## 8. Notifications

> In-app only. No FCM. Notifications written to **PostgreSQL** and delivered via **Go WebSocket** or short polling from Flutter.

- [ ] Receive notification on new message — `Go WebSocket`
- [ ] Receive notification on offer received — `Go WebSocket`
- [ ] Receive notification on offer accepted / rejected — `Go WebSocket`
- [ ] Receive notification on listing sold — `Go WebSocket`
- [ ] Receive notification on new review — `Go WebSocket`
- [ ] Receive notification on new follower — `Go WebSocket`
- [ ] Receive notification on saved search match — `Go SQS consumer → WebSocket`
- [ ] Mark a notification as read — `Go API → PostgreSQL`
- [ ] Mark all notifications as read — `Go API → PostgreSQL`
- [ ] List notifications (unread first) — `Go API → PostgreSQL`

---

## 9. Trust & Safety

> Reports stored in **PostgreSQL**. Admin actions via a simple Go admin API (no separate UI needed initially).

- [ ] Report a listing — `Go API → PostgreSQL`
- [ ] Report a user — `Go API → PostgreSQL`
- [ ] Report a message — `Go API → PostgreSQL`
- [ ] Admin: list pending reports — `Go API → PostgreSQL`
- [ ] Admin: resolve a report — `Go API → PostgreSQL`
- [ ] Admin: dismiss a report — `Go API → PostgreSQL`
- [ ] Admin: deactivate a user — `Go API → Clerk Backend API + PostgreSQL`
- [ ] Admin: remove a listing — `Go API → PostgreSQL`

---

## 10. Social

> Simple follow graph in **PostgreSQL**.

- [ ] Follow a seller — `Go API → PostgreSQL`
- [ ] Unfollow a seller — `Go API → PostgreSQL`
- [ ] List followers — `Go API → PostgreSQL`
- [ ] List following — `Go API → PostgreSQL`
- [ ] View new listings from followed sellers — `Go API → PostgreSQL`

---

## 11. Observability

> **Grafana + Prometheus** running locally. Go backend exposes `/metrics` endpoint.

- [ ] Expose Prometheus metrics endpoint in Go API — `prometheus/client_golang`
- [ ] Track active users, listing count, message rate — `Prometheus counters / gauges`
- [ ] Track API latency per endpoint — `Prometheus histogram`
- [ ] Track error rate — `Prometheus counter`
- [ ] Set up Grafana dashboard for API health — `Grafana + Prometheus datasource`
- [ ] Set up Grafana dashboard for business metrics (listings, transactions, reviews) — `Grafana + PostgreSQL datasource`

---

## Implementation Order

| Phase | Domains | Stack focus |
|---|---|---|
| 1 | Auth, Campus, Categories | Clerk app + JWT validation (JWKS) in Go |
| 2 | Listings, Images | Go CRUD + LocalStack S3 |
| 3 | Messaging, Offers | Go WebSocket + PostgreSQL |
| 4 | Commerce, Reviews | Go API + PostgreSQL |
| 5 | Discovery, Notifications | Go cron + SQS + WebSocket |
| 6 | Social, Trust & Safety | Go API + Clerk Backend API |
| 7 | Observability | Prometheus + Grafana dashboards |

---

## Cut / Simplified Use Cases

| Use case | Original plan | Reason cut / simplified |
|---|---|---|
| Email delivery (verify, reset) | SendGrid / Mailgun separate from auth | Covered by Clerk’s built-in transactional email |
| Payment gateway integration | MoMo / VNPay / Stripe | Paid services — replaced with manual cash recording |
| Push notifications (FCM) | Firebase Cloud Messaging | Requires Google Play Services + account — replaced with in-app WebSocket |
| Image CDN | Cloudinary / AWS S3 | Listing images via LocalStack in dev; S3-compatible or similar in prod when needed |
