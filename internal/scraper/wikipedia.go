package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WikipediaScraper fetches data from Wikipedia REST API
type WikipediaScraper struct {
	httpClient *http.Client
}

// NewWikipediaScraper creates a new Wikipedia scraper
func NewWikipediaScraper() *WikipediaScraper {
	return &WikipediaScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WikipediaPageSummary represents the summary response from Wikipedia
type WikipediaPageSummary struct {
	Title       string `json:"title"`
	Extract     string `json:"extract"`
	Description string `json:"description"`
	Thumbnail   *struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

// FetchPersonSummary fetches a person's summary and photo from Wikipedia
func (w *WikipediaScraper) FetchPersonSummary(ctx context.Context, name string) (*WikipediaPageSummary, error) {
	return w.FetchPageSummary(ctx, name)
}

// FetchFilmSummary fetches a film's summary and poster from Wikipedia
// It tries the film title first, then with " (film)" suffix, then " (2024 film)" suffix
func (w *WikipediaScraper) FetchFilmSummary(ctx context.Context, title string) (*WikipediaPageSummary, error) {
	// Try exact title first
	summary, err := w.FetchPageSummary(ctx, title)
	if err == nil && summary.Thumbnail != nil {
		return summary, nil
	}

	// Try with "(film)" suffix
	summary, err = w.FetchPageSummary(ctx, title+" (film)")
	if err == nil && summary.Thumbnail != nil {
		return summary, nil
	}

	// Try with "(2024 film)" suffix for recent films
	summary, err = w.FetchPageSummary(ctx, title+" (2024 film)")
	if err == nil && summary.Thumbnail != nil {
		return summary, nil
	}

	// Return whatever we got, even without thumbnail
	if summary != nil {
		return summary, nil
	}

	return nil, fmt.Errorf("could not find Wikipedia page for film: %s", title)
}

// FetchPageSummary fetches a Wikipedia page summary by title
func (w *WikipediaScraper) FetchPageSummary(ctx context.Context, title string) (*WikipediaPageSummary, error) {
	// Wikipedia uses underscores for spaces in titles
	wikiTitle := strings.ReplaceAll(title, " ", "_")
	apiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(wikiTitle))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "EGOT-Tracker/1.0 (https://github.com/egot-tracker)")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wikipedia API returned status %d", resp.StatusCode)
	}

	var summary WikipediaPageSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &summary, nil
}

// OscarNomination represents a parsed Oscar nomination
type OscarNomination struct {
	Category  string
	Nominees  []NomineeInfo
}

// NomineeInfo represents info about a nominee
type NomineeInfo struct {
	Name      string
	WorkTitle string
	IsPerson  bool // true if this is a person (actor, director), false if it's a work (film)
}

