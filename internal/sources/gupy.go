package sources

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// gupyBaseURL is the Gupy public portal jobs API. Gupy is the dominant Brazilian
// ATS (Creditas, Afya, Cogna, Omie, …); a company's per-tenant career page lives at
// <subdomain>.gupy.io, but the listing API is this one central host keyed by the
// numeric companyId.
const gupyBaseURL = "https://employability-portal.gupy.io/api/v1/jobs"

// gupyPageLimit is the listing page size. The API caps limit at 100 (limit=200 is a
// 400), so this is also the maximum.
const gupyPageLimit = 100

// gupyMaxPages bounds the offset walk. The stop signal is a short/empty page, but a
// misbehaving API that always returned a full page would otherwise loop forever, so
// this caps the walk at 5000 postings — well above any single company's openings.
const gupyMaxPages = 50

// gupy adapts the Gupy portal API. Its list endpoint carries the description inline,
// so — like Greenhouse — it needs no per-posting detail request. The board id is the
// company's numeric Gupy companyId.
type gupy struct {
	http JSONGetter
}

// NewGupy builds the Gupy adapter over the given HTTP client.
func NewGupy(c JSONGetter) Source { return gupy{http: c} }

func (gupy) Provider() string { return "gupy" }

// gupyJob is one posting from the Gupy listing. The description is inline; jobUrl is an
// absolute public URL; workplaceType is a structured "remote"/"hybrid"/"on-site" enum.
type gupyJob struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	JobURL        string `json:"jobUrl"`
	Description   string `json:"description"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	IsRemoteWork  bool   `json:"isRemoteWork"`
	WorkplaceType string `json:"workplaceType"`
	PublishedDate string `json:"publishedDate"`
}

func (g gupy) Fetch(ctx context.Context, e CompanyEntry) ([]Job, error) {
	postings, err := g.list(ctx, e.Board)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0, len(postings))
	for _, p := range postings {
		if p.JobURL == "" { // url is the dedup key — a posting without one is unusable
			continue
		}
		jobs = append(jobs, Job{
			ExternalID:  strconv.FormatInt(p.ID, 10),
			URL:         p.JobURL,
			Title:       strings.TrimSpace(p.Name),
			Company:     e.Company,
			Location:    joinNonEmpty(p.City, p.State, p.Country),
			Description: sanitizeHTML(p.Description),
			Remote:      p.IsRemoteWork,
			WorkMode:    workplaceTypeMode(p.WorkplaceType),
			PostedAt:    parseRFC3339(p.PublishedDate),
		})
	}
	return jobs, nil
}

// list pages through the company's postings. It stops on an empty or short page rather
// than on pagination.total: when limit == page size, Gupy reports total = min(real, limit),
// so a full first page would falsely look complete. A short page is the reliable last-page
// signal (the same reasoning SmartRecruiters' listPostings relies on totalFound for, but
// Gupy's total can't be trusted).
func (g gupy) list(ctx context.Context, board string) ([]gupyJob, error) {
	var postings []gupyJob
	for offset, page := 0, 0; page < gupyMaxPages; page++ {
		url := fmt.Sprintf("%s?companyId=%s&limit=%d&offset=%d", gupyBaseURL, board, gupyPageLimit, offset)
		var resp struct {
			Data []gupyJob `json:"data"`
		}
		if err := g.http.GetJSON(ctx, url, &resp); err != nil {
			return nil, fmt.Errorf("gupy: list company %s: %w", board, err)
		}
		if len(resp.Data) == 0 {
			break
		}
		postings = append(postings, resp.Data...)
		if len(resp.Data) < gupyPageLimit {
			break // short page = last page
		}
		offset += len(resp.Data)
	}
	return postings, nil
}
