package baseRepo

type transactional[R any] interface {
	WithTxConn(db DbConn) R
	TxBeginner() TxBeginner
}

func InTransaction[R any](r transactional[R], fn func(R) error) error {
	// Начинаем транзакцию
	tx, err := r.TxBeginner().Beginx()
	if err != nil {
		return err
	}

	// Создаем транзакционный репозиторий
	txRepo := r.WithTxConn(tx)

	// Откатить в случае паники
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // пробрасываем panic дальше
		}
	}()

	// Выполняем callback
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Коммитим, если не было ошибок
	return tx.Commit()
}
