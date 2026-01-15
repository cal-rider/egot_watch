package repository

import (
	"context"

	"egot-tracker/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AwardRepository struct {
	pool *pgxpool.Pool
}

func NewAwardRepository(pool *pgxpool.Pool) *AwardRepository {
	return &AwardRepository{pool: pool}
}

func (r *AwardRepository) FindByCelebrityID(ctx context.Context, celebrityID pgtype.UUID) ([]models.Award, error) {
	query := `
		SELECT id, celebrity_id, type, year, work, category, is_winner, ceremony_date, is_upcoming
		FROM awards
		WHERE celebrity_id = $1
		ORDER BY year DESC, type
	`

	rows, err := r.pool.Query(ctx, query, celebrityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []models.Award
	for rows.Next() {
		var award models.Award
		err := rows.Scan(
			&award.ID,
			&award.CelebrityID,
			&award.Type,
			&award.Year,
			&award.Work,
			&award.Category,
			&award.IsWinner,
			&award.CeremonyDate,
			&award.IsUpcoming,
		)
		if err != nil {
			return nil, err
		}
		awards = append(awards, award)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return awards, nil
}

func (r *AwardRepository) CreateBatch(ctx context.Context, celebrityID pgtype.UUID, awards []models.Award) ([]models.Award, error) {
	if len(awards) == 0 {
		return []models.Award{}, nil
	}

	query := `
		INSERT INTO awards (celebrity_id, type, year, work, category, is_winner, ceremony_date, is_upcoming)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, celebrity_id, type, year, work, category, is_winner, ceremony_date, is_upcoming
	`

	created := make([]models.Award, 0, len(awards))
	for _, award := range awards {
		var a models.Award
		err := r.pool.QueryRow(ctx, query,
			celebrityID,
			award.Type,
			award.Year,
			award.Work,
			award.Category,
			award.IsWinner,
			award.CeremonyDate,
			award.IsUpcoming,
		).Scan(
			&a.ID,
			&a.CelebrityID,
			&a.Type,
			&a.Year,
			&a.Work,
			&a.Category,
			&a.IsWinner,
			&a.CeremonyDate,
			&a.IsUpcoming,
		)
		if err != nil {
			return nil, err
		}
		created = append(created, a)
	}

	return created, nil
}
