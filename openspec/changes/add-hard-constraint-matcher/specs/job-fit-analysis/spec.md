## MODIFIED Requirements

### Requirement: Five-dimension scored verdict

The analysis payload SHALL contain five dimensions — Title & role alignment, Experience
relevance, Seniority fit, Skills coverage, and Company & role context — each with an integer
score clamped to 0–100, plus a weighted `overall_score`, a `verdict` label drawn from the
controlled set {Strong Fit, Good Fit, Moderate Fit, Weak Fit, Poor Fit}, a `strengths` array,
a `gaps` array, and a single `recommendation` string. All model output MUST be sanitized: scores
clamped, the verdict coerced to the controlled set, and free-text fields trimmed and length/count
bounded before it is persisted or served. The server-computed `overall_score` MUST additionally be
clamped down to the deterministic hard-constraint ceiling (the minimum score-cap over the caller's
unmet blockers) when any blocker is present, so a posting the caller plainly cannot meet can never
present as a strong fit; the `verdict` label is derived from the capped `overall_score`.

#### Scenario: Out-of-range or invalid model output

- **WHEN** the LLM returns a dimension score above 100 or a verdict outside the controlled set
- **THEN** the score is clamped to 0–100 and the verdict is derived from `overall_score`, so no out-of-vocabulary value is ever persisted or served

#### Scenario: Overall score is the weighted dimensions

- **WHEN** the five dimension scores are known and the caller has no unmet hard-constraint blockers
- **THEN** `overall_score` equals the deterministic weighted average of the dimensions, computed server-side rather than trusting the model's own overall

#### Scenario: Hard-constraint ceiling caps an over-optimistic score

- **WHEN** the weighted average is 88 but the caller has an unmet certification blocker with score-cap 60
- **THEN** the served `overall_score` is 60 and the `verdict` label is derived from that capped value

### Requirement: Per-(user, job) cache with staleness invalidation

The system SHALL cache each analysis per `(user_id, job_id)`, stamped with the CV's upload time,
the job's `content_hash`, the model that produced it at analysis time, and the hard-constraint
dictionary version in effect at analysis time. `GET /api/v1/jobs/:slug/fit`
MUST return a cached analysis only when all four stamps still equal the current CV upload time, job
`content_hash`, model, and hard-constraint dictionary version; when any differs it MUST report the cached
analysis as stale rather than serving it as current. A `content_hash` absent on both the stored stamp
and the live job (a non-board job that is never re-crawled) counts as unchanged, so it does not force
an endless recompute.

#### Scenario: Fresh cache hit

- **WHEN** a user GETs the fit for a job they already analyzed, and neither their CV, the job, the model, nor the hard-constraint dictionary has changed since
- **THEN** the system returns the cached analysis with `stale: false` and makes no LLM call

#### Scenario: Model upgraded since analysis

- **WHEN** a user GETs the fit for a job analyzed under a previous `LLM_MODEL`
- **THEN** the cached analysis is reported with `stale: true` so the improved model can re-analyze on request

#### Scenario: CV changed since analysis

- **WHEN** a user GETs the fit after re-uploading their CV
- **THEN** the cached analysis is reported with `stale: true` so the SPA can offer a recompute, and it is not served as current

#### Scenario: Job re-ingested with changed content

- **WHEN** a user GETs the fit for a job whose `content_hash` changed since the analysis
- **THEN** the cached analysis is reported with `stale: true`

#### Scenario: Hard-constraint dictionary updated since analysis

- **WHEN** a user GETs the fit for a job analyzed under a previous hard-constraint dictionary version
- **THEN** the cached analysis is reported with `stale: true` so the recomputed cap reflects the current dictionary

#### Scenario: No analysis yet

- **WHEN** a user GETs the fit for a job they have never analyzed
- **THEN** the system responds `200` with `has_cv` reflecting CV presence and a null analysis (no LLM call)

## ADDED Requirements

### Requirement: Hard-constraint blockers ground the prompt chain

The prompt chain SHALL include the deterministic hard-constraint blockers as known, already-established constraints so the model explains and respects them rather than re-deriving degree, years, license, or work-authorization requirements. The served analysis MUST expose the blockers alongside the verdict.

#### Scenario: Blockers passed into the prompt and surfaced

- **WHEN** the fit analysis is computed for a caller with an unmet hard constraint
- **THEN** the prompt carries the blocker as a known constraint and the served analysis exposes it beside the dimensions
