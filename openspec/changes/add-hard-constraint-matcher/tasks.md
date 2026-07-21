## 1. Foundation & dictionaries

- [ ] 1.1 Explore call sites: confirm whether the profile-match bar handler already loads the caller's structured résumé + job enrichment, and where `matchanalysis` computes `overall_score` and assembles the Stage-1/Stage-3 prompt and cache stamps. Record findings; resolve the design's open questions (experience tolerance = strict; bar degradation path).
- [ ] 1.2 Create `internal/hardconstraint/credentials`: curated credential vocabulary (alias→canonical slug, IT-first + global), ported/trimmed from JobSentinel's credential groups. Provide `Canonical(alias string) (slug string, ok bool)` and the controlled slug set. Table-test alias→canonical + unknown→drop.
- [ ] 1.3 Add the degree ladder + equivalence dictionary (none<ged<associate<bachelor<master<phd) with `DegreeRank(name string) (rank int, ok bool)`. Table-test equivalence normalization and ranking.

## 2. Core evaluator (internal/hardconstraint)

- [ ] 2.1 Define the wire types in a gen-contracts–visible file: `Blocker{Category, Severity, ScoreCap, Reason, Action, Met}`, the `Category`/`Severity` enums, and the `JobRequirements`/`CVEvidence` input structs. Add `Sanitize`/bounds if any strings are model-derived.
- [ ] 2.2 Implement `Evaluate(job, cv) []Blocker` with the experience and education categories first (numeric compare; degree ladder). RED→GREEN table tests: unmet→blocker, met→satisfied, missing-either-side→skipped, strict tolerance.
- [ ] 2.3 Add the certification category (canonical-slug intersection via the credential vocabulary) and the language category (info-only by default, never a blocker without a comparable level). Table-test met/unmet/skip.
- [ ] 2.4 Add the work-authorization and location-and-work-mode categories (visa_sponsorship + countries vs résumé location / location_preferences; work_mode conflict). Soft severities. Table-test conflict vs geo-uncertainty→skip.
- [ ] 2.5 Implement the severity/score-cap tiers and the `min(ScoreCap)` overall-cap helper. Table-test hardest-blocker-wins and no-unmet→no-cap.

## 3. Schema add-fields (structured-first)

- [ ] 3.1 Enrichment: add `required_certifications []string` (unmarshal + Sanitize/Validate against the credential vocabulary). Add the field to the enrichment prompt and instruct it to emit canonical credential slugs. Test sanitize drops out-of-vocab slugs.
- [ ] 3.2 Enrichment prompt: instruct the extractor to leave `education_level` unset on "degree or equivalent experience" wording. Add/adjust a prompt test asserting the instruction is present.
- [ ] 3.3 Résumé: add `Certifications []string` to `resumeextract.Structured` (wire + Sanitize bounds) and to the extraction prompt. Test sanitize bounds; confirm stale/absent degrades.

## 4. Integration surfaces

- [ ] 4.1 Profile-match bar: assemble `JobRequirements`/`CVEvidence` at the handler and attach `blockers` to the payload when structured inputs are present; degrade to coverage-only otherwise. Handler/integration test both paths; never hide the job.
- [ ] 4.2 matchanalysis: compute blockers before the chain; clamp the server-owned `overall_score` to `min(ScoreCap)`; derive `verdict` from the capped score. Test the cap scenario (weighted 88 → capped 60).
- [ ] 4.3 matchanalysis: inject the blockers into the Stage-1 and Stage-3 prompt context as known constraints; expose blockers in the served analysis. Prompt-content test + analysis-shape test.
- [ ] 4.4 matchanalysis: add the hard-constraint dictionary version as a fourth cache stamp; report `stale: true` when it differs. Cache-staleness test for the new stamp.
- [ ] 4.5 Tailor: pass unmet blockers' `Action` strings into the tailor context as explicit "do not claim unless true" guardrails. Test the action strings reach the tailor input.

## 5. Contracts & frontend

- [ ] 5.1 Run `cmd/gen-contracts`; regenerate TS for `Blocker`, the new enrichment/résumé fields, and the fit payload. Verify the generated diff is limited to these.
- [ ] 5.2 Surface blockers in the profile-match bar UI (advisory warning chips + met ✓), and show the ceiling context on the fit analysis page. Frontend unit test for the pure reducer/formatter; visual check.

## 6. Rollout & verification

- [ ] 6.1 Verify no DB migration is needed (both fields serialize into existing jsonb) and that `required_certifications` is not added to any Meilisearch facet (no reindex). Confirm `go build ./... && go vet ./...` and full `go test ./...` are green.
- [ ] 6.2 Document the post-deploy catalogue re-enrich in waves via `cmd/enrich` (empty field = certification category skipped meanwhile) in the change notes; no backfill blocks the deploy.
- [ ] 6.3 `verification-before-completion`: drive the fit analysis end-to-end for a caller who misses a hard constraint and confirm the served score is capped and the blocker is surfaced.
