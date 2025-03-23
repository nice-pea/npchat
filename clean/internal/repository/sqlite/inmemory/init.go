package inmemory

import (
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	//DSN string
}
type SQLiteInMemory struct {
	db *sqlx.DB
}

func Init(config Config) (*SQLiteInMemory, error) {
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	sqlin := &SQLiteInMemory{
		db: db,
	}
	if err = sqlin.migrate(); err != nil {
		return nil, err
	}

	return sqlin, nil
}

func (m *SQLiteInMemory) migrate() error {
	scriptsFolder := "../../../../migrations/repository/sqlite/inmemory"
	entries, err := os.ReadDir(scriptsFolder)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}
	//scripts := make([]string, len(entries))
	for _, entry := range entries {
		file, err := os.ReadFile(scriptsFolder + string(os.PathSeparator) + entry.Name())
		if err != nil {
			return err
		}
		m.db.MustExec(string(file))
	}

	return nil
}
