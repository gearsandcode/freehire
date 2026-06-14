## 1. Pin current behavior (characterization)

- [ ] 1.1 Add a handler-level characterization test that drives the current view/apply/save/unsave/track endpoints and captures the exact JSON body + status for representative cases: viewed-only, applied, saved, tracked (stage+notes), track stage-only, track notes-only, invalid stage (400), empty track body (400), unsave-when-absent (zero interaction), unknown slug (404).
- [ ] 1.2 Run it against the current (pre-refactor) handler and confirm it passes — this is the green baseline the refactor must preserve.

## 2. jobtracking package — contract and rules (TDD, fake repo)

- [ ] 2.1 Create `internal/jobtracking` with the domain `Interaction` struct (`JobID int64`, `ViewedAt/SavedAt/AppliedAt *time.Time`, `Stage/Notes *string`) and the sentinel errors `ErrJobNotFound`, `ErrInvalidStage`, `ErrEmptyTrack`.
- [ ] 2.2 Define the narrow `Repository` interface (`JobIDBySlug`, `RecordView`, `MarkApplied`, `SaveJob`, `UnsaveJob`, `TrackJob`) in the package.
- [ ] 2.3 Write failing unit tests for `Service` against a fake `Repository`: record view, mark applied, save, each returning the mapped `Interaction`.
- [ ] 2.4 Implement `Service.RecordView/MarkApplied/SaveJob` (resolve slug→id via repo; `ErrJobNotFound` on unknown slug) to pass 2.3.
- [ ] 2.5 Write failing tests for `Unsave` idempotency (no row → zero `Interaction{JobID}`, nil error) and implement it.
- [ ] 2.6 Write failing tests for `Track`: invalid stage → `ErrInvalidStage`; neither stage nor notes → `ErrEmptyTrack`; stage-only / notes-only / both leave the other field unchanged. Implement `Service.Track` (using `userjob.ValidStage`) to pass.

## 3. Persistence adapter

- [ ] 3.1 Add a `Repository` adapter in `internal/jobtracking` wrapping `*db.Queries`, converting `db.UserJob` (pgtype) → `Interaction` and `db` errors → the package sentinels (`pgx.ErrNoRows` → `ErrJobNotFound` for slug lookup; no-row on unsave handled in the service).
- [ ] 3.2 Confirm the package builds and `go test ./internal/jobtracking/` is green.

## 4. Thin the handler

- [ ] 4.1 Construct the service in `handler.Register` (`jobtracking.New(jobtracking.NewQueriesRepository(queries))`) and store it on `Handler`.
- [ ] 4.2 Rewrite `RecordView/MarkApplied/SaveJob/UnsaveJob/TrackJob` to: read `auth.UserID`, parse the body (track only), call the service with `(userID, slug, …)`, map domain errors (`ErrJobNotFound`→404, `ErrInvalidStage`/`ErrEmptyTrack`→400), and render `Interaction` → the existing `interactionResponse` JSON.
- [ ] 4.3 Remove the now-dead inline rules from `user_jobs.go` (`validStages` is already gone; remove `interactionParams`/`trackRequest` plumbing superseded by the service, keeping only what transport still needs). The handler no longer calls `h.queries` for these endpoints.

## 5. Verify

- [ ] 5.1 Run the characterization test from 1.1 against the refactored handler — must be byte-identical (no wire change).
- [ ] 5.2 `go build ./...`, `go vet ./...`, `go test ./...` all green.
- [ ] 5.3 Self-review: confirm no `pgtype`/`pgx`/`fiber` import in `internal/jobtracking`, and no remaining business rule (stage check, partial-update, idempotency, slug resolution) in `user_jobs.go`.
