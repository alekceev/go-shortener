package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/alekceev/go-shortener/app/repos/target"
	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ target.URLStore = &PgStore{}

type PgURL struct {
	ID           uuid.UUID `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	URL          string    `db:"url"`
	ShortURL     string    `db:"short_url"`
	NumRedirects int       `db:"num_redirects"`
}

type PgStore struct {
	db *sql.DB
}

// NewPgStore takes DSN string and trying to ping server.
func NewPgStore(dsn string) (*PgStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	ps := &PgStore{db: db}
	if err = ps.migrate(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *PgStore) migrate() error {
	_, err := s.db.Exec(`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS urls (
		id            uuid NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
		created_at    timestamp with time zone,
		url           varchar,
		short_url     varchar UNIQUE,
		num_redirects bigint default 0
	);`)

	return err
}

func (s *PgStore) Close() error {
	return s.db.Close()
}

func (s *PgStore) Create(ctx context.Context, u target.URL) (*uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO urls (id, created_at, url, short_url) VALUES ($1, $2, $3, $4)`, u.ID, time.Now(), u.URL, u.ShortURL)

	if err != nil {
		return &uuid.UUID{}, err
	}

	return &u.ID, nil
}

func (s *PgStore) Update(ctx context.Context, u target.URL) error {
	pgURL := &PgURL{
		ID:           u.ID,
		ShortURL:     u.ShortURL,
		NumRedirects: u.NumRedirects,
	}
	//TODO update other columns
	_, err := s.db.ExecContext(ctx, "UPDATE urls SET num_redirects = $1 WHERE id = $2", pgURL.NumRedirects, pgURL.ID)
	return err
}

func (s *PgStore) GetURL(ctx context.Context, shortURL string) (*target.URL, error) {
	pgURL := &PgURL{}

	row := s.db.QueryRowContext(ctx, `SELECT id, url, short_url, num_redirects
		FROM urls WHERE short_url = $1`, shortURL)
	err := row.Scan(&pgURL.ID, &pgURL.URL, &pgURL.ShortURL, &pgURL.NumRedirects)
	if err != nil {
		return nil, err
	}
	return &target.URL{
		ID:           pgURL.ID,
		URL:          pgURL.URL,
		ShortURL:     pgURL.ShortURL,
		NumRedirects: pgURL.NumRedirects,
	}, nil
}
