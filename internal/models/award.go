package models

import "github.com/jackc/pgx/v5/pgtype"

type AwardType string

const (
	AwardTypeEmmy  AwardType = "Emmy"
	AwardTypeGrammy AwardType = "Grammy"
	AwardTypeOscar AwardType = "Oscar"
	AwardTypeTony  AwardType = "Tony"
)

type Award struct {
	ID           pgtype.UUID `json:"id" db:"id"`
	CelebrityID  pgtype.UUID `json:"celebrity_id" db:"celebrity_id"`
	Type         AwardType   `json:"type" db:"type"`
	Year         int         `json:"year" db:"year"`
	Work         string      `json:"work" db:"work"`
	Category     string      `json:"category" db:"category"`
	IsWinner     bool        `json:"is_winner" db:"is_winner"`
	CeremonyDate pgtype.Date `json:"ceremony_date,omitempty" db:"ceremony_date"`
	IsUpcoming   bool        `json:"is_upcoming" db:"is_upcoming"`
}
