# Campus Marketplace — Database Schema

## Overview

This schema is designed for a long-term, scalable campus marketplace. All primary keys use `uuid` to prevent enumeration attacks and support future multi-region sharding. Soft deletes (`is_active` / `deleted_at`) are preferred everywhere to support dispute auditing.

---

## Table of Contents

1. [Core Identity](#1-core-identity)
   - [campuses](#campuses)
   - [users](#users)
2. [Catalog](#2-catalog)
   - [categories](#categories)
   - [listings](#listings)
   - [listing_images](#listing_images)
   - [tags](#tags)
   - [listing_tags](#listing_tags)
3. [Messaging & Offers](#3-messaging--offers)
   - [conversations](#conversations)
   - [messages](#messages)
   - [offers](#offers)
4. [Commerce](#4-commerce)
   - [transactions](#transactions)
   - [reviews](#reviews)
5. [Discovery](#5-discovery)
   - [wishlists](#wishlists)
   - [wishlist_items](#wishlist_items)
   - [saved_searches](#saved_searches)
6. [Trust & Safety](#6-trust--safety)
   - [reports](#reports)
   - [notifications](#notifications)
7. [Social](#7-social)
   - [user_follows](#user_follows)
8. [Relationships Summary](#8-relationships-summary)
9. [Design Decisions](#9-design-decisions)
10. [Suggested Index Strategy](#10-suggested-index-strategy)
11. [Implementation Order](#11-implementation-order)

---

## 1. Core Identity

### `campuses`

Represents a university or college. Users and listings are scoped to a campus.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `name` | `varchar(255)` | NOT NULL | e.g. "Ho Chi Minh City University of Technology" |
| `slug` | `varchar(100)` | UNIQUE, NOT NULL | e.g. `hcmut` |
| `domain` | `varchar(100)` | UNIQUE | e.g. `hcmut.edu.vn` — used for email auto-verification |
| `country` | `varchar(100)` | NOT NULL | |
| `city` | `varchar(100)` | NOT NULL | |
| `is_active` | `boolean` | DEFAULT true | Soft disable a campus |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `users`

A student, staff member, or admin on a campus. **Authentication is delegated to Clerk** (Phase 1 per `USE_CASES.md`); this row is the app’s profile + campus binding, not the password store.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `campus_id` | `uuid` | FK → campuses.id, NOT NULL | |
| `clerk_user_id` | `varchar(128)` | UNIQUE, NOT NULL | Clerk user id (`sub` claim) |
| `email` | `varchar(255)` | UNIQUE, NOT NULL | Mirrored from Clerk / registration |
| `full_name` | `varchar(255)` | NOT NULL | Display name |
| `avatar_url` | `text` | | Object storage URL (e.g. S3) |
| `role` | `varchar(20)` | DEFAULT `'student'` | `student`, `staff`, `admin` |
| `is_verified` | `boolean` | DEFAULT false | e.g. campus email validated in Go middleware |
| `is_active` | `boolean` | DEFAULT true | Soft delete / ban |
| `created_at` | `timestamptz` | DEFAULT now() | |
| `updated_at` | `timestamptz` | DEFAULT now() | |

Legacy self-hosted-auth designs may use `hashed_password`; with Clerk **omit passwords** entirely on this row.

---

## 2. Catalog

### `categories`

Self-referencing tree for nested categories (e.g. Electronics → Phones).

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `parent_id` | `uuid` | FK → categories.id, NULLABLE | NULL = top-level category |
| `name` | `varchar(100)` | NOT NULL | |
| `slug` | `varchar(100)` | UNIQUE, NOT NULL | |
| `icon_url` | `text` | | |
| `sort_order` | `integer` | DEFAULT 0 | Controls display order |
| `is_active` | `boolean` | DEFAULT true | |

---

### `listings`

The core entity — an item for sale by a student or staff member.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `seller_id` | `uuid` | FK → users.id, NOT NULL | |
| `category_id` | `uuid` | FK → categories.id, NOT NULL | |
| `campus_id` | `uuid` | FK → campuses.id, NOT NULL | Denormalized for fast filtering |
| `title` | `varchar(255)` | NOT NULL | |
| `description` | `text` | | |
| `price` | `numeric(12,2)` | NOT NULL | |
| `currency` | `varchar(10)` | DEFAULT `'VND'` | |
| `condition` | `varchar(20)` | NOT NULL | `new`, `like_new`, `good`, `fair`, `parts_only` |
| `status` | `varchar(20)` | DEFAULT `'draft'` | `draft`, `active`, `sold`, `archived`, `removed` |
| `view_count` | `integer` | DEFAULT 0 | |
| `is_negotiable` | `boolean` | DEFAULT false | |
| `location` | `point` | | PostGIS point for proximity search |
| `expires_at` | `timestamptz` | | Auto-archive when past this date |
| `created_at` | `timestamptz` | DEFAULT now() | |
| `updated_at` | `timestamptz` | DEFAULT now() | |

**Listing status lifecycle:**
```
draft → active → sold → archived
              ↘ removed (by admin/report)
```

---

### `listing_images`

Ordered images attached to a listing.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `listing_id` | `uuid` | FK → listings.id, NOT NULL | |
| `url` | `text` | NOT NULL | CDN URL |
| `cdn_key` | `varchar(500)` | | Storage key for deletion |
| `sort_order` | `integer` | DEFAULT 0 | |
| `is_primary` | `boolean` | DEFAULT false | Only one per listing |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `tags`

Freeform searchable tags shared across all listings.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `name` | `varchar(50)` | UNIQUE, NOT NULL | lowercase, trimmed |
| `usage_count` | `integer` | DEFAULT 0 | Denormalized counter |

---

### `listing_tags`

Join table linking listings to tags.

| Column | Type | Constraints |
|---|---|---|
| `listing_id` | `uuid` | FK → listings.id |
| `tag_id` | `uuid` | FK → tags.id |

**Primary key:** `(listing_id, tag_id)`

---

## 3. Messaging & Offers

### `conversations`

A chat thread between a buyer and seller about a specific listing.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `listing_id` | `uuid` | FK → listings.id, NOT NULL | |
| `buyer_id` | `uuid` | FK → users.id, NOT NULL | |
| `seller_id` | `uuid` | FK → users.id, NOT NULL | |
| `status` | `varchar(20)` | DEFAULT `'active'` | `active`, `closed`, `archived` |
| `last_message_at` | `timestamptz` | | For inbox sorting |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `messages`

Individual messages inside a conversation.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `conversation_id` | `uuid` | FK → conversations.id, NOT NULL | |
| `sender_id` | `uuid` | FK → users.id, NOT NULL | |
| `body` | `text` | NOT NULL | |
| `message_type` | `varchar(20)` | DEFAULT `'text'` | `text`, `image`, `offer`, `system` |
| `is_read` | `boolean` | DEFAULT false | |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `offers`

A price offer made by a buyer inside a conversation.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `conversation_id` | `uuid` | FK → conversations.id, NOT NULL | |
| `listing_id` | `uuid` | FK → listings.id, NOT NULL | |
| `buyer_id` | `uuid` | FK → users.id, NOT NULL | |
| `amount` | `numeric(12,2)` | NOT NULL | |
| `currency` | `varchar(10)` | DEFAULT `'VND'` | |
| `status` | `varchar(20)` | DEFAULT `'pending'` | `pending`, `accepted`, `rejected`, `withdrawn`, `expired` |
| `note` | `text` | | Optional buyer message |
| `expires_at` | `timestamptz` | | Auto-expire open offers |
| `created_at` | `timestamptz` | DEFAULT now() | |
| `updated_at` | `timestamptz` | DEFAULT now() | |

---

## 4. Commerce

### `transactions`

Created when an offer is accepted and a payment is initiated. Financial record of a completed sale.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `listing_id` | `uuid` | FK → listings.id, NOT NULL | |
| `buyer_id` | `uuid` | FK → users.id, NOT NULL | |
| `seller_id` | `uuid` | FK → users.id, NOT NULL | |
| `offer_id` | `uuid` | FK → offers.id, NULLABLE | NULL if no offer was made |
| `amount` | `numeric(12,2)` | NOT NULL | |
| `currency` | `varchar(10)` | DEFAULT `'VND'` | |
| `payment_method` | `varchar(50)` | | `cash`, `momo`, `vnpay`, `bank_transfer` |
| `payment_status` | `varchar(20)` | DEFAULT `'pending'` | `pending`, `completed`, `failed`, `refunded` |
| `external_tx_id` | `varchar(255)` | | ID from payment gateway |
| `metadata` | `jsonb` | | Raw payment gateway response |
| `completed_at` | `timestamptz` | | |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `reviews`

Mutual rating after a completed transaction. Both buyer and seller can review each other.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `transaction_id` | `uuid` | FK → transactions.id, NOT NULL | |
| `reviewer_id` | `uuid` | FK → users.id, NOT NULL | |
| `reviewee_id` | `uuid` | FK → users.id, NOT NULL | |
| `rating` | `smallint` | CHECK (1–5), NOT NULL | |
| `comment` | `text` | | |
| `created_at` | `timestamptz` | DEFAULT now() | |

**Unique constraint:** `(transaction_id, reviewer_id)` — one review per person per transaction.

---

## 5. Discovery

### `wishlists`

A named collection of saved listings for a user.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `user_id` | `uuid` | FK → users.id, NOT NULL | |
| `name` | `varchar(100)` | DEFAULT `'Saved items'` | |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

### `wishlist_items`

Individual listings inside a wishlist.

| Column | Type | Constraints |
|---|---|---|
| `wishlist_id` | `uuid` | FK → wishlists.id |
| `listing_id` | `uuid` | FK → listings.id |
| `added_at` | `timestamptz` | DEFAULT now() |

**Primary key:** `(wishlist_id, listing_id)`

---

### `saved_searches`

Stored search filters that can trigger notifications when a matching listing is posted.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `user_id` | `uuid` | FK → users.id, NOT NULL | |
| `name` | `varchar(100)` | | User-defined label |
| `filters` | `jsonb` | NOT NULL | e.g. `{"category":"phones","price_max":500000,"condition":"like_new"}` |
| `notify_on_match` | `boolean` | DEFAULT false | Trigger a push/email notification |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

## 6. Trust & Safety

### `reports`

Polymorphic report — can target a listing, user, or message.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `reporter_id` | `uuid` | FK → users.id, NOT NULL | |
| `target_type` | `varchar(20)` | NOT NULL | `listing`, `user`, `message` |
| `target_id` | `uuid` | NOT NULL | ID of the reported entity |
| `reason` | `varchar(50)` | NOT NULL | `spam`, `fraud`, `inappropriate`, `wrong_price`, `other` |
| `status` | `varchar(20)` | DEFAULT `'pending'` | `pending`, `reviewed`, `resolved`, `dismissed` |
| `note` | `text` | | Reporter's description |
| `resolved_by` | `uuid` | FK → users.id, NULLABLE | Admin who resolved |
| `created_at` | `timestamptz` | DEFAULT now() | |
| `resolved_at` | `timestamptz` | | |

---

### `notifications`

In-app notifications for any user event.

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `uuid` | PK | |
| `user_id` | `uuid` | FK → users.id, NOT NULL | |
| `type` | `varchar(50)` | NOT NULL | `new_message`, `offer_received`, `listing_sold`, `review_received`, etc. |
| `title` | `varchar(255)` | NOT NULL | |
| `body` | `text` | | |
| `data` | `jsonb` | | Contextual payload e.g. `{"listing_id":"...", "offer_id":"..."}` |
| `is_read` | `boolean` | DEFAULT false | |
| `created_at` | `timestamptz` | DEFAULT now() | |

---

## 7. Social

### `user_follows`

Lets a user follow a seller to get notified of their new listings.

| Column | Type | Constraints |
|---|---|---|
| `follower_id` | `uuid` | FK → users.id |
| `following_id` | `uuid` | FK → users.id |
| `created_at` | `timestamptz` |  DEFAULT now() |

**Primary key:** `(follower_id, following_id)`
**Constraint:** `follower_id != following_id`

---

## 8. Relationships Summary

```
campuses        ──< users
campuses        ──< listings
users           ──< listings         (seller)
users           ──< conversations    (buyer)
users           ──< messages
users           ──< reviews
users           ──< transactions
users           ──< wishlists
users           ──< notifications
users           ──< saved_searches
users           >──< users           (user_follows)
categories      ──< categories       (parent)
categories      ──< listings
listings        ──< listing_images
listings        >──< tags            (listing_tags)
listings        ──< conversations
conversations   ──< messages
conversations   ──< offers
offers          ──| transactions
transactions    ──| reviews
wishlists       ──< wishlist_items
listings        ──< wishlist_items
```

---

## 9. Design Decisions

| Decision | Rationale |
|---|---|
| UUID primary keys | Prevents enumeration attacks; supports sharding |
| `campus_id` on `listings` | Denormalized for fast campus-scoped filtering without joins |
| `jsonb` on `transactions.metadata` | Stores raw payment gateway responses without migration risk |
| `jsonb` on `saved_searches.filters` | Filter schema can evolve freely |
| `jsonb` on `notifications.data` | Attach any context (listing_id, offer_id, etc.) without new columns |
| Polymorphic `reports` | One reports table handles listings, users, and messages |
| Reviews tied to transactions | Both parties can rate; avoids fake reviews from non-buyers |
| Soft deletes everywhere | Full audit trail; needed for dispute resolution |
| `offers` separate from `transactions` | Keeps negotiation records clean and financial records separate |
| Self-referencing `categories` | Supports unlimited nesting depth without schema changes |

---

## 10. Suggested Index Strategy

```sql
-- Core lookups
CREATE INDEX idx_users_campus ON users(campus_id);
CREATE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_clerk_user_id ON users(clerk_user_id);

-- Listing filtering (most-used query path)
CREATE INDEX idx_listings_campus_status ON listings(campus_id, status);
CREATE INDEX idx_listings_category ON listings(category_id);
CREATE INDEX idx_listings_seller ON listings(seller_id);
CREATE INDEX idx_listings_expires ON listings(expires_at) WHERE status = 'active';

-- Messaging
CREATE INDEX idx_conversations_buyer ON conversations(buyer_id);
CREATE INDEX idx_conversations_seller ON conversations(seller_id);
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at);

-- Notifications inbox
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read, created_at);

-- Reports queue
CREATE INDEX idx_reports_status ON reports(status, created_at);

-- Full-text search on listings (PostgreSQL)
CREATE INDEX idx_listings_search ON listings USING gin(to_tsvector('english', title || ' ' || coalesce(description, '')));
```

---

## 11. Implementation Order

Build tables in this order to satisfy foreign key constraints:

1. `campuses`
2. `users`
3. `categories`
4. `listings`
5. `listing_images`, `tags`, `listing_tags`
6. `conversations`, `messages`, `offers`
7. `transactions`, `reviews`
8. `wishlists`, `wishlist_items`, `saved_searches`
9. `reports`, `notifications`
10. `user_follows`
