-- Durable work queue for job enrichment: one row per (job_id, target_version)
-- that needs enriching. Reference-only — job_id + version + bookkeeping, NOT a
-- copy of the job; `jobs` stays canonical. A future ingest path will insert here
-- in the same transaction as the job upsert (transactional outbox); until then
-- rows are enqueued by an idempotent backfill from jobs' provenance columns.
--
-- Applied by Postgres on first volume init (same as 0001-0003) and read by sqlc.
CREATE TABLE IF NOT EXISTS enrichment_outbox (
    id             BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id         BIGINT      NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    -- The enrich.Version this entry should be brought up to. Lets a version bump
    -- enqueue re-enrichment without colliding with an old, already-done entry.
    target_version INT         NOT NULL,
    -- Bookkeeping for the lease/retry/dead-letter machinery.
    attempts       INT         NOT NULL DEFAULT 0,
    claimed_at     TIMESTAMPTZ,            -- lease stamp; NULL = unleased
    failed_at      TIMESTAMPTZ,            -- non-NULL = dead-lettered, never reclaimed
    last_error     TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- One entry per job per target version; makes the backfill enqueue idempotent.
    UNIQUE (job_id, target_version)
);

-- The claim scans only live (not dead-lettered) work, oldest first.
CREATE INDEX IF NOT EXISTS enrichment_outbox_claimable_idx
    ON enrichment_outbox (id)
    WHERE failed_at IS NULL;
