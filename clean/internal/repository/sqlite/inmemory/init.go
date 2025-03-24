package inmemory

import (
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	MigrationsDir string
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
	if err = migrate(sqlin.db, config); err != nil {
		return nil, err
	}

	return sqlin, nil
}

func migrate(db *sqlx.DB, cfg Config) error {
	entries, err := os.ReadDir(cfg.MigrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}
	for _, entry := range entries {
		filename := cfg.MigrationsDir + string(os.PathSeparator) + entry.Name()
		var file []byte
		if file, err = os.ReadFile(filename); err != nil {
			return err
		}
		if _, err = db.Exec(string(file)); err != nil {
			return err
		}
	}

	return nil
}
