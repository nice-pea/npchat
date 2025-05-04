package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	MigrationsDir string
	//DSN           string
}
type RepositoryFactory struct {
	db *sqlx.DB
}

func InitRepositoryFactory(config Config) (*RepositoryFactory, error) {
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	sqliteMemory := &RepositoryFactory{
		db: db,
	}
	if err = migrate(sqliteMemory.db, config.MigrationsDir); err != nil {
		return nil, err
	}

	return sqliteMemory, nil
}

func (m *RepositoryFactory) Close() error {
	return m.db.Close()
}

func migrate(db *sqlx.DB, migrationsDir string) error {
	if migrationsDir == "" {
		return nil
	}
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		path := filepath.Join(migrationsDir, entry.Name())
		var file []byte
		if file, err = os.ReadFile(path); err != nil {
			return err
		}
		if _, err = db.Exec(string(file)); err != nil {
			return err
		}
	}

	return nil
}
