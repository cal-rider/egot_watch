package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"egot-tracker/internal/models"
)

// Wikidata IDs for EGOT award types
const (
	EmmyAwardID   = "Q123538"
	GrammyAwardID = "Q41254"
	OscarAwardID  = "Q19020"
	TonyAwardID   = "Q191874"
)

// WikidataScraper fetches celebrity award data from Wikidata
type WikidataScraper struct {
	httpClient *http.Client
}

// NewWikidataScraper creates a new Wikidata scraper instance
func NewWikidataScraper() *WikidataScraper {
	return &WikidataScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchPerson searches for a person by name and returns their Wikidata ID
// Uses disambiguation to prefer people with EGOT awards
func (w *WikidataScraper) SearchPerson(ctx context.Context, name string) (*WikidataPersonInfo, error) {
	baseURL := "https://www.wikidata.org/w/api.php"
	params := url.Values{}
	params.Set("action", "wbsearchentities")
	params.Set("search", name)
	params.Set("language", "en")
	params.Set("type", "item")
	params.Set("limit", "5") // Get top 5 results for disambiguation
	params.Set("format", "json")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "EGOT-Tracker/1.0 (https://github.com/egot-tracker)")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search Wikidata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wikidata API returned status %d", resp.StatusCode)
	}

	var searchResp WikidataSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(searchResp.Search) == 0 {
		return nil, fmt.Errorf("no results found for: %s", name)
	}

	// Try to find a result with EGOT awards (disambiguation)
	for _, result := range searchResp.Search {
		hasEGOT, err := w.hasEGOTAwards(ctx, result.ID)
		if err != nil {
			continue // Skip on error, try next result
		}
		if hasEGOT {
			return &WikidataPersonInfo{
				WikidataID: result.ID,
				Name:       result.Label,
			}, nil
		}
	}

	// Fallback: return first result if none have EGOT awards
	result := searchResp.Search[0]
	return &WikidataPersonInfo{
		WikidataID: result.ID,
		Name:       result.Label,
	}, nil
}

