package enrich

import (
	"context"
	"encoding/json"
	"fmt"
)

// Claimed is one outbox entry leased to this run.
type Claimed struct {
	OutboxID      int64
	JobID         int64
	TargetVersion int
	Attempts      int
}

// Store is the persistence the runner needs, expressed in domain terms so the
// runner is independent of the DB layer. The real implementation wraps the
// generated queries and a connection pool (and runs Complete in a transaction);
// tests use an in-memory fake.
type Store interface {
	// Enqueue adds outbox entries for jobs not yet enriched to targetVersion.
	Enqueue(ctx context.Context, targetVersion int) (int64, error)
	// Claim leases up to batch live, unleased entries.
	Claim(ctx context.Context, batch, leaseSeconds int) ([]Claimed, error)
	// Job returns the source fields a Provider reads for the given job id.
	Job(ctx context.Context, id int64) (JobInput, error)
	// Complete writes the enrichment payload + provenance stamp to the job and
	// deletes the outbox entry, atomically.
	Complete(ctx context.Context, outboxID, jobID int64, payload json.RawMessage, version int) error
	// Fail records a failed attempt; it returns whether the entry was dead-lettered.
	Fail(ctx context.Context, outboxID int64, errMsg string, maxAttempts int) (deadLettered bool, err error)
}

// RunOptions are the per-run knobs.
type RunOptions struct {
	TargetVersion int
	BatchSize     int
	LeaseSeconds  int
	MaxAttempts   int
}

// Stats reports what a run did.
type Stats struct {
	Enriched     int
	Failed       int
	DeadLettered int
}

// Runner drives the enrichment process: enqueue pending jobs, then drain claimed
// batches, enriching and writing back each entry.
type Runner struct {
	Provider Provider
	Store    Store
}

// Run enqueues pending jobs and drains the queue until no claimable entries remain.
// A failure on a single entry is recorded and never aborts the run.
func (r Runner) Run(ctx context.Context, opt RunOptions) (Stats, error) {
	if _, err := r.Store.Enqueue(ctx, opt.TargetVersion); err != nil {
		return Stats{}, fmt.Errorf("enqueue: %w", err)
	}

	var stats Stats
	for {
		batch, err := r.Store.Claim(ctx, opt.BatchSize, opt.LeaseSeconds)
		if err != nil {
			return stats, fmt.Errorf("claim: %w", err)
		}
		if len(batch) == 0 {
			return stats, nil
		}
		for _, entry := range batch {
			r.process(ctx, entry, opt, &stats)
		}
	}
}

// process handles one claimed entry. Any failure routes to recordFailure so the
// run continues with the remaining entries.
func (r Runner) process(ctx context.Context, entry Claimed, opt RunOptions, stats *Stats) {
	job, err := r.Store.Job(ctx, entry.JobID)
	if err != nil {
		r.recordFailure(ctx, entry, opt, stats, fmt.Errorf("load job: %w", err))
		return
	}

	enr, err := r.enrich(ctx, job)
	if err != nil {
		r.recordFailure(ctx, entry, opt, stats, err)
		return
	}

	payload, err := json.Marshal(enr)
	if err != nil {
		r.recordFailure(ctx, entry, opt, stats, fmt.Errorf("marshal: %w", err))
		return
	}

	if err := r.Store.Complete(ctx, entry.OutboxID, entry.JobID, payload, opt.TargetVersion); err != nil {
		r.recordFailure(ctx, entry, opt, stats, fmt.Errorf("write back: %w", err))
		return
	}
	stats.Enriched++
}

// enrich asks the provider for a payload and validates it, retrying once. An
// invalid payload is treated as an error so it is never persisted.
func (r Runner) enrich(ctx context.Context, job JobInput) (Enrichment, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		enr, err := r.Provider.Enrich(ctx, job)
		if err != nil {
			lastErr = err
			continue
		}
		if err := enr.Validate(); err != nil {
			lastErr = err
			continue
		}
		return enr, nil
	}
	return Enrichment{}, lastErr
}

func (r Runner) recordFailure(ctx context.Context, entry Claimed, opt RunOptions, stats *Stats, cause error) {
	dead, err := r.Store.Fail(ctx, entry.OutboxID, cause.Error(), opt.MaxAttempts)
	if err != nil {
		// The attempt could not even be recorded; count it as failed so the run
		// reports honestly and moves on.
		stats.Failed++
		return
	}
	if dead {
		stats.DeadLettered++
		return
	}
	stats.Failed++
}
