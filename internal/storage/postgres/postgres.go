package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/storage"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(conn string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := OpenDB(conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt2, err := db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func OpenDB(conn string) (*sql.DB, error) {
	const op = "storage.postgres.OpenDB"
	count := 0

	for {
		count++
		db, err := sql.Open("pgx", conn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		if count > 8 {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		time.Sleep(time.Second * 2)
	}
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const op = "storage.postgres.SaveURL"
	query := "INSERT INTO URL(url, alias) VALUES($1, $2)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetRUL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	query := "SELECT url FROM url WHERE alias = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(alias)
	var url string
	err = row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	query := "DELETE FROM url WHERE alias = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
