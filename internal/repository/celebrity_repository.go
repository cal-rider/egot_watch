package repository

import (
	"context"
	"errors"
	"strings"

	"egot-tracker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCelebrityNotFound = errors.New("celebrity not found")

type CelebrityRepository struct {
	pool *pgxpool.Pool
}

func NewCelebrityRepository(pool *pgxpool.Pool) *CelebrityRepository {
	return &CelebrityRepository{pool: pool}
}

func (r *CelebrityRepository) FindByName(ctx context.Context, name string) (*models.Celebrity, error) {
	query := `
		SELECT id, name, slug, photo_url, summary, last_updated
		FROM celebrities
		WHERE LOWER(name) = LOWER($1)
	`

	var celebrity models.Celebrity
	err := r.pool.QueryRow(ctx, query, strings.TrimSpace(name)).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Slug,
		&celebrity.PhotoURL,
		&celebrity.Summary,
		&celebrity.LastUpdated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCelebrityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &celebrity, nil
}

func (r *CelebrityRepository) FindByID(ctx context.Context, id pgtype.UUID) (*models.Celebrity, error) {
	query := `
		SELECT id, name, slug, photo_url, summary, last_updated
		FROM celebrities
		WHERE id = $1
	`

	var celebrity models.Celebrity
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Slug,
		&celebrity.PhotoURL,
		&celebrity.Summary,
		&celebrity.LastUpdated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCelebrityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &celebrity, nil
}

func (r *CelebrityRepository) Search(ctx context.Context, query string, limit int) ([]models.Celebrity, error) {
	sql := `
		SELECT id, name, slug, photo_url, summary, last_updated
		FROM celebrities
		WHERE LOWER(name) LIKE LOWER($1)
		ORDER BY name
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, sql, "%"+strings.TrimSpace(query)+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var celebrities []models.Celebrity
	for rows.Next() {
		var c models.Celebrity
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.PhotoURL, &c.Summary, &c.LastUpdated)
		if err != nil {
			return nil, err
		}
		celebrities = append(celebrities, c)
	}

	return celebrities, rows.Err()
}

// FindCloseToEGOT returns celebrities with exactly 3 unique EGOT award wins
func (r *CelebrityRepository) FindCloseToEGOT(ctx context.Context, limit int) ([]models.CelebrityWithEGOTProgress, error) {
	query := `
		WITH celebrity_wins AS (
			SELECT
				c.id,
				c.name,
				c.slug,
				c.photo_url,
				c.summary,
				c.last_updated,
				COUNT(DISTINCT a.type) as egot_win_count,
				ARRAY_AGG(DISTINCT a.type::text ORDER BY a.type::text) as won_awards
			FROM celebrities c
			INNER JOIN awards a ON c.id = a.celebrity_id
			WHERE a.is_winner = true AND a.is_upcoming = false
			GROUP BY c.id, c.name, c.slug, c.photo_url, c.summary, c.last_updated
			HAVING COUNT(DISTINCT a.type) = 3
		)
		SELECT id, name, slug, photo_url, summary, last_updated, egot_win_count, won_awards
		FROM celebrity_wins
		ORDER BY name
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var celebrities []models.CelebrityWithEGOTProgress
	for rows.Next() {
		var c models.CelebrityWithEGOTProgress
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.PhotoURL,
			&c.Summary,
			&c.LastUpdated,
			&c.EGOTWinCount,
			&c.WonAwards,
		)
		if err != nil {
			return nil, err
		}
		celebrities = append(celebrities, c)
	}

	return celebrities, rows.Err()
}

// FindEGOTWinners returns celebrities with all 4 unique EGOT award wins
func (r *CelebrityRepository) FindEGOTWinners(ctx context.Context, limit int) ([]models.CelebrityWithEGOTProgress, error) {
	query := `
		WITH celebrity_wins AS (
			SELECT
				c.id,
				c.name,
				c.slug,
				c.photo_url,
				c.summary,
				c.last_updated,
				COUNT(DISTINCT a.type) as egot_win_count,
				ARRAY_AGG(DISTINCT a.type::text ORDER BY a.type::text) as won_awards
			FROM celebrities c
			INNER JOIN awards a ON c.id = a.celebrity_id
			WHERE a.is_winner = true AND a.is_upcoming = false
			GROUP BY c.id, c.name, c.slug, c.photo_url, c.summary, c.last_updated
			HAVING COUNT(DISTINCT a.type) = 4
		)
		SELECT id, name, slug, photo_url, summary, last_updated, egot_win_count, won_awards
		FROM celebrity_wins
		ORDER BY name
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var celebrities []models.CelebrityWithEGOTProgress
	for rows.Next() {
		var c models.CelebrityWithEGOTProgress
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.PhotoURL,
			&c.Summary,
			&c.LastUpdated,
			&c.EGOTWinCount,
			&c.WonAwards,
		)
		if err != nil {
			return nil, err
		}
		celebrities = append(celebrities, c)
	}

	return celebrities, rows.Err()
}

// FindNoAwards returns celebrities with no awards
func (r *CelebrityRepository) FindNoAwards(ctx context.Context, limit int) ([]models.Celebrity, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.photo_url, c.summary, c.last_updated
		FROM celebrities c
		LEFT JOIN awards a ON c.id = a.celebrity_id
		WHERE a.id IS NULL
		ORDER BY c.last_updated ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var celebrities []models.Celebrity
	for rows.Next() {
		var c models.Celebrity
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.PhotoURL, &c.Summary, &c.LastUpdated)
		if err != nil {
			return nil, err
		}
		celebrities = append(celebrities, c)
	}

	return celebrities, rows.Err()
}

func (r *CelebrityRepository) Create(ctx context.Context, celebrity *models.Celebrity) (*models.Celebrity, error) {
	query := `
		INSERT INTO celebrities (name, slug, photo_url, summary)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, photo_url, summary, last_updated
	`

	var created models.Celebrity
	err := r.pool.QueryRow(ctx, query,
		celebrity.Name,
		celebrity.Slug,
		celebrity.PhotoURL,
		celebrity.Summary,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Slug,
		&created.PhotoURL,
		&created.Summary,
		&created.LastUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
