-- Per-user search profiles: a named record of who the user is professionally — one
-- specialization (a job category) and a non-empty set of skills. Unlike a saved
-- search (which stores a serialized filter query), a profile is structured: the
-- `specialization` is validated by the service against the category vocabulary
-- (enrich.CategoryValues) and `skills` holds canonical lowercase tokens. A user may
-- keep several profiles (e.g. "Go backend", "DevOps"), so names are unique per user
-- and length-bounded as a backstop to the service-layer validation. How a profile is
-- consumed (match scoring, ranked feeds, notifications) is intentionally out of scope
-- for this table. Applied automatically by Postgres on first volume init (same as
-- 0001) and also serves as schema source for sqlc. Existing volumes/prod need a manual
-- apply (the versioned-migration-runner seam from AGENT.md remains open).

CREATE TABLE IF NOT EXISTS search_profiles (
    id             BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id        BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    -- Display name, trimmed and bounded by the service; the CHECK is the backstop.
    name           TEXT        NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 100),
    -- One job category (validated against enrich.CategoryValues in the service, the
    -- same way the rest of the app validates its vocabularies — the table stays
    -- decoupled from a moving Go enum).
    specialization TEXT        NOT NULL CHECK (length(trim(specialization)) > 0),
    -- Canonical lowercase skill tokens; non-empty (a profile without skills has no
    -- meaning). The CHECK is the backstop to the service's normalization.
    skills         TEXT[]      NOT NULL CHECK (cardinality(skills) > 0),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- Distinct names per user so the profile picker has stable labels.
    UNIQUE (user_id, name)
);

-- List-by-owner, most-recently-updated first (the picker order).
CREATE INDEX IF NOT EXISTS search_profiles_user_updated_idx
    ON search_profiles (user_id, updated_at DESC);
