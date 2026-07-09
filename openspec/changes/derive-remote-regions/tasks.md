## 1. Derive remote_regions in the recompute

- [ ] 1.1 RED: extend the recompute integration test (`internal/db`,
      `//go:build integration`) — seed a company with an open `remote` job in one
      region and an open `onsite` job in another; assert after `RefreshCompanyFacets`
      that `remote_regions` is the remote job's region only, while `regions` is the
      union. Flip the old "leaves curated remote_regions untouched" guard.
- [ ] 1.2 GREEN: add a remote-scoped `remote_reg` CTE to `RefreshCompanyFacets`
      (`internal/db/queries/companies.sql`) — `array_agg(DISTINCT r) ... FROM open
      jobs WHERE work_mode='remote' CROSS JOIN unnest(regions)`; add `remote_regions`
      to the SET list and to the `IS DISTINCT FROM` change-guard. Regenerate via `make sqlc`.
- [ ] 1.3 REFACTOR + simplify; recompute tests stay green.

## 2. Remove the curated backfill machinery

- [ ] 2.1 Delete `internal/remoteregion/` (package + tests).
- [ ] 2.2 Delete `cmd/backfill-remote-regions/` (worker + tests).
- [ ] 2.3 Delete `sources/remote-companies.csv`.
- [ ] 2.4 Remove the `SetCompanyRemoteRegions` query from
      `internal/db/queries/companies.sql`; regenerate via `make sqlc`. Delete the
      curated-only integration test (`company_remote_regions_integration_test.go`),
      keeping/moving the recompute coverage into task 1.
- [ ] 2.5 `go build ./... && go vet ./...` — confirm no dangling references.

## 3. Docs

- [ ] 3.1 Rewrite the AGENT.md convention + layout/commands entries: `remote_regions`
      is now a job-derived facet (remote-scoped `regions`) maintained by the
      recompute; drop the `internal/remoteregion` / `cmd/backfill-remote-regions` /
      `sources/remote-companies.csv` references.

## 4. Finish + deploy

- [ ] 4.1 `go build ./... && go vet ./... && go test ./...` and the `internal/db`
      integration tests green.
- [ ] 4.2 Deploy (host-2 `release.sh freehire`), then run `cmd/recount-companies`
      once so `remote_regions` repopulates from remote jobs; one-off
      `UPDATE companies SET company_info = company_info - 'remote_regions_raw'` to
      drop the stale audit field. Verify the facet live.
