package pgsqlRepository

import (
	"github.com/nice-pea/npchat/internal/domain/chatt"
	baseRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/base_repo"
)

type ChattRepository struct {
	baseRepo.BaseRepo
}

func (r *ChattRepository) List(filter chatt.Filter) ([]chatt.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ChattRepository) Upsert(chat chatt.Chat) error {
	//TODO implement me
	panic("implement me")
}

func (r *ChattRepository) withTxConn(txConn baseRepo.DbConn) chatt.Repository {
	return &ChattRepository{
		BaseRepo: r.WithTxConn(txConn),
	}
}

func (r *ChattRepository) InTransaction(fn func(txRepo chatt.Repository) error) error {
	return inTransaction(r, fn)
}
