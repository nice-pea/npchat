package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type SessionnRepository struct {
	db *sqlx.DB
}

func (r SessionnRepository) List(filter sessionn.Filter) ([]sessionn.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (r SessionnRepository) Upsert(session sessionn.Session) error {
	//TODO implement me
	panic("implement me")
}
