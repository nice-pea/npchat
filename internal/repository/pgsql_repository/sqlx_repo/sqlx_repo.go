package sqlxRepo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SqlxRepo struct {
	db         DB
	txBeginner TxBeginner
	isTx       bool // Признак выполнения в транзакции
}

func New(sqlxDB interface {
	DB
	TxBeginner
}) SqlxRepo {
	return SqlxRepo{
		db:         sqlxDB,
		txBeginner: sqlxDB,
		isTx:       false,
	}
}

func (r SqlxRepo) withTx(tx DB) SqlxRepo {
	return SqlxRepo{
		db:         tx,
		txBeginner: r.txBeginner,
		isTx:       true,
	}
}

func (r SqlxRepo) DB() DB {
	return r.db
}

func (r SqlxRepo) IsTx() bool {
	return r.isTx
}

func (r SqlxRepo) InTransaction(f func(tx SqlxRepo) error) error {
	// Начинаем транзакцию
	tx, err := r.txBeginner.Beginx()
	if err != nil {
		return err
	}

	// Создаем транзакционный репозиторий
	txRepo := r.withTx(tx)

	// Откатить в случае паники
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // пробрасываем panic дальше
		}
	}()

	// Выполняем callback
	if err := f(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Коммитим, если не было ошибок
	return tx.Commit()
}

type DB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

type TxBeginner interface {
	Beginx() (*sqlx.Tx, error)
}
