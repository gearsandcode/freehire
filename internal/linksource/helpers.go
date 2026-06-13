package linksource

import "time"

// parseDate parses a date-only timestamp ("2006-01-02", as Habr's datePosted emits),
// returning nil for an empty or unparseable value (posted_at is nullable).
func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

// parseRFC3339 parses an RFC3339 timestamp (as RemoteYeah's datePosted emits), in UTC,
// returning nil for an empty or unparseable value.
func parseRFC3339(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	t = t.UTC()
	return &t
}
