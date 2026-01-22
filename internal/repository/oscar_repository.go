package repository

import (
	"context"
	"errors"

	"egot-tracker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCeremonyNotFound = errors.New("ceremony not found")

type OscarRepository struct {
	pool *pgxpool.Pool
}

func NewOscarRepository(pool *pgxpool.Pool) *OscarRepository {
	return &OscarRepository{pool: pool}
}

// CreateCeremony creates a new Oscar ceremony
func (r *OscarRepository) CreateCeremony(ctx context.Context, ceremony *models.OscarCeremony) (*models.OscarCeremony, error) {
	query := `
		INSERT INTO oscar_ceremonies (year, ceremony_name, ceremony_date, is_complete)
		VALUES ($1, $2, $3, $4)
		RETURNING id, year, ceremony_name, ceremony_date, is_complete, created_at
	`

	var created models.OscarCeremony
	err := r.pool.QueryRow(ctx, query,
		ceremony.Year,
		ceremony.CeremonyName,
		ceremony.CeremonyDate,
		ceremony.IsComplete,
	).Scan(
		&created.ID,
		&created.Year,
		&created.CeremonyName,
		&created.CeremonyDate,
		&created.IsComplete,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetCeremonyByYear fetches a ceremony by year
func (r *OscarRepository) GetCeremonyByYear(ctx context.Context, year int) (*models.OscarCeremony, error) {
	query := `
		SELECT id, year, ceremony_name, ceremony_date, is_complete, created_at
		FROM oscar_ceremonies
		WHERE year = $1
	`

	var ceremony models.OscarCeremony
	err := r.pool.QueryRow(ctx, query, year).Scan(
		&ceremony.ID,
		&ceremony.Year,
		&ceremony.CeremonyName,
		&ceremony.CeremonyDate,
		&ceremony.IsComplete,
		&ceremony.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCeremonyNotFound
	}
	if err != nil {
		return nil, err
	}

	return &ceremony, nil
}

// GetAllCeremonyYears returns all tracked Oscar years
func (r *OscarRepository) GetAllCeremonyYears(ctx context.Context) ([]int, error) {
	query := `SELECT year FROM oscar_ceremonies ORDER BY year DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			return nil, err
		}
		years = append(years, year)
	}

	return years, rows.Err()
}

// CreateCategory creates a new Oscar category
func (r *OscarRepository) CreateCategory(ctx context.Context, category *models.OscarCategory) (*models.OscarCategory, error) {
	query := `
		INSERT INTO oscar_categories (ceremony_id, name, display_order, winner_announced)
		VALUES ($1, $2, $3, $4)
		RETURNING id, ceremony_id, name, display_order, winner_announced
	`

	var created models.OscarCategory
	err := r.pool.QueryRow(ctx, query,
		category.CeremonyID,
		category.Name,
		category.DisplayOrder,
		category.WinnerAnnounced,
	).Scan(
		&created.ID,
		&created.CeremonyID,
		&created.Name,
		&created.DisplayOrder,
		&created.WinnerAnnounced,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetCategoriesByCeremony fetches all categories for a ceremony
func (r *OscarRepository) GetCategoriesByCeremony(ctx context.Context, ceremonyID pgtype.UUID) ([]models.OscarCategory, error) {
	query := `
		SELECT id, ceremony_id, name, display_order, winner_announced
		FROM oscar_categories
		WHERE ceremony_id = $1
		ORDER BY display_order
	`

	rows, err := r.pool.Query(ctx, query, ceremonyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.OscarCategory
	for rows.Next() {
		var c models.OscarCategory
		err := rows.Scan(&c.ID, &c.CeremonyID, &c.Name, &c.DisplayOrder, &c.WinnerAnnounced)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

// CreateNominee creates a new Oscar nominee
func (r *OscarRepository) CreateNominee(ctx context.Context, nominee *models.OscarNominee) (*models.OscarNominee, error) {
	query := `
		INSERT INTO oscar_nominees (category_id, celebrity_id, name, photo_url, work_title, is_winner, display_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, category_id, celebrity_id, name, photo_url, work_title, is_winner, display_order
	`

	var created models.OscarNominee
	err := r.pool.QueryRow(ctx, query,
		nominee.CategoryID,
		nominee.CelebrityID,
		nominee.Name,
		nominee.PhotoURL,
		nominee.WorkTitle,
		nominee.IsWinner,
		nominee.DisplayOrder,
	).Scan(
		&created.ID,
		&created.CategoryID,
		&created.CelebrityID,
		&created.Name,
		&created.PhotoURL,
		&created.WorkTitle,
		&created.IsWinner,
		&created.DisplayOrder,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetNomineesByCategory fetches all nominees for a category
func (r *OscarRepository) GetNomineesByCategory(ctx context.Context, categoryID pgtype.UUID) ([]models.OscarNominee, error) {
	query := `
		SELECT id, category_id, celebrity_id, name, photo_url, work_title, is_winner, display_order
		FROM oscar_nominees
		WHERE category_id = $1
		ORDER BY display_order
	`

	rows, err := r.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nominees []models.OscarNominee
	for rows.Next() {
		var n models.OscarNominee
		err := rows.Scan(&n.ID, &n.CategoryID, &n.CelebrityID, &n.Name, &n.PhotoURL, &n.WorkTitle, &n.IsWinner, &n.DisplayOrder)
		if err != nil {
			return nil, err
		}
		nominees = append(nominees, n)
	}

	return nominees, rows.Err()
}

// GetFullCeremony fetches a ceremony with all categories and nominees
func (r *OscarRepository) GetFullCeremony(ctx context.Context, year int) (*models.OscarCeremonyFull, error) {
	// Get the ceremony
	ceremony, err := r.GetCeremonyByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	// Get all categories
	categories, err := r.GetCategoriesByCeremony(ctx, ceremony.ID)
	if err != nil {
		return nil, err
	}

	// Build full result
	result := &models.OscarCeremonyFull{
		OscarCeremony: *ceremony,
		Categories:    make([]models.OscarCategoryWithNominees, len(categories)),
	}

	// Get nominees for each category
	for i, cat := range categories {
		nominees, err := r.GetNomineesByCategory(ctx, cat.ID)
		if err != nil {
			return nil, err
		}

		result.Categories[i] = models.OscarCategoryWithNominees{
			OscarCategory: cat,
			Nominees:      nominees,
		}
	}

	return result, nil
}

// SetNomineeAsWinner marks a nominee as the winner
func (r *OscarRepository) SetNomineeAsWinner(ctx context.Context, nomineeID pgtype.UUID) error {
	// First, get the category ID for this nominee
	var categoryID pgtype.UUID
	err := r.pool.QueryRow(ctx, "SELECT category_id FROM oscar_nominees WHERE id = $1", nomineeID).Scan(&categoryID)
	if err != nil {
		return err
	}

	// Reset all winners in this category
	_, err = r.pool.Exec(ctx, "UPDATE oscar_nominees SET is_winner = false WHERE category_id = $1", categoryID)
	if err != nil {
		return err
	}

	// Set this nominee as winner
	_, err = r.pool.Exec(ctx, "UPDATE oscar_nominees SET is_winner = true WHERE id = $1", nomineeID)
	if err != nil {
		return err
	}

	// Mark category as winner announced
	_, err = r.pool.Exec(ctx, "UPDATE oscar_categories SET winner_announced = true WHERE id = $1", categoryID)
	return err
}

// DeleteCeremony removes a ceremony and all related data
func (r *OscarRepository) DeleteCeremony(ctx context.Context, year int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM oscar_ceremonies WHERE year = $1", year)
	return err
}
