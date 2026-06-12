## Why

The ingest pipeline supports six ATS platforms (greenhouse, lever, ashby, workable,
recruitee, smartrecruiters), but a live audit of open ATS company datasets
(`kalil0321/ats-scrapers`, MIT) shows the reachable long tail is far larger. Several
more platforms expose a public, no-auth posting feed and together cover tens of
thousands of additional companies we cannot ingest today. Adding adapters for them is
the highest-leverage way to widen catalogue coverage without changing the pipeline.

## What Changes

- Register five new `Source` adapters, each verified live against the platform's public
  feed, so boards on these platforms can be listed in `sources.yml`:
  - **personio** — public XML feed (`{board}.jobs.personio.com/xml`); description inline, single request.
  - **breezy** — public JSON (`{board}.breezy.hr/json`); description inline, single request.
  - **pinpoint** — public JSON (`{board}.pinpointhq.com/postings.json`); description inline, single request.
  - **rippling** — public JSON (`api.rippling.com/platform/api/ats/v1/board/{board}/jobs`); list lacks description, per-posting detail fetch required.
  - **bamboohr** — public JSON (`{board}.bamboohr.com/careers/list` → `/careers/{id}/detail`); list lacks description, per-posting detail fetch required.
- Add a **join.com** adapter (`join.com`, ~23k companies — the largest single source).
  Its feed is not a plain REST endpoint: postings are served via a GraphQL
  `candidate-api` / Next.js `__NEXT_DATA__`. This change first pins the exact request,
  then builds the adapter against a captured fixture.
- Seed a small set of live-validated boards per new provider into `sources.yml` so each
  adapter ingests real postings and tests exercise the real response shape.

## Capabilities

### New Capabilities
<!-- none — this extends the existing source-ingest capability -->

### Modified Capabilities
- `source-ingest`: add a requirement registering `personio`, `breezy`, `pinpoint`,
  `rippling`, `bamboohr`, and `join.com` as providers, each yielding the normalized job
  shape with a sanitized-HTML description, consistent with the existing adapters
  (single-request where the list carries the body; bounded per-posting detail fetch
  where it does not).

## Impact

- **Code**: new `internal/sources/{personio,breezy,pinpoint,rippling,bamboohr,joincom}.go`
  + table-driven `_test.go` per adapter; one registration line each in `sources.All`
  (`internal/sources/source.go`); seed entries in `sources.yml`.
- **Pipeline**: none — adapters slot into the existing registry/`Source` interface; no
  change to `pipeline`, the write path, or `cmd/ingest`.
- **Dependencies**: XML parsing for personio uses the stdlib `encoding/xml`; no new
  third-party dependency. join.com may need a small captured GraphQL query string.
- **Out of scope (seams)**: `gem`, `jazzhr`, `recruiterbox` — no clean public feed found
  (client-rendered / redirected); they need feed discovery (or headless rendering) first.
  Mass slug harvest + live-probe tooling and splitting `sources.yml` into per-provider
  files remain separate, provider-agnostic changes.
