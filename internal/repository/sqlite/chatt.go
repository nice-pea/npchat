package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

type ChattRepository struct {
	db *sqlx.DB
}

func (r *ChattRepository) List(filter chatt.Filter) ([]chatt.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ChattRepository) Upsert(chat chatt.Chat) error {
	//TODO implement me
	panic("implement me")
}
