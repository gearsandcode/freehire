## Why

Our fit signal is split between a deterministic skill-coverage bar (`internal/jobmatch`, instant/free) and an LLM match-analysis (`internal/matchanalysis`, opt-in). Neither enforces **hard constraints** — a posting that requires 5+ years, a specific degree, a license/certification, work authorization, or on-site presence can still score high when the caller plainly does not meet it, because skill coverage ignores those axes and the LLM's score drifts run-to-run and cannot be trusted as a guardrail. Both sides of the comparison are already structured (job enrichment carries `experience_years_min`/`education_level`/`english_level`/`visa_sponsorship`/`countries`/`work_mode`; the résumé is parsed into `total_years`/`education[]`/`languages`/`location`), so a cheap, explainable, deterministic checker can surface real blockers and cap an over-optimistic score — no new LLM calls on the hot path.

## What Changes

- Add `internal/hardconstraint` — a pure, deterministic, dict-only package (same discipline as `jobmatch`/`classify`) that evaluates a job's structured requirements against the caller's structured résumé across six categories (experience-years, education, language, work-authorization, location-and-work-mode, certification) and emits typed `Blocker{Category, Severity, ScoreCap, Reason, Action, Met}`.
- Add a shared curated **credential vocabulary** (alias→canonical, IT-first + global; ported/trimmed from JobSentinel) and a **degree ladder / equivalence** dictionary, both under the new package.
- Add `required_certifications []string` to the job enrichment payload (canonical credential slugs), and instruct the enrichment prompt to **omit `education_level` when the posting offers "or equivalent experience"** (suppression at the source, not a downstream regex).
- Add `Certifications []string` to the résumé structured shape.
- Surface blockers on three surfaces: (1) the instant `jobmatch` bar renders blockers beside skill coverage; (2) `matchanalysis` caps the server-owned `overall_score` by `min(ScoreCap)` over unmet blockers, injects the known blockers into its Stage-1 and Stage-3 context, and stamps the analysis cache with the hard-constraint dictionary version; (3) `tailor` consumes the anti-hallucination `Action` strings.
- Precision rule #1 — **never a false blocker**: a category is evaluated only when both sides carry data; a missing enrichment field or absent résumé structure skips that category silently. Language is info-only by default (CV languages rarely carry a CEFR level); location/work-mode is soft.
- No behaviour is gated/hidden — blockers warn and cap only; the job stays visible and clickable.

## Capabilities

### New Capabilities
- `hard-constraint-matching`: the deterministic evaluator, `Blocker` model, six requirement categories, severity/score-cap tiers, credential vocabulary, degree ladder, and the "never a false blocker" precision rules.

### Modified Capabilities
- `job-enrichment`: adds the `required_certifications` payload field (controlled credential vocabulary) and the "omit `education_level` on 'or equivalent'" extraction rule.
- `resume-structured-profile`: adds the `certifications` field to the structured résumé shape.
- `job-profile-match`: the profile-match bar additionally surfaces hard-constraint blockers alongside skill coverage.
- `job-fit-analysis`: `overall_score` is capped by the deterministic blockers, the blockers are injected into the prompt chain, and the analysis cache is invalidated on hard-constraint dictionary-version change.
- `cv-tailoring`: the tailor flow consumes the blockers' anti-hallucination `Action` strings as guardrails.

## Impact

- **New package:** `internal/hardconstraint` (+ credential/degree dictionaries) and its generated TS `Blocker` wire shape via `cmd/gen-contracts`.
- **Schema (no DB migration):** `required_certifications` lives in the existing `jobs.enrichment` jsonb and `certifications` in the `users.resume_structured` jsonb — no DDL. `required_certifications` stays **out of Meilisearch facets** in v1, so **no reindex**.
- **Rollout:** the enrichment field populates on new jobs automatically; the existing catalogue is re-enriched in waves via `cmd/enrich` (an empty field simply skips the certification category meanwhile). The résumé field self-heals on the next CV upload. Both are additive and degrade gracefully.
- **Touched code:** `internal/enrich` (add-field: unmarshal + Sanitize/Validate + prompt), `internal/resumeextract` (add-field + prompt), `internal/jobmatch` or its handler (blocker surface), `internal/matchanalysis` (score-cap + prompt context + cache stamp), the tailor context assembly, and the relevant HTTP handlers/wire shapes.
