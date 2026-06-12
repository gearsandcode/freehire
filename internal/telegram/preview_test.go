package telegram

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestParsePreview(t *testing.T) {
	html, err := os.ReadFile("testdata/preview.html")
	if err != nil {
		t.Fatal(err)
	}

	posts, err := ParsePreview("hrlunapark", html)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(posts) != 3 {
		t.Fatalf("posts = %d, want 3", len(posts))
	}

	p := posts[1] // hrlunapark/392 — the Claimsorted multi-vacancy post
	if p.MsgID != 392 {
		t.Errorf("msg_id = %d, want 392", p.MsgID)
	}
	want := time.Date(2026, 5, 28, 12, 3, 7, 0, time.UTC)
	if !p.PostedAt.Equal(want) {
		t.Errorf("posted_at = %v, want %v", p.PostedAt, want)
	}
	if !strings.Contains(p.Text, "ML & full-stack engineers") {
		t.Errorf("text lacks the tl;dr line; got %.120q", p.Text)
	}
	// Entities decoded, tags stripped, <br> became newlines.
	if strings.Contains(p.Text, "<") || strings.Contains(p.Text, "&amp;") {
		t.Errorf("text still contains markup: %.200q", p.Text)
	}
	if !strings.Contains(p.Text, "\n") {
		t.Error("text has no newlines — <br> separation lost")
	}
}

func TestParsePreviewYieldsPostsInPageOrder(t *testing.T) {
	html, _ := os.ReadFile("testdata/preview.html")
	posts, err := ParsePreview("hrlunapark", html)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(posts); i++ {
		if posts[i].MsgID <= posts[i-1].MsgID {
			t.Errorf("posts out of order: %d after %d", posts[i].MsgID, posts[i-1].MsgID)
		}
	}
}

func TestParsePreviewZeroPostsIsAnError(t *testing.T) {
	// A page with no recognizable posts means the markup drifted or the channel
	// has no preview — the caller must count the channel failed, loudly.
	_, err := ParsePreview("foo", []byte("<html><body>nothing here</body></html>"))
	if err == nil {
		t.Fatal("want error for a page with zero parseable posts")
	}
}
