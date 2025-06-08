package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

type UserrRepository struct {
	db *sqlx.DB
}

func (r *RepositoryFactory) NewUserrRepository() *UserrRepository {
	return &UserrRepository{
		db: r.db,
	}
}

func (u *UserrRepository) List(filter userr.Filter) ([]userr.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserrRepository) Upsert(u2 userr.User) error {
	//TODO implement me
	panic("implement me")
}
