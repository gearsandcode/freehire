## ADDED Requirements

### Requirement: Analysed-jobs list endpoint

The system SHALL provide an authenticated `GET /api/v1/me/tracking/analyses` endpoint that lists the jobs the caller has run the AI fit analysis on, newest first, without invoking the LLM. Each item MUST carry the job's public slug, title, company, a `closed` flag, the analysis `overall_score` and `verdict`, the analysed timestamp, and a `stale` flag (true when the caller's CV, the job content, or the model changed since the analysis was computed). The response MUST include the caller's fit-analysis `quota` (used/limit/remaining) in `meta`. The endpoint accepts a session cookie or an API key.

#### Scenario: List returns analysed jobs with quota

- **WHEN** a signed-in caller who has analysed two jobs requests `GET /api/v1/me/tracking/analyses`
- **THEN** the response is `{ "data": [<two items newest first>], "meta": { "quota": { "used": 2, "limit": 10, "remaining": 8 } } }`, each item carrying slug/title/company/closed/overall_score/verdict/analysed-at/stale, and no LLM call is made

#### Scenario: Closed analysed job is retained with a flag

- **WHEN** the caller analysed a job that has since closed
- **THEN** it still appears in the list with `closed: true`

#### Scenario: Stale analysis is flagged

- **WHEN** the caller's CV was re-uploaded after an analysis was computed
- **THEN** that item is returned with `stale: true`

### Requirement: Tracking routes with back-compat aliases

The per-user tracking endpoints SHALL be served canonically under `/api/v1/me/tracking` (`""`, `/viewed`, `/pipeline`, `/swipe`, `/analyses`), and the previous `/api/v1/me/jobs*` paths MUST continue to work as aliases to the same handlers so already-released API clients (the freehire-cli) are not broken.

#### Scenario: Canonical tracking path

- **WHEN** a caller requests `GET /api/v1/me/tracking`
- **THEN** it returns the caller's tracked jobs exactly as the legacy `/api/v1/me/jobs` did

#### Scenario: Legacy alias still works

- **WHEN** an existing client requests `GET /api/v1/me/jobs`
- **THEN** the system serves the same response as `GET /api/v1/me/tracking` (no breakage)

### Requirement: Tracking section renamed with URL redirects

The frontend personal-jobs section SHALL be presented as **Tracking** and served under `/my/tracking/*` (Board, Pipeline, History, AI fit). Requests to the previous `/my/jobs/*` URLs MUST redirect (HTTP 308) to the corresponding `/my/tracking/*` path so existing bookmarks and inbound links keep working.

#### Scenario: Old URL redirects to the new section

- **WHEN** a user opens `/my/jobs/pipeline`
- **THEN** the app redirects to `/my/tracking/pipeline`

#### Scenario: Section labelled Tracking

- **WHEN** a signed-in user opens the tracking section
- **THEN** the navigation and heading read "Tracking", with tabs for Board, Pipeline, History, and AI fit
