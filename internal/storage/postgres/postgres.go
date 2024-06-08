package postgres

import (
	"database/sql"
	"fmt"
	"time"

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
			id INTEGER PRIMARY KEY,
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
