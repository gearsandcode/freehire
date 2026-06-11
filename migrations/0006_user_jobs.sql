-- Per-user job interactions: one row per (user, job) pair recording when the
-- user viewed a job and (optionally) when they marked it applied.
-- Applied automatically by Postgres on first volume init (same as 0001) and also
-- serves as schema source for sqlc.

CREATE TABLE IF NOT EXISTS user_jobs (
    user_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    job_id     BIGINT      NOT NULL REFERENCES jobs (id)  ON DELETE CASCADE,
    -- Most-recent view: RecordJobView refreshes this on every revisit. Add a
    -- separate first_viewed_at later if a timeline ever needs the original.
    viewed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- NULL = viewed but not applied. Set when the user confirms an application;
    -- a non-null value is the entry point for the future stage pipeline.
    applied_at TIMESTAMPTZ,
    -- The composite key is also the dedup key: at most one interaction (and so at
    -- most one application) per (user, job).
    PRIMARY KEY (user_id, job_id)
);
