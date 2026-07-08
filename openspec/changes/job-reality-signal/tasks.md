## 1. Repost fingerprint (schema + write path)

- [ ] 1.1 Add migration: `jobs.role_fingerprint text` + partial index on `(company_slug, role_fingerprint) WHERE closed_at IS NULL` (backs the concurrent mass-posting count); note dev volume recreate + prod manual-apply-before-deploy in the migration comment.
- [ ] 1.2 Add a narrow fingerprint function (in `internal/jobhash` or a new `internal/rolefingerprint`) over `company_slug` + normalized title + normalized description, **excluding** `posted_at`/url/slug; unit tests: bumped `posted_at` → same fingerprint, differing title/description → different fingerprint.
- [ ] 1.3 Wire the fingerprint into the write path: add `role_fingerprint` to `UpsertJob` in `internal/db/queries/jobs.sql`, run `make sqlc`, and set it from the normalize/ingest path where `content_hash` is computed. Verify it is NOT `content_hash`.
- [ ] 1.4 Integration test (`//go:build integration`): a repost under a new `external_id` with refreshed `posted_at` resolves to the same `role_fingerprint` as the original.

## 2. Reality classifier (pure `internal/jobreality`)

- [ ] 2.1 Curated evergreen-text dictionary (EN + RU full surface forms: "always hiring", "talent community", "building a pipeline", RU equivalents) + whole-word/phrase matcher; tests: known phrase matches, unmatched text emits nothing (never guesses).
- [ ] 2.2 Pure classifier: input `{now, createdAt, postedAt, closedAt, repostCount, massPostingCount, evergreenText}` → `{Class, Evidence{ageDays, repostCount, massPostingCount, fakeFreshness}}`. Rule: `fresh` if age ≤ freshWindow and no evergreen signal; `likely-evergreen` only when ≥ N of {old-age, repost≥k, massPosting≥m, evergreenText} converge; else `stale`.
- [ ] 2.3 Classifier tests covering every spec scenario: fresh, stale, age-alone-is-not-evergreen, convergence-is-evergreen, fake-freshness recorded when `postedAt` recent but `createdAt` old, determinism.

## 3. Derive, serve, and index

- [ ] 3.1 sqlc queries for the two counts per `(company_slug, role_fingerprint)`: distinct `external_id`s of any status (repost history) and of open jobs only (concurrent mass-posting).
- [ ] 3.2 Derive `reality` where the other facets are derived (`internal/jobderive` / the reindex build); include it in `cmd/backfill-derive`'s single pass. Also backfill `role_fingerprint` for existing rows (one-off) so counts are meaningful.
- [ ] 3.3 `internal/jobview`: add a top-level `reality: {class, ageDays, repostCount, massPostingCount, fakeFreshness}` field, served dict-only (never from `enrichment.*`); unit test `FromRow`.
- [ ] 3.4 `internal/search`: index `reality.class` as a filterable Meili facet attribute in `FromJob` + index settings; note the "new attribute 500s `/jobs` until first reindex" window.

## 4. Web surface

- [ ] 4.1 Regenerate TS contracts for the new `jobview` field (`cmd/gen-contracts`) so the SPA type matches.
- [ ] 4.2 Reality badge on the job card + detail: `fresh` hidden/neutral, `stale`/`likely-evergreen` render **facts** ("open 240 days · reposted 6×"), not a bare accusation; verify visually.
- [ ] 4.3 Add a `reality` filter chip to the filter modal (opt-in include/exclude; **not** hidden by default).

## 5. Verify and reconcile

- [ ] 5.1 `go build ./... && go vet ./... && go test ./...` green; run the integration test with the DB tag.
- [ ] 5.2 Verification-before-completion: exercise the end-to-end path (ingest a reposted job → derive → search facet → badge) and confirm the class + evidence render; document the ship order (apply migration → deploy → backfill `role_fingerprint` → `cmd/backfill-derive` → `make reindex`).
