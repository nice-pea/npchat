package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
)

type SessionnRepository struct {
	db *sqlx.DB
}

func (s SessionnRepository) List(filter sessionn.Filter) ([]sessionn.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (s SessionnRepository) Upsert(s2 sessionn.Session) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryFactory) NewSessionnRepository() sessionn.Repository {
	return &SessionnRepository{
		db: r.db,
	}
}
