-- name: ListJobs :many
SELECT *
FROM jobs
ORDER BY posted_at DESC NULLS LAST, id DESC
LIMIT $1 OFFSET $2;

-- name: GetJob :one
SELECT *
FROM jobs
WHERE id = $1;

-- name: CountJobs :one
SELECT count(*)
FROM jobs;

-- name: ListJobsByCompany :many
SELECT *
FROM jobs
WHERE company_slug = $1
ORDER BY posted_at DESC NULLS LAST, id DESC
LIMIT $2 OFFSET $3;

-- name: UpsertJob :one
-- Single atomic write: upsert the company (only when the slug is non-empty,
-- via the WHERE on the SELECT) and the job together, keeping the "one write =
-- one job" property of the pipeline's write path.
-- NOTE: enrichment must be a non-nil json.RawMessage (pass []byte("{}") for an
-- un-enriched job, never nil) — the column is NOT NULL and the '{}' default does
-- not apply to an explicit NULL on INSERT.
WITH company_upsert AS (
    INSERT INTO companies (slug, name)
    SELECT sqlc.arg(company_slug), sqlc.arg(company)
    WHERE sqlc.arg(company_slug) <> ''
    ON CONFLICT (slug) DO UPDATE SET
        name       = EXCLUDED.name,
        updated_at = now()
)
INSERT INTO jobs (
    source, external_id, url, title, company, company_slug, location, remote, description, posted_at,
    enrichment, enriched_at, enrichment_version
) VALUES (
    sqlc.arg(source), sqlc.arg(external_id), sqlc.arg(url), sqlc.arg(title),
    sqlc.arg(company), sqlc.arg(company_slug), sqlc.arg(location), sqlc.arg(remote),
    sqlc.arg(description), sqlc.arg(posted_at),
    sqlc.arg(enrichment), sqlc.arg(enriched_at), sqlc.arg(enrichment_version)
)
ON CONFLICT (source, external_id) DO UPDATE SET
    url          = EXCLUDED.url,
    title        = EXCLUDED.title,
    company      = EXCLUDED.company,
    company_slug = EXCLUDED.company_slug,
    location     = EXCLUDED.location,
    remote       = EXCLUDED.remote,
    description  = EXCLUDED.description,
    posted_at    = EXCLUDED.posted_at,
    -- Full-replace, consistent with the raw fields above. Seam for phase 2: when
    -- the ingest path (which carries no enrichment) and the enrichment path are
    -- separated, decide whether a source re-ingest preserves existing enrichment.
    enrichment         = EXCLUDED.enrichment,
    enriched_at        = EXCLUDED.enriched_at,
    enrichment_version = EXCLUDED.enrichment_version,
    updated_at   = now()
RETURNING *;

-- name: SetJobEnrichment :exec
-- Targeted enrichment write used by the enrichment command: set only the payload
-- and the provenance stamp, touching no raw source field. Kept separate from
-- UpsertJob (the ingest full-upsert path) so ingest and enrichment stay decoupled.
UPDATE jobs
SET enrichment         = sqlc.arg(enrichment),
    enriched_at        = sqlc.arg(enriched_at),
    enrichment_version = sqlc.arg(enrichment_version),
    updated_at         = now()
WHERE id = sqlc.arg(id);
