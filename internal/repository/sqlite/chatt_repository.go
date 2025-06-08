package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

type ChattRepository struct {
	db *sqlx.DB
}

func (r *RepositoryFactory) NewChattRepository() chatt.Repository {
	return &ChattRepository{
		db: r.db,
	}
}

func (c *ChattRepository) List(filter chatt.Filter) ([]chatt.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChattRepository) Upsert(c2 chatt.Chat) error {
	//TODO implement me
	panic("implement me")
}