// GetOscarNominations2025 returns hardcoded 2025 Oscar nominations
// In a production system, this would parse from Wikipedia
func GetOscarNominations2025() []OscarNomination {
	return []OscarNomination{
		{
			Category: "Best Picture",
			Nominees: []NomineeInfo{
				{Name: "Anora", WorkTitle: "Anora", IsPerson: false},
				{Name: "The Brutalist", WorkTitle: "The Brutalist", IsPerson: false},
				{Name: "A Complete Unknown", WorkTitle: "A Complete Unknown", IsPerson: false},
				{Name: "Conclave", WorkTitle: "Conclave", IsPerson: false},
				{Name: "Dune: Part Two", WorkTitle: "Dune: Part Two", IsPerson: false},
				{Name: "Emilia Pérez", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "I'm Still Here", WorkTitle: "I'm Still Here", IsPerson: false},
				{Name: "Nickel Boys", WorkTitle: "Nickel Boys", IsPerson: false},
				{Name: "The Substance", WorkTitle: "The Substance", IsPerson: false},
				{Name: "Wicked", WorkTitle: "Wicked", IsPerson: false},
			},
		},
		{
			Category: "Best Director",
			Nominees: []NomineeInfo{
				{Name: "Sean Baker", WorkTitle: "Anora", IsPerson: true},
				{Name: "Brady Corbet", WorkTitle: "The Brutalist", IsPerson: true},
				{Name: "James Mangold", WorkTitle: "A Complete Unknown", IsPerson: true},
				{Name: "Jacques Audiard", WorkTitle: "Emilia Pérez", IsPerson: true},
				{Name: "Coralie Fargeat", WorkTitle: "The Substance", IsPerson: true},
			},
		},
		{
			Category: "Best Actor",
			Nominees: []NomineeInfo{
				{Name: "Adrien Brody", WorkTitle: "The Brutalist", IsPerson: true},
				{Name: "Timothée Chalamet", WorkTitle: "A Complete Unknown", IsPerson: true},
				{Name: "Colman Domingo", WorkTitle: "Sing Sing", IsPerson: true},
				{Name: "Ralph Fiennes", WorkTitle: "Conclave", IsPerson: true},
				{Name: "Sebastian Stan", WorkTitle: "The Apprentice", IsPerson: true},
			},
		},
		{
			Category: "Best Actress",
			Nominees: []NomineeInfo{
				{Name: "Cynthia Erivo", WorkTitle: "Wicked", IsPerson: true},
				{Name: "Karla Sofía Gascón", WorkTitle: "Emilia Pérez", IsPerson: true},
				{Name: "Mikey Madison", WorkTitle: "Anora", IsPerson: true},
				{Name: "Demi Moore", WorkTitle: "The Substance", IsPerson: true},
				{Name: "Fernanda Torres", WorkTitle: "I'm Still Here", IsPerson: true},
			},
		},
		{
			Category: "Best Supporting Actor",
			Nominees: []NomineeInfo{
				{Name: "Yura Borisov", WorkTitle: "Anora", IsPerson: true},
				{Name: "Kieran Culkin", WorkTitle: "A Real Pain", IsPerson: true},
				{Name: "Edward Norton", WorkTitle: "A Complete Unknown", IsPerson: true},
				{Name: "Guy Pearce", WorkTitle: "The Brutalist", IsPerson: true},
				{Name: "Jeremy Strong", WorkTitle: "The Apprentice", IsPerson: true},
			},
		},
		{
			Category: "Best Supporting Actress",
			Nominees: []NomineeInfo{
				{Name: "Monica Barbaro", WorkTitle: "A Complete Unknown", IsPerson: true},
				{Name: "Ariana Grande", WorkTitle: "Wicked", IsPerson: true},
				{Name: "Felicity Jones", WorkTitle: "The Brutalist", IsPerson: true},
				{Name: "Isabella Rossellini", WorkTitle: "Conclave", IsPerson: true},
				{Name: "Zoe Saldaña", WorkTitle: "Emilia Pérez", IsPerson: true},
			},
		},
		{
			Category: "Best Original Screenplay",
			Nominees: []NomineeInfo{
				{Name: "Anora", WorkTitle: "Anora", IsPerson: false},
				{Name: "The Brutalist", WorkTitle: "The Brutalist", IsPerson: false},
				{Name: "A Real Pain", WorkTitle: "A Real Pain", IsPerson: false},
				{Name: "September 5", WorkTitle: "September 5", IsPerson: false},
				{Name: "The Substance", WorkTitle: "The Substance", IsPerson: false},
			},
		},
		{
			Category: "Best Adapted Screenplay",
			Nominees: []NomineeInfo{
				{Name: "A Complete Unknown", WorkTitle: "A Complete Unknown", IsPerson: false},
				{Name: "Conclave", WorkTitle: "Conclave", IsPerson: false},
				{Name: "Emilia Pérez", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "Nickel Boys", WorkTitle: "Nickel Boys", IsPerson: false},
				{Name: "Sing Sing", WorkTitle: "Sing Sing", IsPerson: false},
			},
		},
		{
			Category: "Best Animated Feature",
			Nominees: []NomineeInfo{
				{Name: "Flow", WorkTitle: "Flow", IsPerson: false},
				{Name: "Inside Out 2", WorkTitle: "Inside Out 2", IsPerson: false},
				{Name: "Memoir of a Snail", WorkTitle: "Memoir of a Snail", IsPerson: false},
				{Name: "Wallace & Gromit: Vengeance Most Fowl", WorkTitle: "Wallace & Gromit: Vengeance Most Fowl", IsPerson: false},
				{Name: "The Wild Robot", WorkTitle: "The Wild Robot", IsPerson: false},
			},
		},
		{
			Category: "Best International Feature Film",
			Nominees: []NomineeInfo{
				{Name: "I'm Still Here", WorkTitle: "I'm Still Here", IsPerson: false},
				{Name: "The Girl with the Needle", WorkTitle: "The Girl with the Needle", IsPerson: false},
				{Name: "Emilia Pérez", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "The Seed of the Sacred Fig", WorkTitle: "The Seed of the Sacred Fig", IsPerson: false},
				{Name: "Flow", WorkTitle: "Flow", IsPerson: false},
			},
		},
		{
			Category: "Best Original Score",
			Nominees: []NomineeInfo{
				{Name: "The Brutalist", WorkTitle: "The Brutalist", IsPerson: false},
				{Name: "Conclave", WorkTitle: "Conclave", IsPerson: false},
				{Name: "Emilia Pérez", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "Wicked", WorkTitle: "Wicked", IsPerson: false},
				{Name: "The Wild Robot", WorkTitle: "The Wild Robot", IsPerson: false},
			},
		},
		{
			Category: "Best Original Song",
			Nominees: []NomineeInfo{
				{Name: "El Mal", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "The Journey", WorkTitle: "The Six Triple Eight", IsPerson: false},
				{Name: "Like a Bird", WorkTitle: "Sing Sing", IsPerson: false},
				{Name: "Mi Camino", WorkTitle: "Emilia Pérez", IsPerson: false},
				{Name: "Never Too Late", WorkTitle: "Elton John: Never Too Late", IsPerson: false},
			},
		},
	}
}

// GetCeremonyName returns the ceremony name for a given year
func GetCeremonyName(year int) string {
	// Calculate the ceremony number (1st Academy Awards was in 1929 for 1927/28 films)
	ceremonyNumber := year - 1928
	suffix := "th"
	if ceremonyNumber%10 == 1 && ceremonyNumber%100 != 11 {
		suffix = "st"
	} else if ceremonyNumber%10 == 2 && ceremonyNumber%100 != 12 {
		suffix = "nd"
	} else if ceremonyNumber%10 == 3 && ceremonyNumber%100 != 13 {
		suffix = "rd"
	}
	return fmt.Sprintf("%d%s Academy Awards", ceremonyNumber, suffix)
}
