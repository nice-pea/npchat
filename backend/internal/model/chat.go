package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Chat struct {
	ID        uint          `gorm:"primaryKey" json:"id" db:"id"`
	Name      string        `json:"name" db:"name"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	CreatorID pgtype.Uint32 `json:"creator_id" db:"creator_id"`
}
