# EGOT Tracker

A web application that tracks celebrity progress toward achieving EGOT status (Emmy, Grammy, Oscar, Tony awards).

## Features

- Search for any celebrity and see their EGOT progress
- View celebrities who are "close to EGOT" (3 of 4 awards)
- Automatic data fetching from Wikidata
- Old Hollywood-themed UI

## Tech Stack

- **Backend**: Go 1.22+ with net/http
- **Frontend**: Next.js 14, React, Tailwind CSS
- **Database**: PostgreSQL
- **Data**: Wikidata SPARQL API

## Setup

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 16

### Quick Start

```bash
# All commands should be run from the project root directory
cd egot_watch

# 1. Start PostgreSQL
brew services start postgresql@16

# 2. Create database
# Note: Use /usr/local/opt/postgresql@16/bin/ for Intel Macs
/opt/homebrew/opt/postgresql@16/bin/createdb egot_tracker
/opt/homebrew/opt/postgresql@16/bin/psql -d egot_tracker -f setup.sql

# If you get errors about missing columns (ceremony_date, is_upcoming), run migrations:
# /opt/homebrew/opt/postgresql@16/bin/psql -d egot_tracker -f migration_add_ceremony_date.sql

# 3. Configure environment
echo "DATABASE_URL=postgresql://localhost:5432/egot_tracker?sslmode=disable" > .env
echo "PORT=8080" >> .env

# 4. Start backend
go run ./cmd/api

# 5. Setup and start frontend (new terminal)
# Install Node.js 18.20.4 if not already installed (asdf will auto-detect from .tool-versions)
asdf install nodejs 18.20.4

cd frontend && npm install && npm run dev
```

Visit http://localhost:3210

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/celebrity/search?q=NAME` | Search for a celebrity |
| `GET /api/celebrity/autocomplete?q=QUERY` | Autocomplete suggestions |
| `GET /api/celebrity/close-to-egot` | Get celebrities with 3/4 awards |
| `GET /health` | Health check |

## License

MIT
