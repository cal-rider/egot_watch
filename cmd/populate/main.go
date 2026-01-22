package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"

	"egot-tracker/internal/config"
	"egot-tracker/internal/database"
	"egot-tracker/internal/repository"
	"egot-tracker/internal/scraper"
	"egot-tracker/internal/service"
)

// List of celebrities to populate (deduplicated)
var celebrities = []string{

	"Donald Trump",
	"Richard Rodgers",
	"Helen Hayes",
	"Rita Moreno",
	"John Gielgud",
	"Audrey Hepburn",
	"Marvin Hamlisch",
	"Jonathan Tunick",
	"Mel Brooks",
	"Mike Nichols",
	"Whoopi Goldberg",
	"Scott Rudin",
	"Robert Lopez",
	"Andrew Lloyd Webber",
	"Tim Rice",
	"John Legend",
	"Alan Menken",
	"Jennifer Hudson",
	"Viola Davis",
	"Elton John",
	"Benj Pasek",
	"Justin Paul",
	"Barbra Streisand",
	"Liza Minnelli",
	"James Earl Jones",
	"Harry Belafonte",
	"Quincy Jones",
	"Frank Marshall",
	"Lin-Manuel Miranda",
	"Stephen Sondheim",
	"Cyndi Lauper",
	"Hugh Jackman",
	"Billy Porter",
	"Audra McDonald",
	"Bette Midler",
	"Cher",
	"Kate Winslet",
	"Jessica Lange",
	"Jeremy Irons",
	"Al Pacino",
	"Helen Mirren",
	"Maggie Smith",
	"Frances McDormand",
	"Christopher Plummer",
	"Vanessa Redgrave",
	"Ellen Burstyn",
	"Glenda Jackson",
	"Common",
	"John Williams",
	"Randy Newman",
	"Julie Andrews",
	"Dick Van Dyke",
	"Burt Bacharach",
	"Cynthia Erivo",
	"Adele",
	"H.E.R.",
	"Paul McCartney",
	"Ringo Starr",
	"Martin Scorsese",
	"Trent Reznor",
	"Atticus Ross",
	"Hildur Guðnadóttir",
	"Ben Platt",
	"Lily Tomlin",
	"Billie Eilish",
	"Finneas O'Connell",
	"Lady Gaga",
	"Beyoncé",
	"Eminem",
	"Kendrick Lamar",
	"Ludwig Göransson",
	"Allison Janney",
	"Meryl Streep",
	"Denzel Washington",
	"Anne Hathaway",
	"Cynthia Nixon",
	"Bryan Cranston",
	"Octavia Spencer",
	"Regina King",
	"Laura Dern",
	"J.K. Simmons",
	"Sam Rockwell",
	"Dustin Hoffman",
	"Morgan Freeman",
	"Donald Glover",
	"Mark Ronson",
	"Hans Zimmer",
	"Aretha Franklin",
	"Stevie Wonder",
	"Paul Simon",
	"Bruce Springsteen",
	"Prince",
	"Michael Jackson",
	"Madonna",
	"Whitney Houston",
	"Mariah Carey",
	"Celine Dion",
	"Dolly Parton",
	"Billy Joel",
	"Kevin Spacey",
	"Angelina Jolie",
	"Gwyneth Paltrow",
	"Judi Dench",
	"Susan Sarandon",
	"Tom Hanks",
	"Emma Thompson",
	"Anthony Hopkins",
	"Michael Douglas",
	"Shirley MacLaine",
	"Bette Davis",
	"Sidney Poitier",
	"William Holden",
	"Yul Brynner",
	"Jon Batiste",
	"Questlove",
	"Riz Ahmed",
	"George Clooney",
	"Brad Pitt",
	"Joaquin Phoenix",
	"Olivia Colman",
	"Mark Rylance",
	"Julianne Moore",
	"Patricia Arquette",
	"Lupita Nyong'o",
	"Reese Witherspoon",
	"Nicole Kidman",
	"Halle Berry",
	"Julia Roberts",
	"Casey Affleck",
	"Christian Bale",
	"Philip Seymour Hoffman",
	"Robin Williams",
	"Gene Hackman",
	"Robert De Niro",
	"Jack Nicholson",
	"Samuel L. Jackson",
	"Zendaya",
	"Jeremy Allen White",
	"Ayo Edebiri",
	"Sarah Snook",
	"Kieran Culkin",
	"Jennifer Coolidge",
	"Quinta Brunson",
	"Steven Yeun",
	"Ali Wong",
	"Tina Fey",
	"Amy Poehler",
	"Jerry Seinfeld",
	"Phoebe Waller-Bridge",
	"Ariana Grande",
	"Drake",
	"SZA",
	"Miley Cyrus",
	"Harry Styles",
	"Dua Lipa",
	"Ed Sheeran",
	"Bruno Mars",
	"Miles Davis",
	"Mick Jagger",
	"David Bowie",
	"Dave Grohl",
	"Nicolas Cage",
	"Jodie Foster",
	"Kathy Bates",
	"Joe Pesci",
	"Sean Connery",
	"Paul Newman",
	"William Hurt",
	"Robert Duvall",
	"Marlon Brando",
	"Elizabeth Taylor",
	"Rex Harrison",
	"Gregory Peck",
	"Sophia Loren",
	"Burt Lancaster",
	"Alec Guinness",
	"Ernest Borgnine",
	"Grace Kelly",
	"Frank Sinatra",
	"Gary Cooper",
	"Humphrey Bogart",
	"Vivien Leigh",
	"Laurence Olivier",
	"Bing Crosby",
	"Spencer Tracy",
	"James Stewart",
	"Ginger Rogers",
	"Janet Gaynor",
	"Victoria Monét",
	"Ben Affleck",
	"Matt Damon",
	"Mahershala Ali",
	"Brie Larson",
	"Alicia Vikander",
	"Matthew McConaughey",
	"Cate Blanchett",
	"Jared Leto",
	"Jennifer Lawrence",
	"Colin Firth",
	"Natalie Portman",
	"Sandra Bullock",
	"Sean Penn",
	"Tilda Swinton",
	"Forest Whitaker",
	"Jamie Foxx",
	"Hilary Swank",
	"Robert Downey Jr.",
	"Cillian Murphy",
	"Da'Vine Joy Randolph",
	"Michelle Yeoh",
	"Brendan Fraser",
	"Ke Huy Quan",
	"Will Smith",
	"Jessica Chastain",
	"Troy Kotsur",
	"Yuh-jung Youn",
	"Renee Zellweger",
	"Gary Oldman",
	"Emma Stone",
	"Leonardo DiCaprio",
	"Eddie Redmayne",
	"Daniel Day-Lewis",
	"Christoph Waltz",
	"Jean Dujardin",
	"Melissa Leo",
	"Jeff Bridges",
	"Mo'Nique",
	"Heath Ledger",
	"Penelope Cruz",
	"Alan Arkin",
	"Rachel Weisz",
	"Tim Robbins",
	"Adrien Brody",
	"Chris Cooper",
	"Catherine Zeta-Jones",
	"Jim Broadbent",
	"Jennifer Connelly",
	"Russell Crowe",
	"Benicio del Toro",
	"Marcia Gay Harden",
	"Annette Bening",
	"Michael Caine",
	"Roberto Benigni",
	"James Coburn",
	"Helen Hunt",
	"Kim Basinger",
	"Geoffrey Rush",
	"Cuba Gooding Jr.",
	"Juliette Binoche",
	"Mira Sorvino",
	"Martin Landau",
	"Dianne Wiest",
	"Holly Hunter",
	"Tommy Lee Jones",
	"Anna Paquin",
	"Marisa Tomei",
	"Jack Palance",
	"Mercedes Ruehl",
	"Brenda Fricker",
	"Kevin Kline",
	"Geena Davis",
	"Olympia Dukakis",
	"Marlee Matlin",
	"Anjelica Huston",
	"F. Murray Abraham",
	"Haing S. Ngor",
	"Peggy Ashcroft",
	"Linda Hunt",
	"Ben Kingsley",
	"Louis Gossett Jr.",
	"Henry Fonda",
	"Maureen Stapleton",
	"Sissy Spacek",
	"Timothy Hutton",
	"Mary Steenburgen",
	"Jon Voight",
	"Jane Fonda",
	"Christopher Walken",
	"Richard Dreyfuss",
	"Diane Keaton",
	"Jason Robards",
	"Peter Finch",
	"Faye Dunaway",
	"Beatrice Straight",
	"Louise Fletcher",
	"George Burns",
	"Lee Grant",
	"Art Carney",
	"John Houseman",
	"Tatum O'Neal",
	"Joel Grey",
	"Eileen Heckart",
	"Ben Johnson",
	"Cloris Leachman",
	"George C. Scott",
	"John Mills",
	"Gig Young",
	"Goldie Hawn",
	"Cliff Robertson",
	"Ruth Gordon",
	"Katharine Hepburn",
	"Rod Steiger",
	"George Kennedy",
	"Estelle Parsons",
	"Paul Scofield",
	"Walter Matthau",
	"Sandy Dennis",
	"Lee Marvin",
	"Julie Christie",
	"Martin Balsam",
	"Shelley Winters",
	"Lila Kedrova",
	"Peter Ustinov",
	"Patricia Neal",
	"Margaret Rutherford",
	"Maximilian Schell",
	"Anne Bancroft",
	"Ed Begley",
	"Patty Duke",
	"Joan Fontaine",
	"Thomas Mitchell",
	"Hattie McDaniel",
	"Charles Laughton",
	"Marie Dressler",
	"Mary Pickford",
	"Mikey Madison",
	"Sean Baker",
	"Christopher Nolan",
	"Zoe Saldaña",
}

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Use a longer timeout for the entire operation
	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	celebrityRepo := repository.NewCelebrityRepository(pool)
	awardRepo := repository.NewAwardRepository(pool)
	wikidataScraper := scraper.NewWikidataScraper()
	celebrityService := service.NewCelebrityService(celebrityRepo, awardRepo, wikidataScraper)

	successCount := 0
	skipCount := 0
	failCount := 0

	log.Printf("Starting population of %d celebrities...\n", len(celebrities))

	for i, name := range celebrities {
		log.Printf("[%d/%d] Processing: %s", i+1, len(celebrities), name)

		// Check if already exists
		existing, _ := celebrityRepo.FindByName(ctx, name)
		if existing != nil {
			log.Printf("  → Skipped (already exists)")
			skipCount++
			continue
		}

		// Create a context with timeout for each request
		reqCtx, cancel := context.WithTimeout(ctx, 60*time.Second)

		// Use the service to search (which triggers Wikidata fetch)
		result, err := celebrityService.SearchCelebrity(reqCtx, name)
		cancel()

		if err != nil {
			log.Printf("  → Failed: %v", err)
			failCount++
		} else {
			log.Printf("  → Success: %d awards found", len(result.Awards))
			successCount++
		}

		// Delay between requests to be respectful to APIs
		if i < len(celebrities)-1 {
			time.Sleep(1500 * time.Millisecond)
		}
	}

	log.Println("\n=== Population Complete ===")
	log.Printf("Success: %d", successCount)
	log.Printf("Skipped: %d", skipCount)
	log.Printf("Failed:  %d", failCount)
	log.Printf("Total:   %d", len(celebrities))
}
