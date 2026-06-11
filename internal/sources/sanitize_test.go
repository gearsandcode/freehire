package sources

import (
	"strings"
	"testing"
)

func TestSanitizeHTML(t *testing.T) {
	in := `<h2>Role</h2><p>Lead the <strong>backend</strong> team.</p>` +
		`<ul><li>Ship features</li></ul>` +
		`<a href="https://example.com" onclick="steal()">apply</a>` +
		`<img src="https://evil.example/track.gif">` +
		`<script>alert(1)</script>`

	got := sanitizeHTML(in)

	// Structural formatting is preserved.
	for _, want := range []string{"<h2>Role</h2>", "<strong>backend</strong>", "<li>Ship features</li>", `href="https://example.com"`} {
		if !strings.Contains(got, want) {
			t.Errorf("sanitizeHTML dropped expected markup %q\ngot: %s", want, got)
		}
	}

	// Active content and external request vectors are stripped.
	for _, bad := range []string{"<script", "onclick", "alert(1)", "<img", "track.gif"} {
		if strings.Contains(got, bad) {
			t.Errorf("sanitizeHTML kept unsafe content %q\ngot: %s", bad, got)
		}
	}

	// Links are defanged so untrusted postings cannot pass link authority.
	if !strings.Contains(got, `rel="nofollow"`) {
		t.Errorf("sanitizeHTML should mark links nofollow\ngot: %s", got)
	}
}
