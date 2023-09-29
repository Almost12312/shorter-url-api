package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"short-url-api/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) (id int64, err error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// Todo: ref
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed get last id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url from url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return res, nil
}

func (s *Storage) DeleteUrl(alias string) (bool, error) {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url where alias = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	r, err := stmt.Exec(alias)
	if err != nil {
		return false, fmt.Errorf("%s: cannot delete url: %w", op, err)
	}

	res, err := r.RowsAffected()

	return res > 0, nil
}
