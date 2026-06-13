package linksource

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/strelov1/freehire/internal/sources"
)

// remoteYeah resolves RemoteYeah vacancies. Posts link directly to
// remoteyeah.com/jobs/<slug>; the slug is the stable id.
type remoteYeah struct {
	http Client
}

// NewRemoteYeah builds the RemoteYeah link-source adapter.
func NewRemoteYeah(c Client) LinkSource { return remoteYeah{http: c} }

func (remoteYeah) Source() string { return "remoteyeah" }

// Match handles remoteyeah.com/jobs/<slug> links only — the homepage and other paths are
// not vacancies.
func (remoteYeah) Match(u *url.URL) bool {
	return host(u) == "remoteyeah.com" && strings.HasPrefix(u.Path, "/jobs/") &&
		strings.TrimPrefix(u.Path, "/jobs/") != ""
}

// remoteYeahPosting selects the JobPosting ld+json fields RemoteYeah publishes.
type remoteYeahPosting struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	DatePosted         string `json:"datePosted"`
	JobLocationType    string `json:"jobLocationType"`
	HiringOrganization struct {
		Name string `json:"name"`
	} `json:"hiringOrganization"`
	BaseSalary struct {
		Currency string `json:"currency"`
		Value    struct {
			MinValue float64 `json:"minValue"`
			MaxValue float64 `json:"maxValue"`
		} `json:"value"`
	} `json:"baseSalary"`
}

// Resolve fetches the job page and parses its JobPosting ld+json. The slug from the link
// path is the id; the page carries no native identifier.
func (r remoteYeah) Resolve(ctx context.Context, raw string) (sources.Job, bool, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return sources.Job{}, false, err
	}
	slug := strings.Trim(strings.TrimPrefix(u.Path, "/jobs/"), "/")
	if slug == "" {
		return sources.Job{}, false, nil
	}

	node, err := r.http.GetHTML(ctx, raw)
	if err != nil {
		return sources.Job{}, false, err
	}
	var p remoteYeahPosting
	if !sources.LDJobPosting(node, &p) {
		return sources.Job{}, false, fmt.Errorf("linksource: remoteyeah job %s has no JobPosting ld+json", slug)
	}

	desc := sources.SanitizeHTML(p.Description)
	if salary := remoteYeahSalary(p); salary != "" {
		// Sanitize the salary fragment too: its currency is third-party JSON-LD text and
		// the description is rendered with {@html}, so an unsanitized prefix is stored XSS.
		desc = sources.SanitizeHTML(salary) + desc
	}
	return sources.Job{
		ExternalID:  slug,
		URL:         "https://remoteyeah.com/jobs/" + slug,
		Title:       p.Title,
		Company:     p.HiringOrganization.Name,
		Description: desc,
		Remote:      strings.EqualFold(p.JobLocationType, "TELECOMMUTE"),
		PostedAt:    parseRFC3339(p.DatePosted),
	}, true, nil
}

// remoteYeahSalary renders a structured baseSalary range as a leading paragraph, or "" when
// the page states no amount. Folding it into the description keeps it visible and lets
// enrichment pick it up (sources.Job has no dedicated salary field).
func remoteYeahSalary(p remoteYeahPosting) string {
	min, max := p.BaseSalary.Value.MinValue, p.BaseSalary.Value.MaxValue
	if min <= 0 && max <= 0 {
		return ""
	}
	cur := p.BaseSalary.Currency
	switch {
	case min > 0 && max > 0:
		return fmt.Sprintf("<p>Salary: %.0f–%.0f %s</p>", min, max, cur)
	case min > 0:
		return fmt.Sprintf("<p>Salary: from %.0f %s</p>", min, cur)
	default:
		return fmt.Sprintf("<p>Salary: up to %.0f %s</p>", max, cur)
	}
}
