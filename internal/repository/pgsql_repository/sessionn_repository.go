package pgsqlRepository

import (
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	baseRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/base_repo"
)

type SessionnRepository struct {
	baseRepo.BaseRepo
}

func (r *SessionnRepository) List(filter sessionn.Filter) ([]sessionn.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SessionnRepository) Upsert(session sessionn.Session) error {
	//TODO implement me
	panic("implement me")
}

func (r *SessionnRepository) WithTxConn(txConn baseRepo.DbConn) sessionn.Repository {
	return &SessionnRepository{
		BaseRepo: r.BaseRepo.WithTxConn(txConn),
	}
}

func (r *SessionnRepository) InTransaction(fn func(txRepo sessionn.Repository) error) error {
	return baseRepo.InTransaction(r, fn)
}
