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

func (r *SessionnRepository) withTxConn(txConn baseRepo.DbConn) sessionn.Repository {
	return &SessionnRepository{
		BaseRepo: r.WithTxConn(txConn),
	}
}

func (r *SessionnRepository) InTransaction(fn func(txRepo sessionn.Repository) error) error {
	return inTransaction(r, fn)
}
