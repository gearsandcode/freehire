## Context

`hire` is an early backend: a Fiber v2 server over Postgres (sqlc, no ORM),
serving read-only `jobs` and `companies` endpoints to a Svelte SPA that runs on
a separate origin (`localhost:5173` → API `8080`, already CORS-enabled). There
is no user concept anywhere. This change introduces the first authenticated
surface. Per project conventions, sqlc is the only DB layer, migrations apply
via Postgres initdb, and response shapes follow `{"data": ...}`.

## Goals / Non-Goals

**Goals:**
- A `users` identity with secure password storage (bcrypt).
- Stateless JWT auth that works cleanly cross-origin for the existing SPA.
- `register`, `login`, `me` endpoints plus a reusable "require auth" middleware
  that future protected routes can adopt without further wiring.
- Stay within existing conventions (sqlc, config-from-env, `handler.Register`).

**Non-Goals:**
- Roles/authorization tiers (admin vs user) — not needed until a mutating or
  privileged endpoint exists. Noted as a seam.
- Refresh tokens / token revocation / logout — out of scope; the trade-off is
  accepted below.
- Email verification, password reset, OAuth, magic-link — explicitly excluded
  this iteration (OAuth/magic-link are the announced *next* task; the data model
  is shaped to absorb them additively, but none of their tables or flows are
  built here).
- Gating any existing read endpoint behind auth.

## Decisions

### Stateless JWT (HS256) over server-side sessions

Tokens are signed with a shared secret (`JWT_SECRET`) and carry the user id as
`sub` plus an `exp`. The SPA stores the token and sends it as
`Authorization: Bearer <token>`.

*Why over DB-backed sessions:* the SPA is cross-origin, so cookie sessions would
require `SameSite=None`, `Secure`, and CORS `AllowCredentials` — more moving
parts and a CSRF surface. A bearer header sidesteps all of that and needs no new
table. *Alternative considered:* opaque session IDs in a `sessions` table
(revocable, but adds a table, a lookup per request, and the cookie/CORS
complexity). Rejected for MVP.

*Library:* `github.com/golang-jwt/jwt/v5` — the de-facto Go JWT library;
maintained, v5 has the safer parsing API.

### bcrypt for password hashing

`golang.org/x/crypto/bcrypt` with the default cost. The salt and cost are
embedded in the hash string, so the `users` table needs a single
`password_hash` column — no separate salt column. *Alternative:* argon2id
(stronger, PHC winner) but requires hand-tuning salt/time/memory params; bcrypt
is the simpler, proven default for an MVP.

### New `internal/auth` package, separate from handlers

A focused package owns the security primitives behind small interfaces:
- password hashing/verification (`HashPassword`, `CheckPassword`),
- token issue/verify (`Issuer` wrapping secret + TTL: `Issue(userID)` /
  `Parse(token) → userID`),
- a Fiber middleware `RequireAuth` that validates the bearer token and stores
  the user id in `c.Locals`.

Handlers (`internal/handler/auth.go`) stay thin: parse/validate input, call
sqlc + `internal/auth`, shape the `{"data": ...}` response. This keeps crypto
and token logic testable in isolation (unit tests with no DB) and the security
boundary easy to hold in context. `handler.Register` grows parameters for the
JWT secret/TTL (mirroring how `frontendOrigin` is already threaded in).

### Data model

`migrations/0005_users.sql`:
```
id            bigint generated always as identity primary key
email         text not null
password_hash text                       -- NULLABLE: passwordless users have none
created_at    timestamptz not null default now()
unique (lower(email))   -- case-insensitive uniqueness
```
Queries in `internal/db/queries/users.sql`: `CreateUser`, `GetUserByEmail`,
`GetUserByID`. The `password_hash` column is selected only where needed
(login/registration) and never serialized into a handler response — the API
user type omits it.

`password_hash` is **nullable** on purpose (see "Forward-compatibility" below):
the announced next iteration adds Google OAuth and magic-link sign-in, where a
user has no password at all. Relaxing `NOT NULL` later would mean a second
migration — and this project has no migration runner yet (changing `0005` does
not re-apply to a live volume), so a `NOT NULL → NULL` change on a persistent DB
is genuinely painful. Getting the nullability right now is free. The password
login path treats a row with a null hash as "this account has no password" and
rejects it with the same generic `401` as a wrong password.

### Email as the canonical account key

`UNIQUE (lower(email))` makes email the one identity per account. This is the
deliberate foundation for "one account, multiple ways to authenticate":
password today; Google and magic-link later all resolve to (or link against) the
same email-keyed user. The JWT carries only the user id (`sub`), so it is
already provider-agnostic — no token, middleware, or `/me` change is needed when
new sign-in methods land.

### Forward-compatibility for OAuth / magic-link (seam, not built here)

The model is shaped so the next iteration is purely additive — nothing in
`users` needs reworking:
- **Google OAuth / external identities** → a future `user_identities` table
  (`user_id`, `provider`, `provider_user_id`, ...). Linking, not a `users`
  change.
- **Magic link** → a future short-lived login-token table; passwordless, so it
  relies on the nullable `password_hash` above.
- **Email verification** → a future additive `email_verified` column (or derived
  from a verified identity).

Building any of these now would be infrastructure ahead of need — explicitly
deferred (see Non-Goals). The only thing done *now* is removing the single
narrowing barrier (`password_hash NOT NULL`); everything else is added later
without touching existing rows.

### Config

Add `JWTSecret` and `JWTTTL` to `config.Settings`. `JWT_SECRET` has no safe
default — startup MUST fail fast if it is empty, so a server never boots with a
guessable signing key. `JWT_TTL` defaults to a sensible value (e.g. 24h).

## Risks / Trade-offs

- **No token revocation** (stateless JWT) → a leaked or post-logout token stays
  valid until `exp`. Mitigation: keep TTL modest (24h); the refresh/revocation
  design is a known seam — if revocation becomes a requirement, introduce a
  short access TTL + a `refresh_tokens` table without changing the public
  contract.
- **Single shared secret** → rotating it invalidates all live tokens.
  Acceptable at MVP; mitigation is fail-fast on empty secret so it is always set
  deliberately via env.
- **No rate limiting on login** → brute-force surface. Out of scope here;
  generic `401` (no email/password distinction) and bcrypt's cost slow attacks.
  Note the seam for a future rate-limit middleware.
- **No migration runner yet** (existing project gotcha) → `0005_users.sql`
  applies only on a fresh Postgres volume; an existing dev volume needs
  `docker compose down -v && make up`. Same constraint as all current
  migrations; documented, not solved here.

## Migration Plan

1. Add `0005_users.sql`; recreate the dev DB volume to apply it.
2. Add the new go deps, `users.sql` queries, run `make sqlc`, commit generated
   code.
3. Implement `internal/auth`, then handlers, then wire `Register`.
4. Set `JWT_SECRET` in the environment / compose before running the server.

Rollback: the change is additive (new table, new routes, new package). Reverting
the code and dropping the `users` table fully removes it; no existing data or
endpoint is touched.

## Open Questions

- None blocking. Role-based authorization and refresh tokens are deferred by
  decision, not left open.
