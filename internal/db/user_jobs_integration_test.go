//go:build integration

// Integration tests for the user_jobs upsert queries — view recording and apply
// marking are ON CONFLICT semantics that can only be verified against a real
// Postgres. Run with: go test -tags=integration ./internal/db/
// Requires Docker (testcontainers spins up a throwaway Postgres with the migrations).
package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func insertUser(t *testing.T, pool *pgxpool.Pool, email string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO users (email) VALUES ($1) RETURNING id`, email).Scan(&id); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return id
}

func TestUserJobs(t *testing.T) {
	pool := startPostgres(t)
	q := New(pool)
	ctx := context.Background()

	countRows := func(t *testing.T) int {
		t.Helper()
		var n int
		if err := pool.QueryRow(ctx, "SELECT count(*) FROM user_jobs").Scan(&n); err != nil {
			t.Fatalf("count: %v", err)
		}
		return n
	}

	t.Run("RecordJobView creates then refreshes one row", func(t *testing.T) {
		if _, err := pool.Exec(ctx, "TRUNCATE user_jobs, users, jobs RESTART IDENTITY CASCADE"); err != nil {
			t.Fatalf("truncate: %v", err)
		}
		uid := insertUser(t, pool, "viewer@example.test")
		jid := insertJob(t, pool, "view-job")

		first, err := q.RecordJobView(ctx, RecordJobViewParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("first view: %v", err)
		}
		if !first.ViewedAt.Valid || first.AppliedAt.Valid {
			t.Errorf("first view: viewed=%v applied=%v, want viewed/not-applied", first.ViewedAt.Valid, first.AppliedAt.Valid)
		}

		second, err := q.RecordJobView(ctx, RecordJobViewParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("second view: %v", err)
		}
		if n := countRows(t); n != 1 {
			t.Errorf("rows = %d, want 1 (one per (user, job))", n)
		}
		if second.ViewedAt.Time.Before(first.ViewedAt.Time) {
			t.Errorf("second viewed_at %v is before first %v — not refreshed", second.ViewedAt.Time, first.ViewedAt.Time)
		}
		if second.AppliedAt.Valid {
			t.Error("a view must not set applied_at")
		}
	})

	t.Run("MarkJobApplied sets applied_at without a prior view and is idempotent", func(t *testing.T) {
		if _, err := pool.Exec(ctx, "TRUNCATE user_jobs, users, jobs RESTART IDENTITY CASCADE"); err != nil {
			t.Fatalf("truncate: %v", err)
		}
		uid := insertUser(t, pool, "applier@example.test")
		jid := insertJob(t, pool, "apply-job")

		// No prior RecordJobView: the insert path must still populate viewed_at.
		first, err := q.MarkJobApplied(ctx, MarkJobAppliedParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("first apply: %v", err)
		}
		if !first.AppliedAt.Valid || !first.ViewedAt.Valid {
			t.Errorf("first apply: applied=%v viewed=%v, want both set", first.AppliedAt.Valid, first.ViewedAt.Valid)
		}

		second, err := q.MarkJobApplied(ctx, MarkJobAppliedParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("second apply: %v", err)
		}
		if n := countRows(t); n != 1 {
			t.Errorf("rows = %d, want 1 (idempotent)", n)
		}
		if second.AppliedAt.Time.Before(first.AppliedAt.Time) {
			t.Errorf("second applied_at %v is before first %v", second.AppliedAt.Time, first.AppliedAt.Time)
		}
	})

	t.Run("apply after view keeps the single row and preserves viewed_at", func(t *testing.T) {
		if _, err := pool.Exec(ctx, "TRUNCATE user_jobs, users, jobs RESTART IDENTITY CASCADE"); err != nil {
			t.Fatalf("truncate: %v", err)
		}
		uid := insertUser(t, pool, "both@example.test")
		jid := insertJob(t, pool, "both-job")

		viewed, err := q.RecordJobView(ctx, RecordJobViewParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("view: %v", err)
		}
		applied, err := q.MarkJobApplied(ctx, MarkJobAppliedParams{UserID: uid, JobID: jid})
		if err != nil {
			t.Fatalf("apply: %v", err)
		}
		if n := countRows(t); n != 1 {
			t.Errorf("rows = %d, want 1 (view then apply is the same row)", n)
		}
		if !applied.AppliedAt.Valid {
			t.Error("apply after view must set applied_at")
		}
		if !applied.ViewedAt.Time.Equal(viewed.ViewedAt.Time) {
			t.Errorf("apply changed viewed_at: %v -> %v", viewed.ViewedAt.Time, applied.ViewedAt.Time)
		}
	})
}
