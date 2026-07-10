// Command backfill-descriptions repairs job descriptions that were stored still
// percent-encoded. The Taleo adapter decoded its description HTML with the strict
// url.PathUnescape, which rejects the whole string on a single stray "%" — common in
// Word-pasted postings (CSS "line-height:115%") — and the old fallback stored the raw,
// fully percent-encoded blob (rendered as literal "%3Cp class=%22..."). The adapter now
// decodes leniently (internal/sources.LenientPercentUnescape); this one-off worker fixes the
// rows already in the catalogue.
//
// It pages by keyset and re-decodes every row whose description still carries the "%3C"
// (encoded "<") marker, re-running the same sanitize+decode pipeline the fixed adapter uses and
// refreshing content_hash so the row re-indexes. The marker is source-agnostic: any
// percent-encoded description is repaired the same way, open or closed. Idempotent — a
// re-decoded row no longer matches the marker, so a second run rewrites nothing.
//
// With no argument it sweeps the whole catalogue (universal). Pass a source name as the first
// argument (e.g. `backfill-descriptions taleo`) to scope the sweep to that provider — the
// affected rows are known to be a single provider's, and scoping skips detoasting every other
// row, so a targeted repair does not read the whole table.
//
// Follow it with `make reindex` so the search/recommendation index picks up the fixed text.
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/strelov1/freehire/internal/db"
	"github.com/strelov1/freehire/internal/jobhash"
	"github.com/strelov1/freehire/internal/sources"
	"github.com/strelov1/freehire/internal/worker"
)

// backfillBatchSize bounds how many rows are read per keyset page.
const backfillBatchSize = 500

// encodedMarker is the fingerprint of a still-percent-encoded description: the encoded "<"
// ("%3C") that opens every mangled HTML blob. A cleanly decoded description never contains it,
// so it precisely selects the rows to repair without a content-scanning SQL predicate.
const encodedMarker = "%3C"

// jobStore is the slice of the data layer the backfill needs: page rows by keyset (whole table or
// one source) and rewrite a row's description + content_hash. *db.Queries satisfies it; the test
// uses a fake.
type jobStore interface {
	ListJobsByIDAfter(ctx context.Context, arg db.ListJobsByIDAfterParams) ([]db.Job, error)
	ListJobsBySourceAfter(ctx context.Context, arg db.ListJobsBySourceAfterParams) ([]db.Job, error)
	UpdateJobDescription(ctx context.Context, arg db.UpdateJobDescriptionParams) (int64, error)
}

func main() {
	worker.Main(run)
}

func run() int {
	ctx, _, pool, cleanup, err := worker.Bootstrap(context.Background())
	if err != nil {
		log.Printf("database: %v", err)
		return 1
	}
	defer cleanup()

	source := ""
	if len(os.Args) > 1 {
		source = os.Args[1]
	}

	scanned, updated, err := backfillAll(ctx, db.New(pool), source)
	if err != nil {
		log.Printf("backfill-descriptions: %v", err)
		return 1
	}
	scope := "all sources"
	if source != "" {
		scope = "source=" + source
	}
	log.Printf("backfill-descriptions done (%s): scanned=%d updated=%d (run `make reindex` to refresh the index)", scope, scanned, updated)
	return 0
}

// backfillAll pages jobs by keyset (id > last seen, so concurrent writes cannot skip or repeat
// rows) and re-decodes the ones whose stored description still carries the encoded marker. An
// empty source sweeps the whole table; a non-empty source scopes the sweep to that provider. The
// decode reproduces the fixed adapter's pipeline exactly (LenientPercentUnescape then
// SanitizeHTML), so the recomputed content_hash matches what a re-ingest would produce.
func backfillAll(ctx context.Context, store jobStore, source string) (scanned, updated int, err error) {
	var afterID int64
	for {
		jobs, err := pageJobs(ctx, store, source, afterID)
		if err != nil {
			return scanned, updated, err
		}
		if len(jobs) == 0 {
			break
		}
		afterID = jobs[len(jobs)-1].ID

		for _, j := range jobs {
			scanned++
			if !strings.Contains(j.Description, encodedMarker) {
				continue
			}
			desc := sources.SanitizeHTML(sources.LenientPercentUnescape(j.Description))
			if desc == j.Description {
				continue // nothing recovered — leave it be (defensive; a marker row always changes)
			}
			hash := jobhash.Of(hashParams(j, desc))
			if _, err := store.UpdateJobDescription(ctx, db.UpdateJobDescriptionParams{
				ID:          j.ID,
				Description: desc,
				ContentHash: pgtype.Text{String: hash, Valid: true},
			}); err != nil {
				return scanned, updated, err
			}
			updated++
		}

		if len(jobs) < backfillBatchSize {
			break
		}
	}
	return scanned, updated, nil
}

// pageJobs fetches one keyset page after afterID — from the whole table when source is empty, or
// scoped to a single provider otherwise (which lets the query use the (source, ...) index and
// skip detoasting every other row).
func pageJobs(ctx context.Context, store jobStore, source string, afterID int64) ([]db.Job, error) {
	if source == "" {
		return store.ListJobsByIDAfter(ctx, db.ListJobsByIDAfterParams{
			AfterID:   afterID,
			BatchSize: backfillBatchSize,
		})
	}
	return store.ListJobsBySourceAfter(ctx, db.ListJobsBySourceAfterParams{
		Source:    source,
		AfterID:   afterID,
		BatchSize: backfillBatchSize,
	})
}

// hashParams builds the content_hash inputs for a row with a replaced description — the exact
// indexed fields jobhash.Of fingerprints (see internal/jobhash), so the recomputed hash matches
// what the ingest write path would produce for the same row.
func hashParams(j db.Job, description string) db.UpsertJobParams {
	return db.UpsertJobParams{
		URL:                j.URL,
		Title:              j.Title,
		Company:            j.Company,
		CompanySlug:        j.CompanySlug,
		Location:           j.Location,
		Remote:             j.Remote,
		Description:        description,
		PostedAt:           j.PostedAt,
		PublicSlug:         j.PublicSlug,
		Countries:          j.Countries,
		Regions:            j.Regions,
		WorkMode:           j.WorkMode,
		Skills:             j.Skills,
		Seniority:          j.Seniority,
		Category:           j.Category,
		PostingLanguage:    j.PostingLanguage,
		EmploymentType:     j.EmploymentType,
		EducationLevel:     j.EducationLevel,
		ExperienceYearsMin: j.ExperienceYearsMin,
	}
}
