package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"

	"egot-tracker/internal/config"
	"egot-tracker/internal/database"
	"egot-tracker/internal/models"
	"egot-tracker/internal/repository"
	"egot-tracker/internal/scraper"
	"egot-tracker/internal/service"
)

func main() {
	// Parse flags
	year := flag.Int("year", 2025, "Oscar ceremony year")
	reset := flag.Bool("reset", false, "Delete existing ceremony data for this year before creating")
	flag.Parse()

	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories and services
	oscarRepo := repository.NewOscarRepository(pool)
	celebrityRepo := repository.NewCelebrityRepository(pool)
	oscarService := service.NewOscarService(oscarRepo, celebrityRepo)
	wikiScraper := scraper.NewWikipediaScraper()

	log.Printf("Setting up Oscar race for %d...\n", *year)

	// Check if ceremony already exists
	existing, _ := oscarRepo.GetCeremonyByYear(ctx, *year)
	if existing != nil {
		if *reset {
			log.Printf("Deleting existing ceremony data for %d...", *year)
			if err := oscarService.DeleteCeremony(ctx, *year); err != nil {
				log.Fatalf("Failed to delete existing ceremony: %v", err)
			}
		} else {
			log.Fatalf("Ceremony for %d already exists. Use --reset to recreate.", *year)
		}
	}

	// Get nominations data
	var nominations []scraper.OscarNomination
	if *year == 2025 {
		nominations = scraper.GetOscarNominations2025()
	} else {
		log.Fatalf("No nomination data available for year %d. Only 2025 is currently supported.", *year)
	}

	// Create ceremony
	ceremonyName := scraper.GetCeremonyName(*year)
	ceremony, err := oscarService.CreateCeremony(ctx, &models.OscarCeremony{
		Year: *year,
		CeremonyName: pgtype.Text{
			String: ceremonyName,
			Valid:  true,
		},
		IsComplete: false,
	})
	if err != nil {
		log.Fatalf("Failed to create ceremony: %v", err)
	}
	log.Printf("Created ceremony: %s\n", ceremonyName)

	// Track stats
	categoriesCreated := 0
	nomineesCreated := 0
	celebritiesCreated := 0

	// Create categories and nominees
	for i, nom := range nominations {
		// Create category
		category, err := oscarService.CreateCategory(ctx, &models.OscarCategory{
			CeremonyID:      ceremony.ID,
			Name:            nom.Category,
			DisplayOrder:    i,
			WinnerAnnounced: false,
		})
		if err != nil {
			log.Printf("Failed to create category %s: %v", nom.Category, err)
			continue
		}
		categoriesCreated++
		log.Printf("  Category: %s (%d nominees)", nom.Category, len(nom.Nominees))

		// Create nominees
		for j, nomineeInfo := range nom.Nominees {
			nominee := &models.OscarNominee{
				CategoryID:   category.ID,
				Name:         nomineeInfo.Name,
				DisplayOrder: j,
				IsWinner:     false,
			}

			// Set work title
			if nomineeInfo.WorkTitle != "" {
				nominee.WorkTitle = pgtype.Text{
					String: nomineeInfo.WorkTitle,
					Valid:  true,
				}
			}

			// If this is a person, try to find/create celebrity and fetch their photo
			if nomineeInfo.IsPerson {
				// Try to fetch from Wikipedia
				summary, err := wikiScraper.FetchPersonSummary(ctx, nomineeInfo.Name)
				if err != nil {
					log.Printf("    Warning: Could not fetch Wikipedia data for %s: %v", nomineeInfo.Name, err)
				}

				photoURL := ""
				bio := ""
				if summary != nil {
					if summary.Thumbnail != nil {
						photoURL = summary.Thumbnail.Source
					}
					bio = summary.Extract
				}

				// Find or create celebrity
				celebrity, err := oscarService.FindOrCreateCelebrity(ctx, nomineeInfo.Name, photoURL, bio)
				if err != nil {
					log.Printf("    Warning: Could not create celebrity for %s: %v", nomineeInfo.Name, err)
				} else {
					nominee.CelebrityID = celebrity.ID
					if celebrity.PhotoURL.Valid {
						nominee.PhotoURL = celebrity.PhotoURL
					}
					celebritiesCreated++
				}

				// Small delay to be nice to Wikipedia API
				time.Sleep(200 * time.Millisecond)
			} else {
				// This is a film/work - try to fetch movie poster from Wikipedia
				summary, err := wikiScraper.FetchFilmSummary(ctx, nomineeInfo.Name)
				if err != nil {
					log.Printf("    Warning: Could not fetch Wikipedia data for film %s: %v", nomineeInfo.Name, err)
				} else if summary != nil && summary.Thumbnail != nil {
					nominee.PhotoURL = pgtype.Text{
						String: summary.Thumbnail.Source,
						Valid:  true,
					}
				}

				// Small delay to be nice to Wikipedia API
				time.Sleep(200 * time.Millisecond)
			}

			// Create nominee
			_, err := oscarService.CreateNominee(ctx, nominee)
			if err != nil {
				log.Printf("    Failed to create nominee %s: %v", nomineeInfo.Name, err)
				continue
			}
			nomineesCreated++
		}
	}

	log.Println("\n=== Setup Complete ===")
	log.Printf("Ceremony: %s", ceremonyName)
	log.Printf("Categories: %d", categoriesCreated)
	log.Printf("Nominees: %d", nomineesCreated)
	log.Printf("Celebrities created/linked: %d", celebritiesCreated)
	log.Printf("\nView at: http://localhost:3000/oscar-race/%d", *year)
}
