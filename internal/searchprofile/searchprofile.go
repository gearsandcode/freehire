// Package searchprofile is the per-user search-profile use case: a signed-in user names a
// record of their professional self — one specialization (a job category) and a non-empty
// set of skills — and can list, create, overwrite, and delete those profiles. It owns
// validation (name bounds, the specialization vocabulary, skill normalization, the
// per-user cap); the Repository owns persistence and maps the relevant Postgres conditions
// (unique violation, no row) onto the package sentinels. How a profile is consumed (match
// scoring, ranked feeds, notifications) lives outside this package.
package searchprofile

import (
	"context"
	"errors"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/strelov1/freehire/internal/db"
	"github.com/strelov1/freehire/internal/enrich"
)

// Sentinel errors mapped to HTTP statuses by the handler.
var (
	// ErrInvalidName is a blank or over-long name (mapped to 400).
	ErrInvalidName = errors.New("searchprofile: name must be 1-100 characters")
	// ErrInvalidSpecialization is a specialization outside the category vocabulary
	// (mapped to 400).
	ErrInvalidSpecialization = errors.New("searchprofile: specialization is not a known category")
	// ErrEmptySkills is a profile whose skills are empty after normalization (mapped to 400).
	ErrEmptySkills = errors.New("searchprofile: at least one skill is required")
	// ErrDuplicateName is a name the user already uses (the UNIQUE (user_id, name)
	// constraint; mapped to 409).
	ErrDuplicateName = errors.New("searchprofile: a profile with this name already exists")
	// ErrCapExceeded is a create past the per-user limit (mapped to 409).
	ErrCapExceeded = errors.New("searchprofile: profile limit reached")
	// ErrNotFound is a missing or non-owned target (mapped to 404).
	ErrNotFound = errors.New("searchprofile: not found")
)

const (
	// maxNameLen bounds a profile name; the migration's CHECK is the backstop.
	maxNameLen = 100
	// maxPerUser caps how many profiles a single user may keep.
	maxPerUser = 50
)

// Repository is the persistence contract for search profiles. Every method is
// user-scoped. Create maps a unique violation to ErrDuplicateName; Update maps a unique
// violation to ErrDuplicateName and a missing owner-scoped row to ErrNotFound; Delete
// maps "no row affected" to ErrNotFound.
type Repository interface {
	List(ctx context.Context, userID int64) ([]db.SearchProfile, error)
	Count(ctx context.Context, userID int64) (int64, error)
	Create(ctx context.Context, p db.CreateSearchProfileParams) (db.SearchProfile, error)
	Update(ctx context.Context, p db.UpdateSearchProfileParams) (db.SearchProfile, error)
	Delete(ctx context.Context, p db.DeleteSearchProfileParams) error
}

// Service implements the search-profile use cases.
type Service struct {
	repo Repository
}

// New creates a Service backed by the given Repository.
func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// List returns the user's profiles, most recently updated first.
func (s *Service) List(ctx context.Context, userID int64) ([]db.SearchProfile, error) {
	return s.repo.List(ctx, userID)
}

// Create validates and stores a profile for the user. The name is trimmed and bounded; the
// specialization must be a known category; the skills are normalized and must be non-empty;
// the per-user cap is checked before the insert; a duplicate name surfaces as
// ErrDuplicateName (mapped by the repository).
func (s *Service) Create(ctx context.Context, userID int64, name, specialization string, skills []string) (db.SearchProfile, error) {
	name, err := validName(name)
	if err != nil {
		return db.SearchProfile{}, err
	}
	specialization, err = validSpecialization(specialization)
	if err != nil {
		return db.SearchProfile{}, err
	}
	normalized, err := normalizeSkills(skills)
	if err != nil {
		return db.SearchProfile{}, err
	}
	count, err := s.repo.Count(ctx, userID)
	if err != nil {
		return db.SearchProfile{}, err
	}
	if count >= maxPerUser {
		return db.SearchProfile{}, ErrCapExceeded
	}
	return s.repo.Create(ctx, db.CreateSearchProfileParams{
		UserID:         userID,
		Name:           name,
		Specialization: specialization,
		Skills:         normalized,
	})
}

// Update overwrites a profile's name, specialization, and/or skills, scoped to its owner.
// A nil name/specialization pointer or a nil skills slice leaves that column unchanged
// (partial update). A provided name is validated; a provided specialization must be a
// known category; provided skills are normalized and must be non-empty. A missing or
// non-owned row surfaces as ErrNotFound (mapped by the repository).
func (s *Service) Update(ctx context.Context, userID, id int64, name, specialization *string, skills []string) (db.SearchProfile, error) {
	p := db.UpdateSearchProfileParams{ID: id, UserID: userID}
	if name != nil {
		valid, err := validName(*name)
		if err != nil {
			return db.SearchProfile{}, err
		}
		p.Name = pgtype.Text{String: valid, Valid: true}
	}
	if specialization != nil {
		valid, err := validSpecialization(*specialization)
		if err != nil {
			return db.SearchProfile{}, err
		}
		p.Specialization = pgtype.Text{String: valid, Valid: true}
	}
	if skills != nil {
		normalized, err := normalizeSkills(skills)
		if err != nil {
			return db.SearchProfile{}, err
		}
		p.Skills = normalized
	}
	return s.repo.Update(ctx, p)
}

// Delete removes one of the user's profiles. A missing or non-owned row surfaces as
// ErrNotFound (mapped by the repository).
func (s *Service) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, db.DeleteSearchProfileParams{ID: id, UserID: userID})
}

// validName trims the name and enforces the 1..maxNameLen bound (counted in runes, to
// match the DB CHECK's character semantics — names are often Cyrillic), returning the
// trimmed value or ErrInvalidName.
func validName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > maxNameLen {
		return "", ErrInvalidName
	}
	return name, nil
}

// validSpecialization trims the value and checks membership in the controlled category
// vocabulary (the same enum the rest of the app validates against), returning the trimmed
// value or ErrInvalidSpecialization.
func validSpecialization(specialization string) (string, error) {
	specialization = strings.TrimSpace(specialization)
	if !slices.Contains(enrich.CategoryValues, specialization) {
		return "", ErrInvalidSpecialization
	}
	return specialization, nil
}

// normalizeSkills lowercases, trims, and deduplicates the skills (preserving first-seen
// order), dropping blanks. It returns ErrEmptySkills if nothing remains — a profile without
// skills has no meaning.
func normalizeSkills(skills []string) ([]string, error) {
	out := make([]string, 0, len(skills))
	seen := make(map[string]struct{}, len(skills))
	for _, raw := range skills {
		skill := strings.ToLower(strings.TrimSpace(raw))
		if skill == "" {
			continue
		}
		if _, dup := seen[skill]; dup {
			continue
		}
		seen[skill] = struct{}{}
		out = append(out, skill)
	}
	if len(out) == 0 {
		return nil, ErrEmptySkills
	}
	return out, nil
}
