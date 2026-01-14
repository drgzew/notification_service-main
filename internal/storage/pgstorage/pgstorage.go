package pgstorage

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (*PGStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфигурации")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения к БД")
	}

	storage := &PGStorage{db: db}

	// Инициализация таблиц
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGStorage) initTables() error {
	sql := `
	CREATE TABLE IF NOT EXISTS notification (
		id TEXT PRIMARY KEY,
		recipient TEXT NOT NULL,
		message TEXT NOT NULL,
		created_at TIMESTAMPTZ
	);

	CREATE TABLE IF NOT EXISTS notification_status (
		notification_id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		error TEXT,
		sent_at TIMESTAMPTZ
	);`
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "initTables")
	}
	return nil
}