package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
)

type Message struct {
	ID        uint             `json:"id"`
	ChatID    uint             `json:"chat_id"`
	Text      string           `json:"text"`
	AuthorID  null.Uint        `json:"author_id,omitempty"`
	ReplyToID null.Uint        `json:"reply_to_id,omitempty"`
	EditedAt  pgtype.Timestamp `json:"edited_at,omitempty"`
	RemovedAt pgtype.Timestamp `json:"removed_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}
