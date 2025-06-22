package sqlxRepo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SqlxRepo struct {
	db    DbConn
	txBeg TxBeginner
	isTx  bool // Признак выполнения в транзакции
}

func New(sqlxDB interface {
	DbConn
	TxBeginner
}) SqlxRepo {
	return SqlxRepo{
		db:    sqlxDB,
		txBeg: sqlxDB,
		isTx:  false,
	}
}

func (r SqlxRepo) WithTxConn(txConn DbConn) SqlxRepo {
	return SqlxRepo{
		db:    txConn,
		txBeg: r.txBeg,
		isTx:  true,
	}
}

func (r SqlxRepo) DB() DbConn {
	return r.db
}

func (r SqlxRepo) TxBeginner() TxBeginner {
	return r.txBeg
}

func (r SqlxRepo) IsTx() bool {
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
