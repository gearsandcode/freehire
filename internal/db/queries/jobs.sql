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

-- name: UpsertJob :one
INSERT INTO jobs (
    source, external_id, url, title, company, location, remote, description, posted_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (source, external_id) DO UPDATE SET
    url         = EXCLUDED.url,
    title       = EXCLUDED.title,
    company     = EXCLUDED.company,
    location    = EXCLUDED.location,
    remote      = EXCLUDED.remote,
    description = EXCLUDED.description,
    posted_at   = EXCLUDED.posted_at,
    updated_at  = now()
RETURNING *;
