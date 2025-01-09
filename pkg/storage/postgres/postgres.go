package postgres

import (
	"database/sql"
	"fmt"
	"todo_list_service/pkg/config"
)

type Storage struct {
	cfg *config.PgConfig
	db  *sql.DB
}

func generateUrlFromConfig(cfg *config.PgConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
}

func (s *Storage) Close() error {
	err := s.db.Close()
	return err
}

func New(cfg *config.PgConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", generateUrlFromConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{
		db:  db,
		cfg: cfg,
	}

	if err := storage.applyMigrations(cfg.MigrationsDir); err != nil {
		return nil, err
	}

	return storage, nil
}
