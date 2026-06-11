package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/strelov1/freehire/internal/auth"
)

// userJobsApp mounts the view/apply routes behind RequireAuth on a handler with
// no DB. The auth-gate and id-parse cases below all reject before any query
// runs, so the nil queries is never dereferenced. The DB-backed happy path and
// idempotency are covered by the db-package integration test (TestUserJobs).
func userJobsApp() (*fiber.App, *auth.Issuer) {
	iss := auth.NewIssuer("test-secret", time.Hour)
	h := &Handler{issuer: iss}
	app := fiber.New()
	app.Post("/jobs/:id/view", auth.RequireAuth(iss), h.RecordView)
	app.Post("/jobs/:id/apply", auth.RequireAuth(iss), h.MarkApplied)
	return app, iss
}

func postUserJob(t *testing.T, app *fiber.App, path, token string) int {
	t.Helper()
	req := httptest.NewRequest(fiber.MethodPost, path, nil)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	return resp.StatusCode
}

func TestRecordView_RequiresAuth(t *testing.T) {
	app, _ := userJobsApp()
	if got := postUserJob(t, app, "/jobs/1/view", ""); got != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", got)
	}
}

func TestMarkApplied_RequiresAuth(t *testing.T) {
	app, _ := userJobsApp()
	if got := postUserJob(t, app, "/jobs/1/apply", ""); got != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", got)
	}
}

func TestRecordView_RejectsNonNumericID(t *testing.T) {
	app, iss := userJobsApp()
	token, err := iss.Issue(1)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	// Authenticated, so RequireAuth passes; the handler rejects the bad id before
	// any query runs.
	if got := postUserJob(t, app, "/jobs/not-a-number/view", token); got != fiber.StatusBadRequest {
		t.Errorf("status = %d, want 400", got)
	}
}

func TestMarkApplied_RejectsNonNumericID(t *testing.T) {
	app, iss := userJobsApp()
	token, err := iss.Issue(1)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if got := postUserJob(t, app, "/jobs/not-a-number/apply", token); got != fiber.StatusBadRequest {
		t.Errorf("status = %d, want 400", got)
	}
}

// interactionResponse is the only interaction shape that reaches a response. This
// locks the contract: it omits user_id and carries job_id + the two timestamps.
func TestInteractionResponse_Shape(t *testing.T) {
	raw, err := json.Marshal(interactionResponse{JobID: 7})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, leaked := fields["user_id"]; leaked {
		t.Error("interactionResponse must not include user_id")
	}
	for _, want := range []string{"job_id", "viewed_at", "applied_at"} {
		if _, ok := fields[want]; !ok {
			t.Errorf("interactionResponse missing %q", want)
		}
	}
}
