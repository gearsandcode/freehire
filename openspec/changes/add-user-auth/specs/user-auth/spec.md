## ADDED Requirements

### Requirement: User registration

The system SHALL allow a new user to register with an email and password,
creating exactly one account per email and returning a signed authentication
token on success.

- Email MUST be unique (case-insensitive); the stored form is lowercased.
- Password MUST be at least 8 characters; it is stored only as a bcrypt hash,
  never in plaintext and never returned in any response.
- On success the system returns the created user (id, email, created_at) and a
  signed JWT.

#### Scenario: Successful registration

- **WHEN** a client POSTs a unique, well-formed email and an 8+ character password to `/api/v1/auth/register`
- **THEN** the system creates the user, stores a bcrypt hash of the password, and responds `201` with the user (no password hash) and a signed JWT

#### Scenario: Duplicate email

- **WHEN** a client registers with an email that already exists (in any letter case)
- **THEN** the system responds `409` and creates no new account

#### Scenario: Invalid input

- **WHEN** a client submits a malformed email or a password shorter than 8 characters
- **THEN** the system responds `400` and creates no account

### Requirement: User login

The system SHALL authenticate an existing user by email and password and return
a signed authentication token, without revealing whether the email or the
password was the failing factor.

#### Scenario: Successful login

- **WHEN** a client POSTs a registered email and the correct password to `/api/v1/auth/login`
- **THEN** the system responds `200` with the user and a signed JWT

#### Scenario: Wrong password

- **WHEN** a client submits a registered email with an incorrect password
- **THEN** the system responds `401` with a generic "invalid credentials" message and issues no token

#### Scenario: Unknown email

- **WHEN** a client submits an email that has no account
- **THEN** the system responds `401` with the same generic "invalid credentials" message as a wrong password

#### Scenario: Account has no password

- **WHEN** a client attempts password login for an account that has no stored password hash (e.g. one created through a future passwordless sign-in method)
- **THEN** the system responds `401` with the same generic "invalid credentials" message, never treating an absent password as a match

### Requirement: Stateless token authentication

The system SHALL issue stateless JWTs (HS256) on register and login, and SHALL
validate them on protected requests via the `Authorization: Bearer <token>`
header.

- The token SHALL encode the user id as its subject and carry an expiry.
- A protected handler MUST be able to resolve the authenticated user's id from
  the validated token.

#### Scenario: Valid token grants access

- **WHEN** a client calls a protected endpoint with a valid, unexpired Bearer token
- **THEN** the system resolves the user from the token and serves the request

#### Scenario: Missing or malformed token

- **WHEN** a client calls a protected endpoint without an `Authorization: Bearer` header or with a malformed token
- **THEN** the system responds `401` and does not serve the protected resource

#### Scenario: Expired or invalid signature

- **WHEN** a client calls a protected endpoint with an expired token or one whose signature does not verify against the server secret
- **THEN** the system responds `401`

### Requirement: Current user endpoint

The system SHALL expose `GET /api/v1/auth/me` that returns the authenticated
user's profile and is only reachable with a valid token.

#### Scenario: Authenticated request

- **WHEN** an authenticated client calls `GET /api/v1/auth/me`
- **THEN** the system responds `200` with the user (id, email, created_at) and never includes the password hash

#### Scenario: Unauthenticated request

- **WHEN** a client calls `GET /api/v1/auth/me` without a valid token
- **THEN** the system responds `401`

### Requirement: Web client authentication

The Svelte SPA SHALL let a user register, log in, and log out from the
application layout, persist the session across reloads, and reflect the current
auth state in the top bar.

- The auth token SHALL be persisted in `localStorage` and re-loaded on boot; a
  stored token SHALL be validated on startup via `GET /me`, and a token that
  fails validation SHALL be discarded so the user appears signed out.
- Authenticated API requests SHALL attach the token as `Authorization: Bearer`;
  the public jobs/companies requests SHALL remain unauthenticated.
- The top bar SHALL show the signed-in user's email and a logout action when
  authenticated, and Login/Register actions when not.

#### Scenario: Sign in from the layout

- **WHEN** a signed-out user submits valid credentials in the login (or register) form opened from the top bar
- **THEN** the SPA stores the returned token, shows the user's email with a logout action in the top bar, and keeps the user signed in across a page reload

#### Scenario: Log out

- **WHEN** a signed-in user activates the logout action
- **THEN** the SPA clears the stored token and returns the top bar to its Login/Register state

#### Scenario: Stale token on boot

- **WHEN** the SPA boots with a stored token that the server rejects (expired or invalid)
- **THEN** the SPA discards the token and presents the signed-out state without error

### Requirement: Public endpoints remain unauthenticated

The existing read endpoints SHALL remain publicly accessible without a token, so
this change adds authentication without gating current functionality.

#### Scenario: Public read without a token

- **WHEN** a client calls `GET /api/v1/jobs`, `GET /api/v1/jobs/:id`, `GET /api/v1/companies`, or `GET /api/v1/companies/:slug` without any token
- **THEN** the system serves the request as before, unaffected by the auth layer
