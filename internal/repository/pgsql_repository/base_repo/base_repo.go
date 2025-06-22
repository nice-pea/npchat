package baseRepo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type BaseRepo struct {
	db    DbConn
	txBeg TxBeginner
	isTx  bool // Признак выполнения в транзакции
}

func New(conn interface {
	DbConn
	TxBeginner
}) BaseRepo {
	return BaseRepo{
		db:    conn,
		txBeg: conn,
		isTx:  false,
	}
}

func (r BaseRepo) WithTxConn(txConn DbConn) BaseRepo {
	return BaseRepo{
		db:    txConn,
		txBeg: r.txBeg,
		isTx:  true,
	}
}

func (r BaseRepo) DB() DbConn {
	return r.db
}

func (r BaseRepo) TxBeginner() TxBeginner {
	return r.txBeg
}

func (r BaseRepo) IsTx() bool {
	return r.isTx
}

type DbConn interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

type TxBeginner interface {
	Beginx() (*sqlx.Tx, error)
}
