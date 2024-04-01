package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.Postgres.New"

	db, err := sql.Open("postgres", "user=postgres password=postgres host=host.docker.internal dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `CREATE TABLE IF NOT EXISTS url(id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE ,
		url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS index_alias ON url(alias)`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlSave, alias string) error {
	const op = "storage.Postgres.SaveURL"

	query := fmt.Sprintf("INSERT INTO url(url, alias) VALUES ('%s', '%s')", urlSave, alias)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "Storage.Postgres.GetURL"

	query := fmt.Sprintf("SELECT url.url FROM url WHERE alias='%s'", alias)
	rows, err := s.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var res string
	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", storage.ErrURLNotFound
			}
			return "", fmt.Errorf("%s: execute statement: %w", op, err)
		}
	}
	if res == "" {
		return "", storage.ErrURLNotFound
	}
	return res, nil
}

func (s *Storage) GetAlias(ctx context.Context, url string) (string, error) {
	const op = "Storage.Postgres.GetAlias"

	query := fmt.Sprintf("SELECT alias FROM url WHERE url='%s'", url)
	rows, err := s.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var res string
	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", storage.ErrURLNotFound
			}
			return "", fmt.Errorf("%s: execute statement: %w", op, err)
		}
	}
	if res == "" {
		return "", storage.ErrURLNotFound
	}
	return res, nil
}
