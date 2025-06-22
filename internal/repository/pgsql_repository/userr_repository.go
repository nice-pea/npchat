package pgsqlRepository

import (
	"github.com/nice-pea/npchat/internal/domain/userr"
	baseRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/base_repo"
)

type UserrRepository struct {
	baseRepo.BaseRepo
}

func (r *UserrRepository) List(filter userr.Filter) ([]userr.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserrRepository) Upsert(user userr.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserrRepository) WithTxConn(txConn baseRepo.DbConn) userr.Repository {
	return &UserrRepository{
		BaseRepo: r.BaseRepo.WithTxConn(txConn),
	}
}

func (r *UserrRepository) InTransaction(fn func(txRepo userr.Repository) error) error {
	return baseRepo.InTransaction(r, fn)
}
