package sqlite

import (
	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

type UserrRepository struct {
	db *sqlx.DB
}

func (r *UserrRepository) List(filter userr.Filter) ([]userr.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserrRepository) Upsert(user userr.User) error {
	//TODO implement me
	panic("implement me")
}
