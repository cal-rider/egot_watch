package scraper

// WikidataSearchResult represents a search result from Wikidata API
type WikidataSearchResult struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// WikidataSearchResponse represents the response from Wikidata search API
type WikidataSearchResponse struct {
	Search []WikidataSearchResult `json:"search"`
}

// WikidataAward represents an award parsed from SPARQL results
type WikidataAward struct {
	AwardID   string
	AwardName string
	Year      int
	Work      string
	Category  string
	IsWinner  bool
}

// WikidataPersonInfo contains basic person information from Wikidata
type WikidataPersonInfo struct {
	WikidataID string
	Name       string
	PhotoURL   string
}

// SPARQLResponse represents the response from Wikidata SPARQL endpoint
type SPARQLResponse struct {
	Results SPARQLResults `json:"results"`
}

type SPARQLResults struct {
	Bindings []SPARQLBinding `json:"bindings"`
}

type SPARQLBinding struct {
	Award       SPARQLValue `json:"award"`
	AwardLabel  SPARQLValue `json:"awardLabel"`
	Year        SPARQLValue `json:"year"`
	Work        SPARQLValue `json:"workLabel"`
	Image       SPARQLValue `json:"image"`
	PersonLabel SPARQLValue `json:"personLabel"`
}

type SPARQLValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
