## 1. Database

- [ ] 1.1 Add `migrations/0005_users.sql`: `users` table (id identity PK, email text not null, password_hash text **NULLABLE** — passwordless users have none, created_at timestamptz default now()) with `UNIQUE (lower(email))`
- [ ] 1.2 Add `internal/db/queries/users.sql` with `CreateUser` (returns id, email, created_at), `GetUserByEmail` (incl. password_hash, for login), `GetUserByID` (no password_hash, for /me)
- [ ] 1.3 Run `make sqlc` and commit the regenerated `internal/db` code
- [ ] 1.4 Recreate the dev DB volume (`docker compose down -v && make up`) to apply the new migration

## 2. Dependencies & Config

- [ ] 2.1 Add `golang.org/x/crypto/bcrypt` and `github.com/golang-jwt/jwt/v5` to go.mod (`go get`)
- [ ] 2.2 Add `JWTSecret` and `JWTTTL` to `config.Settings`; load `JWT_SECRET` (no default) and `JWT_TTL` (default 24h)
- [ ] 2.3 Fail fast at server startup if `JWT_SECRET` is empty

## 3. Auth package (`internal/auth`)

- [ ] 3.1 Implement `HashPassword(plain) (string, error)` and `CheckPassword(hash, plain) error` over bcrypt
- [ ] 3.2 Implement an `Issuer` (secret + TTL): `Issue(userID) (string, error)` producing an HS256 JWT with `sub` + `exp`, and `Parse(token) (userID, error)` validating signature and expiry
- [ ] 3.3 Implement `RequireAuth` Fiber middleware: parse `Authorization: Bearer`, validate via `Issuer`, store user id in `c.Locals`, return 401 on any failure
- [ ] 3.4 Unit tests: hash round-trip + wrong-password rejection; issue→parse round-trip; expired token and bad-signature rejection

## 4. HTTP handlers (`internal/handler/auth.go`)

- [ ] 4.1 Add an API user type (id, email, created_at) that never includes password_hash, and a JSON-tag-omitted mapping from the db row
- [ ] 4.2 `Register` handler: validate email format + password length (>=8), lowercase email, hash, `CreateUser`, map duplicate (unique violation) to 409, return 201 with `{"data": {user, token}}`
- [ ] 4.3 `Login` handler: `GetUserByEmail`, `CheckPassword`, return generic 401 on unknown email OR wrong password OR account with a null `password_hash` (passwordless), else 200 with `{"data": {user, token}}`
- [ ] 4.4 `Me` handler: read user id from `c.Locals`, `GetUserByID`, return 200 `{"data": user}`
- [ ] 4.5 Wire routes in `handler.Register`: extend its signature with JWT secret/TTL, build the `auth.Issuer`, register `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, and `GET /api/v1/auth/me` (guarded by `RequireAuth`)
- [ ] 4.6 Update `cmd/server/main.go` to pass the new config into `handler.Register`

## 6. Web (SPA) auth integration

- [ ] 6.1 Add auth API to `web/src/lib/api.ts`: `register`, `login`, `me` functions and a way to attach `Authorization: Bearer` (keep public jobs/companies calls unauthenticated); add the `User` wire type to `web/src/lib/types.ts`
- [ ] 6.2 Add `web/src/lib/auth.svelte.ts` auth store mirroring `theme.svelte.ts`: `$state` token + user, persist token in localStorage (`hire.token`), `login`/`register`/`logout` actions, and an `initAuth()` that validates a stored token via `GET /me` on boot (discard on failure); call `initAuth()` in `main.ts`
- [ ] 6.3 Add login + register form components (`web/src/lib/components/`), reachable from the layout (route or modal, matching the existing router pattern), surfacing API errors (e.g. 401/409) inline
- [ ] 6.4 Wire auth controls into `TopBar.svelte`: signed-in shows user email + logout; signed-out shows Login/Register — placed alongside `ThemeToggle` in the `ml-auto` group

## 7. Verification

- [ ] 7.1 `go build ./... && go vet ./... && go test ./...` all pass
- [ ] 7.2 Manual e2e (API): register → receive token; login → receive token; `GET /me` with token → 200; without/expired token → 401; confirm `GET /api/v1/jobs` still works with no token
- [ ] 7.3 Manual e2e (web): register/login from the top bar → email + logout shown; reload stays signed in; logout returns to Login/Register; a tampered/expired stored token boots to signed-out without error
- [ ] 7.4 Update `AGENT.md` (layout + conventions) to document the `internal/auth` package, the auth endpoints, the `JWT_SECRET`/`JWT_TTL` env vars, and the SPA auth store
