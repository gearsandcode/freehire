-- Structured, AI-derived enrichment for jobs. Additive: raw source columns on
-- jobs are untouched. Applied automatically by Postgres on first volume init
-- (same as 0001/0002) and also serves as schema source for sqlc.
--
-- The enriched fields live in one JSONB blob rather than typed columns:
-- filtering will be served by a derived search index (planned Meilisearch),
-- not Postgres, so the blob maps 1:1 to a future search document and lets the
-- field set evolve without a migration per field. The typing discipline lives
-- in the Go contract (internal/enrich), not the database.
ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS enrichment         JSONB       NOT NULL DEFAULT '{}',
    -- Provenance: NULL enriched_at means never enriched; enrichment_version lets
    -- a later enrichment job select rows below the current schema version to
    -- re-run. Columns (not in the blob) because that job queries them directly.
    ADD COLUMN IF NOT EXISTS enriched_at        TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS enrichment_version INT         NOT NULL DEFAULT 0;
