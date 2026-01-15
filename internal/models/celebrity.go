package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Celebrity struct {
	ID          pgtype.UUID      `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Slug        string           `json:"slug" db:"slug"`
	PhotoURL    pgtype.Text      `json:"photo_url" db:"photo_url"`
	LastUpdated pgtype.Timestamp `json:"last_updated" db:"last_updated"`
}

type CelebrityWithAwards struct {
	Celebrity
	Awards []Award `json:"awards"`
}

// CelebrityWithEGOTProgress represents a celebrity with their EGOT win count
type CelebrityWithEGOTProgress struct {
	Celebrity
	EGOTWinCount int      `json:"egot_win_count"`
	WonAwards    []string `json:"won_awards"` // e.g., ["Emmy", "Grammy", "Oscar"]
}

func (c *Celebrity) GetLastUpdated() *time.Time {
	if c.LastUpdated.Valid {
		return &c.LastUpdated.Time
	}
	return nil
}

func (c *Celebrity) GetPhotoURL() *string {
	if c.PhotoURL.Valid {
		return &c.PhotoURL.String
	}
	return nil
}
