## 1. Personio adapter (single-request, XML)

- [ ] 1.1 Add a `GetXML` method to `HTTPClient` (mirror `GetJSON`, stdlib `encoding/xml`) with a test
- [ ] 1.2 Capture a real `…jobs.personio.com/xml` response (a board with ≥1 `<position>`) as a test fixture
- [ ] 1.3 Write a failing table-driven `personio_test.go`: maps `<position>` → `Job`, description from `jobDescriptions` decoded + sanitized, empty feed → zero jobs no error
- [ ] 1.4 Implement `internal/sources/personio.go` (provider `personio`) until tests pass; register in `sources.All`

## 2. Breezy adapter (single-request, JSON)

- [ ] 2.1 Capture a real `…breezy.hr/json` response as a fixture
- [ ] 2.2 Write a failing `breezy_test.go`: maps posting → `Job`, inline body sanitized, empty list → zero jobs no error
- [ ] 2.3 Implement `internal/sources/breezy.go` (provider `breezy`) until tests pass; register in `sources.All`

## 3. Pinpoint adapter (single-request, JSON)

- [ ] 3.1 Capture a real `…pinpointhq.com/postings.json` response as a fixture
- [ ] 3.2 Write a failing `pinpoint_test.go`: maps `data[]` → `Job`, inline body sanitized, empty data → zero jobs no error
- [ ] 3.3 Implement `internal/sources/pinpoint.go` (provider `pinpoint`) until tests pass; register in `sources.All`

## 4. Rippling adapter (list + per-posting detail)

- [ ] 4.1 Capture a real board list (`api.rippling.com/.../board/{board}/jobs`) and one posting detail as fixtures
- [ ] 4.2 Write a failing `rippling_test.go`: list mapped, per-`uuid` detail fetched for description, description sanitized, detail fan-out is bounded
- [ ] 4.3 Implement `internal/sources/rippling.go` (provider `rippling`, reuse the smartrecruiters bounded-concurrency pattern) until tests pass; register in `sources.All`

## 5. BambooHR adapter (list + per-posting detail)

- [ ] 5.1 Capture a real `…/careers/list` and one `…/careers/{id}/detail` as fixtures
- [ ] 5.2 Write a failing `bamboohr_test.go`: list mapped, per-`id` detail fetched for description, location from `joinNonEmpty(city,state,country)`, bounded fan-out
- [ ] 5.3 Implement `internal/sources/bamboohr.go` (provider `bamboohr`) until tests pass; register in `sources.All`

## 6. Join.com adapter (discovery-then-build)

- [ ] 6.1 Capture join.com's postings request for a board (prefer the GraphQL `candidate-api` operation; fall back to `__NEXT_DATA__`); freeze the response as a fixture. If no stable, ToS-clean request exists, drop join.com from this change and note it in the proposal/design
- [ ] 6.2 Write a failing `joincom_test.go` against the captured fixture: postings mapped, description sanitized, empty → zero jobs no error
- [ ] 6.3 Implement `internal/sources/joincom.go` (provider `join.com`) until tests pass; register in `sources.All`

## 7. Seed boards and verification

- [ ] 7.1 Add a few live-validated boards per new provider to `sources.yml` (each confirmed to return ≥1 posting)
- [ ] 7.2 `go build ./... && go vet ./... && go test ./...` all green
- [ ] 7.3 Run `go run ./cmd/ingest` against a dev DB and confirm postings from the new providers are ingested with sanitized descriptions