// hasEGOTAwards checks if a Wikidata entity has any EGOT-type awards
func (w *WikidataScraper) hasEGOTAwards(ctx context.Context, wikidataID string) (bool, error) {
	// Quick SPARQL query to check for any EGOT awards
	query := fmt.Sprintf(`
ASK {
  wd:%s wdt:P166 ?award .
  ?award rdfs:label ?label .
  FILTER(LANG(?label) = "en")
  FILTER(
    CONTAINS(LCASE(?label), "emmy") ||
    CONTAINS(LCASE(?label), "grammy") ||
    CONTAINS(LCASE(?label), "academy award") ||
    CONTAINS(LCASE(?label), "oscar") ||
    CONTAINS(LCASE(?label), "tony award")
  )
}
`, wikidataID)

	sparqlURL := "https://query.wikidata.org/sparql"
	params := url.Values{}
	params.Set("query", query)
	params.Set("format", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sparqlURL+"?"+params.Encode(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "EGOT-Tracker/1.0 (https://github.com/egot-tracker)")
	req.Header.Set("Accept", "application/sparql-results+json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("SPARQL returned status %d", resp.StatusCode)
	}

	var askResp struct {
		Boolean bool `json:"boolean"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&askResp); err != nil {
		return false, err
	}

	return askResp.Boolean, nil
}

// GetPersonWithAwards fetches a person's info and all awards from Wikidata
// EGOT filtering is done in code via classifyAward() for reliability
func (w *WikidataScraper) GetPersonWithAwards(ctx context.Context, wikidataID string) (*WikidataPersonInfo, []WikidataAward, error) {
	// SPARQL query to get person info, photo, and ALL awards (filter in code)
	query := fmt.Sprintf(`
SELECT DISTINCT ?personLabel ?image ?award ?awardLabel ?year ?workLabel WHERE {
  wd:%s p:P166 ?statement .
  ?statement ps:P166 ?award .

  # Get award details
  OPTIONAL { ?statement pq:P585 ?date . BIND(YEAR(?date) AS ?year) }

  # Try multiple properties for the work (P1686=for work, P1411=nominated for, P972=catalog)
  OPTIONAL {
    ?statement pq:P1686|pq:P1411|pq:P972 ?work .
  }

  # Get person's image
  OPTIONAL { wd:%s wdt:P18 ?image }

  SERVICE wikibase:label { bd:serviceParam wikibase:language "en" }
}
ORDER BY DESC(?year)
`, wikidataID, wikidataID)

	sparqlURL := "https://query.wikidata.org/sparql"
	params := url.Values{}
	params.Set("query", query)
	params.Set("format", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sparqlURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SPARQL request: %w", err)
	}
	req.Header.Set("User-Agent", "EGOT-Tracker/1.0 (https://github.com/egot-tracker)")
	req.Header.Set("Accept", "application/sparql-results+json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query SPARQL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("SPARQL endpoint returned status %d", resp.StatusCode)
	}

	var sparqlResp SPARQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&sparqlResp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode SPARQL response: %w", err)
	}

	// Parse results
	var personInfo WikidataPersonInfo
	personInfo.WikidataID = wikidataID

	awards := make([]WikidataAward, 0)
	seenAwards := make(map[string]bool)

	for _, binding := range sparqlResp.Results.Bindings {
		// Get person info from first result
		if personInfo.Name == "" && binding.PersonLabel.Value != "" {
			personInfo.Name = binding.PersonLabel.Value
		}
		if personInfo.PhotoURL == "" && binding.Image.Value != "" {
			personInfo.PhotoURL = binding.Image.Value
		}

		// Parse award
		awardID := binding.Award.Value
		year := 0
		if binding.Year.Value != "" {
			if y, err := strconv.Atoi(binding.Year.Value); err == nil {
				year = y
			}
		}

		// Deduplicate by award + year (allows multiple wins in same category across different years)
		dedupeKey := fmt.Sprintf("%s-%d", awardID, year)
		if seenAwards[dedupeKey] {
			continue
		}
		seenAwards[dedupeKey] = true

		award := WikidataAward{
			AwardID:   awardID,
			AwardName: binding.AwardLabel.Value,
			Work:      binding.Work.Value,
			Category:  binding.AwardLabel.Value, // Use award name as category
			IsWinner:  true,                     // P166 is "award received", so these are wins
			Year:      year,
		}

		awards = append(awards, award)
	}

	return &personInfo, awards, nil
}

// WikipediaSummaryResponse represents the response from Wikipedia REST API
type WikipediaSummaryResponse struct {
	Extract string `json:"extract"`
}

// FetchWikipediaSummary fetches the summary/extract from Wikipedia for a given name
func (w *WikidataScraper) FetchWikipediaSummary(ctx context.Context, name string) (string, error) {
	// Wikipedia uses underscores for spaces in titles
	title := strings.ReplaceAll(name, " ", "_")
	apiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(title))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Wikipedia request: %w", err)
	}
	req.Header.Set("User-Agent", "EGOT-Tracker/1.0 (https://github.com/egot-tracker)")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Wikipedia summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Wikipedia API returned status %d", resp.StatusCode)
	}

	var summaryResp WikipediaSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&summaryResp); err != nil {
		return "", fmt.Errorf("failed to decode Wikipedia response: %w", err)
	}

	return summaryResp.Extract, nil
}

// FetchCelebrity searches for a celebrity and returns their data with awards
func (w *WikidataScraper) FetchCelebrity(ctx context.Context, name string) (*models.Celebrity, []models.Award, error) {
	// Step 1: Search for the person
	personInfo, err := w.SearchPerson(ctx, name)
	if err != nil {
		return nil, nil, err
	}

	// Step 2: Get their awards
	fullInfo, wikidataAwards, err := w.GetPersonWithAwards(ctx, personInfo.WikidataID)
	if err != nil {
		return nil, nil, err
	}

	// Use the name from search if SPARQL didn't return it
	if fullInfo.Name == "" {
		fullInfo.Name = personInfo.Name
	}

	// Step 3: Fetch Wikipedia summary (required - skip if not found)
	summary, err := w.FetchWikipediaSummary(ctx, fullInfo.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Wikipedia lookup failed for %s: %w", fullInfo.Name, err)
	}

	// Step 4: Convert to our models
	celebrity := &models.Celebrity{
		Name: fullInfo.Name,
		Slug: slugify(fullInfo.Name),
	}
	if fullInfo.PhotoURL != "" {
		celebrity.PhotoURL.String = fullInfo.PhotoURL
		celebrity.PhotoURL.Valid = true
	}
	if summary != "" {
		celebrity.Summary.String = summary
		celebrity.Summary.Valid = true
	}

	awards := make([]models.Award, 0, len(wikidataAwards))
	for _, wa := range wikidataAwards {
		awardType := classifyAward(wa.AwardName)
		if awardType == "" {
			continue // Skip non-EGOT awards
		}

		award := models.Award{
			Type:     awardType,
			Year:     wa.Year,
			Work:     wa.Work,
			Category: wa.Category,
			IsWinner: wa.IsWinner,
		}
		awards = append(awards, award)
	}

	return celebrity, awards, nil
}

// classifyAward determines the EGOT award type from the award name
func classifyAward(awardName string) models.AwardType {
	name := strings.ToLower(awardName)

	if strings.Contains(name, "emmy") {
		return models.AwardTypeEmmy
	}
	if strings.Contains(name, "grammy") {
		return models.AwardTypeGrammy
	}
	if strings.Contains(name, "academy award") || strings.Contains(name, "oscar") {
		return models.AwardTypeOscar
	}
	if strings.Contains(name, "tony") {
		return models.AwardTypeTony
	}

	return ""
}

// slugify converts a name to a URL-friendly slug
func slugify(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, "'", "")
	return slug
}
