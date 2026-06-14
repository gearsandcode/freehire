-- Orphan-job liveness (see openspec change probe-orphan-job-liveness).
--
-- liveness_strikes counts CONSECUTIVE "expired" URL probes for jobs the ingest
-- sweep never re-crawls — the non-board sources (telegram, habr_career, geekjob)
-- whose closed_at would otherwise stay NULL forever. The liveness worker
-- increments it on a definitive death signal and closes the job (sets closed_at)
-- once it reaches the threshold; any healthy probe resets it to 0. Board jobs keep
-- the default 0: they are closed by the ingest sweep, not by URL probes.
ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS liveness_strikes INT NOT NULL DEFAULT 0;
