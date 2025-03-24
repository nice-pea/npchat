package memory

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
	if err = migrate(sqlin.db, config.MigrationsDir); err != nil {
		return nil, err
	}

	return sqlin, nil
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
