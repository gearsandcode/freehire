# AGENT.md

Guidance for AI agents working in this repository.

## Working principles

Non-negotiable. Bias toward caution over speed; use judgment on trivial tasks.

- **Think before coding.** Surface assumptions. If multiple interpretations exist, present them — don't pick silently. If something is unclear, ask.
- **Simplicity first.** Minimum code that solves the problem. No features, abstractions, or error handling that wasn't asked for. Prefer a library's intended API over a clever shim.
- **Surgical changes.** Touch only what the task requires; don't refactor unbroken things or rework formatting. Match existing style. Clean up what your change orphaned; leave pre-existing dead code alone. Exception: do the real refactor when a clean change genuinely requires reshaping existing code.
- **Fix root causes, not symptoms.**
- **No overengineering, and no MVP shortcuts.** Hold the middle path: don't build infrastructure before there's a concrete need (note the seam for later instead), and don't ship quick-and-dirty or "for now" hacks. Build each feature correctly and idiomatically — neither gold-plated nor a placeholder.
- **English only.** All code, comments, identifiers, docs, and commits are in English.

## What this is

`hire` is an open-source IT job aggregator backend. Intended shape: many source parsers feed a pipeline that normalizes jobs into one schema, deduplicates them, and enriches them with AI; served over an HTTP API with rich filters.

**Current state: backend scaffold only** — Fiber HTTP server, Postgres via sqlc, a minimal `jobs` table, and `/health` + `/api/v1/jobs[/:id]`. Parsers, pipeline, and AI layer do not exist yet.

Stack: **Go + Fiber v2**, **PostgreSQL**, **sqlc** (generated DB access, no ORM), **Docker Compose**.

## Layout

```
cmd/server/main.go   entry point: Fiber startup + graceful shutdown
internal/
  config/            env config (PORT, DATABASE_URL)
  database/          pgxpool connection pool
  db/                GENERATED sqlc code (do not edit) + queries/*.sql (hand-written)
  handlers/          HTTP handlers (Handler struct + Register wires routes)
migrations/          SQL schema — single source for BOTH sqlc and Postgres initdb
```

Future features slot in here without restructuring: `internal/sources/` (parsers as interface + registry), `internal/pipeline/` (fetch → normalize → dedup → upsert), `internal/ai/` (enrichment). `UpsertJob` already exists as the pipeline's write path.

## Commands

```bash
make up                      # build + start app and postgres in Docker
HIRE_HOST_PORT=8090 make up  # use another host port if 8080 is taken
make down / make logs        # stop containers / tail app logs
make run                     # run server on host (needs a running Postgres)
make psql                    # psql into the DB container
go build ./...  &&  go vet ./...
```

No test suite yet.

## Conventions and gotchas

- **sqlc is the only DB layer.** `internal/db/*.go` is generated — never edit by hand (committed so the repo builds without sqlc installed). To change DB access, edit `internal/db/queries/*.sql` (or `migrations/` for schema), run `make sqlc` (runs sqlc via Docker), and commit the result. Handlers use `*db.Queries`, built once in `handlers.Register`.
- **Migrations apply via Postgres initdb.** `migrations/` is mounted into `/docker-entrypoint-initdb.d`, so Postgres runs each `*.sql` **once, on first volume init only**. Changing a migration does NOT re-apply to an existing volume — recreate it with `docker compose down -v && make up`. The same dir is sqlc's schema source, keeping schema and code in sync. *Known seam:* no versioned migration runner yet; needed before the first schema change ships to a persistent DB.
- **Response shapes.** Lists: `{"data": ..., "meta": {...}}`; single items: `{"data": ...}`. Errors use `fiber.NewError(status, msg)` — no central `ErrorHandler` yet (deferred on purpose; don't hand-roll per-handler error JSON).
- **Dedup key.** `jobs.UNIQUE (source, external_id)` is the dedup key; `UpsertJob` is `ON CONFLICT` on it.
