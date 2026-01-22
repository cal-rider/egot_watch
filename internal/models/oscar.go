package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// OscarCeremony represents an annual Oscar ceremony
type OscarCeremony struct {
	ID           pgtype.UUID      `json:"id" db:"id"`
	Year         int              `json:"year" db:"year"`
	CeremonyName pgtype.Text      `json:"ceremony_name" db:"ceremony_name"`
	CeremonyDate pgtype.Date      `json:"ceremony_date" db:"ceremony_date"`
	IsComplete   bool             `json:"is_complete" db:"is_complete"`
	CreatedAt    pgtype.Timestamp `json:"created_at" db:"created_at"`
}

// OscarCategory represents a category within an Oscar ceremony
type OscarCategory struct {
	ID              pgtype.UUID `json:"id" db:"id"`
	CeremonyID      pgtype.UUID `json:"ceremony_id" db:"ceremony_id"`
	Name            string      `json:"name" db:"name"`
	DisplayOrder    int         `json:"display_order" db:"display_order"`
	WinnerAnnounced bool        `json:"winner_announced" db:"winner_announced"`
}

// OscarNominee represents a nominee in an Oscar category
type OscarNominee struct {
	ID           pgtype.UUID `json:"id" db:"id"`
	CategoryID   pgtype.UUID `json:"category_id" db:"category_id"`
	CelebrityID  pgtype.UUID `json:"celebrity_id,omitempty" db:"celebrity_id"`
	Name         string      `json:"name" db:"name"`
	PhotoURL     pgtype.Text `json:"photo_url" db:"photo_url"`
	WorkTitle    pgtype.Text `json:"work_title" db:"work_title"`
	IsWinner     bool        `json:"is_winner" db:"is_winner"`
	DisplayOrder int         `json:"display_order" db:"display_order"`
}

// OscarCategoryWithNominees combines a category with its nominees
type OscarCategoryWithNominees struct {
	OscarCategory
	Nominees []OscarNominee `json:"nominees"`
}

// OscarCeremonyFull represents a ceremony with all categories and nominees
type OscarCeremonyFull struct {
	OscarCeremony
	Categories []OscarCategoryWithNominees `json:"categories"`
}
