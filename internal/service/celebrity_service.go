package service

import (
	"context"
	"errors"
	"log"

	"egot-tracker/internal/models"
	"egot-tracker/internal/repository"
	"egot-tracker/internal/scraper"
)

var ErrCelebrityNotFound = errors.New("celebrity not found")

type CelebrityService struct {
	celebrityRepo *repository.CelebrityRepository
	awardRepo     *repository.AwardRepository
	scraper       *scraper.WikidataScraper
}

func NewCelebrityService(
	celebrityRepo *repository.CelebrityRepository,
	awardRepo *repository.AwardRepository,
	scraper *scraper.WikidataScraper,
) *CelebrityService {
	return &CelebrityService{
		celebrityRepo: celebrityRepo,
		awardRepo:     awardRepo,
		scraper:       scraper,
	}
}

// SearchCelebrity implements the cache-aside pattern:
// 1. Check if celebrity exists in database
// 2. If found, fetch awards and return
// 3. If not found, scrape from Wikidata, save to DB, and return
func (s *CelebrityService) SearchCelebrity(ctx context.Context, name string) (*models.CelebrityWithAwards, error) {
	// Step 1: Check database for celebrity
	celebrity, err := s.celebrityRepo.FindByName(ctx, name)
	if err != nil && !errors.Is(err, repository.ErrCelebrityNotFound) {
		return nil, err
	}

	// Step 2: If found in DB, fetch awards and return
	if celebrity != nil {
		awards, err := s.awardRepo.FindByCelebrityID(ctx, celebrity.ID)
		if err != nil {
			return nil, err
		}
		if awards == nil {
			awards = []models.Award{}
		}
		return &models.CelebrityWithAwards{
			Celebrity: *celebrity,
			Awards:    awards,
		}, nil
	}

	// Step 3: Not in DB - scrape from Wikidata
	log.Printf("Celebrity not in DB, fetching from Wikidata: %s", name)

	scrapedCelebrity, scrapedAwards, err := s.scraper.FetchCelebrity(ctx, name)
	if err != nil {
		log.Printf("Failed to fetch from Wikidata: %v", err)
		return nil, ErrCelebrityNotFound
	}

	// Step 4: Save celebrity to database
	savedCelebrity, err := s.celebrityRepo.Create(ctx, scrapedCelebrity)
	if err != nil {
		log.Printf("Failed to save celebrity: %v", err)
		return nil, err
	}

	// Step 5: Save awards to database
	savedAwards, err := s.awardRepo.CreateBatch(ctx, savedCelebrity.ID, scrapedAwards)
	if err != nil {
		log.Printf("Failed to save awards: %v", err)
		return nil, err
	}

	log.Printf("Saved %s with %d awards from Wikidata", savedCelebrity.Name, len(savedAwards))

	return &models.CelebrityWithAwards{
		Celebrity: *savedCelebrity,
		Awards:    savedAwards,
	}, nil
}

// Autocomplete returns celebrities matching the query from the local database
func (s *CelebrityService) Autocomplete(ctx context.Context, query string, limit int) ([]models.Celebrity, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.celebrityRepo.Search(ctx, query, limit)
}

// GetCloseToEGOT returns celebrities with 3 of 4 EGOT awards
func (s *CelebrityService) GetCloseToEGOT(ctx context.Context, limit int) ([]models.CelebrityWithEGOTProgress, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.celebrityRepo.FindCloseToEGOT(ctx, limit)
}

// GetEGOTWinners returns celebrities with all 4 EGOT awards
func (s *CelebrityService) GetEGOTWinners(ctx context.Context, limit int) ([]models.CelebrityWithEGOTProgress, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.celebrityRepo.FindEGOTWinners(ctx, limit)
}

// GetNoAwards returns celebrities with no awards
func (s *CelebrityService) GetNoAwards(ctx context.Context, limit int) ([]models.Celebrity, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.celebrityRepo.FindNoAwards(ctx, limit)
}
