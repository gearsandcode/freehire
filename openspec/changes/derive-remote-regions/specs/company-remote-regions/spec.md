## REMOVED Requirements

### Requirement: Companies carry a curated remote-hiring-regions facet

**Reason**: `remote_regions` is no longer a curated backfilled column; it is now a
job-derived facet (union of `regions` over open remote jobs), owned by the facet
recompute. See the `companies` capability's modified derived-facets and recompute
requirements.
**Migration**: The `companies.remote_regions` column and its `remote_regions`
filter facet are retained unchanged; only the data source moves from the backfill
to `RefreshCompanyFacets`. No consumer change. The stale
`company_info.remote_regions_raw` audit field is dropped by a one-off update.

### Requirement: A curated dataset maps company names to remote-hiring region strings

**Reason**: The signal is derived from our own remote postings, so the external
`sources/remote-companies.csv` directory is no longer an input.
**Migration**: Delete `sources/remote-companies.csv`. Nothing consumes it after
the backfill worker is removed.

### Requirement: A pure dictionary maps region strings to macro-region codes

**Reason**: Region derivation for remote hiring now reuses the existing
`jobs.regions` (produced by `internal/location` at ingest); a separate free-text
mapping dictionary is unnecessary.
**Migration**: Delete `internal/remoteregion`. Job region codes already come from
`internal/location` / `enrich.RegionValues`.

### Requirement: The backfill annotates existing companies only, by slug

**Reason**: No backfill worker exists; `remote_regions` is maintained by the
periodic recompute across all companies, not by a slug-matched external load.
**Migration**: Delete `cmd/backfill-remote-regions` and the `SetCompanyRemoteRegions`
query. Values are (re)computed on the next `cmd/recount-companies` run.
