package pgsqlRepository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	sqlxRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/sqlx_repo"
)

// Config представляет собой конфигурацию репозитория
type Config struct {
	DSN string
}

// Factory используется для создания репозиториев реализованных с помощью postgresql
type Factory struct {
	db *sqlx.DB
}

func InitFactory(cfg Config) (*Factory, error) {
	conn, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("conn.Ping: %w", err)
	}

	slog.Info("Успешно подключен к PostgreSQL")

	return &Factory{
		db: conn,
	}, nil
}

// Close закрывает соединение с базой данных
func (f *Factory) Close() error {
	return f.db.Close()
}

// Cleanup очищает все сохраненные записи
func (f *Factory) Cleanup() error {
	return sqlxRepo.New(f.db).InTransaction(func(tx sqlxRepo.SqlxRepo) error {
		if _, err := tx.DB().Exec("DELETE FROM sessions"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}
		if _, err := tx.DB().Exec("DELETE FROM oauth_users"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}
		if _, err := tx.DB().Exec("DELETE FROM users"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}
		if _, err := tx.DB().Exec("DELETE FROM participants"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}
		if _, err := tx.DB().Exec("DELETE FROM invitations"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}
		if _, err := tx.DB().Exec("DELETE FROM chats"); err != nil {
			return fmt.Errorf("tx.DB().Exec: %w", err)
		}

		return nil
	})
}

// NewUserrRepository создает репозиторий пользователей
func (f *Factory) NewUserrRepository() userr.Repository {
	return &UserrRepository{
		SqlxRepo: sqlxRepo.New(f.db),
	}
}

// NewChattRepository создает репозиторий чатов
func (f *Factory) NewChattRepository() chatt.Repository {
	return &ChattRepository{
		SqlxRepo: sqlxRepo.New(f.db),
	}
}

// NewSessionnRepository создает репозиторий сессий
func (f *Factory) NewSessionnRepository() sessionn.Repository {
	return &SessionnRepository{
		SqlxRepo: sqlxRepo.New(f.db),
	}
}
