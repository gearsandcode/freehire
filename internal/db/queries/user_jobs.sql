-- name: RecordJobView :one
-- Record (or refresh) a user's view of a job. Idempotent on (user_id, job_id):
-- the first view creates the row, a repeat view touches viewed_at. Returns the
-- row so the caller learns the current applied_at in the same round-trip.
INSERT INTO user_jobs (user_id, job_id)
VALUES ($1, $2)
ON CONFLICT (user_id, job_id) DO UPDATE SET viewed_at = now()
RETURNING *;

-- name: MarkJobApplied :one
-- Mark a job as applied for a user. Idempotent and independent of a prior view:
-- it inserts the row (viewed_at defaults) or updates applied_at in place.
INSERT INTO user_jobs (user_id, job_id, applied_at)
VALUES ($1, $2, now())
ON CONFLICT (user_id, job_id) DO UPDATE SET applied_at = now()
RETURNING *;
