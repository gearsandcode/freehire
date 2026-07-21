## ADDED Requirements

### Requirement: Required certifications in the enrichment payload

The enrichment payload SHALL carry a `required_certifications` field: a list of canonical credential slugs the posting requires, drawn from the shared credential vocabulary's controlled set. The field MUST pass the same sanitize/validate gate as other enrichment fields — an out-of-vocabulary slug is dropped, never persisted or served — and it MUST live in the existing enrichment jsonb with no new database column and no Meilisearch facet.

#### Scenario: Recognized credential is captured as a canonical slug

- **WHEN** a posting states it requires an "AWS Solutions Architect certification"
- **THEN** enrichment records the canonical slug for that credential in `required_certifications`

#### Scenario: Unknown credential string is dropped

- **WHEN** the model emits a credential slug outside the controlled vocabulary
- **THEN** the sanitize/validate gate removes it and nothing invalid is persisted

### Requirement: Education level omitted on equivalent-experience wording

When a posting offers a degree "or equivalent experience" (or equivalent wording), the enrichment extraction SHALL leave `education_level` unset rather than recording a degree requirement, so that a downstream hard-constraint check does not raise a false education blocker. The suppression is applied at extraction time, not by a downstream text patch.

#### Scenario: "Degree or equivalent" yields no education requirement

- **WHEN** a posting says "Bachelor's degree or equivalent experience"
- **THEN** enrichment leaves `education_level` unset
