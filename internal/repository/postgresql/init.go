package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var configPsql = Config{
	DSN: "",
}

type Config struct {
	DSN string
}

type RepositoryFactory struct {
	db *sqlx.DB
}

func InitRepositoryFactory(config Config) (*RepositoryFactory, error) {
	db, err := sqlx.Connect("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlx.DB.Ping: %w", err)
	}

	pqMemory := &RepositoryFactory{
		db: db,
	}

	return pqMemory, nil
}

func (m *RepositoryFactory) Close() error {
	return m.db.Close()
}
