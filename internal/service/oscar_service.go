package service

import (
	"context"
	"errors"

	"egot-tracker/internal/models"
	"egot-tracker/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

var ErrCeremonyNotFound = errors.New("ceremony not found")

type OscarService struct {
	oscarRepo     *repository.OscarRepository
	celebrityRepo *repository.CelebrityRepository
}

func NewOscarService(oscarRepo *repository.OscarRepository, celebrityRepo *repository.CelebrityRepository) *OscarService {
	return &OscarService{
		oscarRepo:     oscarRepo,
		celebrityRepo: celebrityRepo,
	}
}

// GetCeremony returns a full ceremony with all categories and nominees
func (s *OscarService) GetCeremony(ctx context.Context, year int) (*models.OscarCeremonyFull, error) {
	ceremony, err := s.oscarRepo.GetFullCeremony(ctx, year)
	if errors.Is(err, repository.ErrCeremonyNotFound) {
		return nil, ErrCeremonyNotFound
	}
	return ceremony, err
}

// GetAllYears returns all tracked Oscar years
func (s *OscarService) GetAllYears(ctx context.Context) ([]int, error) {
	return s.oscarRepo.GetAllCeremonyYears(ctx)
}

// SetWinner marks a nominee as the winner for their category
func (s *OscarService) SetWinner(ctx context.Context, nomineeID pgtype.UUID) error {
	return s.oscarRepo.SetNomineeAsWinner(ctx, nomineeID)
}

// CreateCeremony creates a new Oscar ceremony
func (s *OscarService) CreateCeremony(ctx context.Context, ceremony *models.OscarCeremony) (*models.OscarCeremony, error) {
	return s.oscarRepo.CreateCeremony(ctx, ceremony)
}

// CreateCategory creates a new category for a ceremony
func (s *OscarService) CreateCategory(ctx context.Context, category *models.OscarCategory) (*models.OscarCategory, error) {
	return s.oscarRepo.CreateCategory(ctx, category)
}

// CreateNominee creates a new nominee for a category
func (s *OscarService) CreateNominee(ctx context.Context, nominee *models.OscarNominee) (*models.OscarNominee, error) {
	return s.oscarRepo.CreateNominee(ctx, nominee)
}

// DeleteCeremony removes a ceremony and all related data
func (s *OscarService) DeleteCeremony(ctx context.Context, year int) error {
	return s.oscarRepo.DeleteCeremony(ctx, year)
}

// FindOrCreateCelebrity finds a celebrity by name or creates them if they don't exist
func (s *OscarService) FindOrCreateCelebrity(ctx context.Context, name, photoURL, summary string) (*models.Celebrity, error) {
	// Try to find existing celebrity
	celebrity, err := s.celebrityRepo.FindByName(ctx, name)
	if err == nil {
		return celebrity, nil
	}

	// Create new celebrity if not found
	if errors.Is(err, repository.ErrCelebrityNotFound) {
		newCelebrity := &models.Celebrity{
			Name: name,
			Slug: slugify(name),
		}
		if photoURL != "" {
			newCelebrity.PhotoURL.String = photoURL
			newCelebrity.PhotoURL.Valid = true
		}
		if summary != "" {
			newCelebrity.Summary.String = summary
			newCelebrity.Summary.Valid = true
		}
		return s.celebrityRepo.Create(ctx, newCelebrity)
	}

	return nil, err
}

// slugify converts a name to a URL-friendly slug
func slugify(name string) string {
	slug := name
	// Simple slugification - replace spaces with dashes, lowercase
	result := ""
	for _, c := range slug {
		if c >= 'a' && c <= 'z' {
			result += string(c)
		} else if c >= 'A' && c <= 'Z' {
			result += string(c + 32) // lowercase
		} else if c == ' ' || c == '-' {
			result += "-"
		}
	}
	return result
}
