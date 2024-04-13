package sqlite

import (
	"database/sql"
	"fmt"
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

	s := &Storage{db: db}

	err = s.initTables()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (s *Storage) initTables() error {
	q := `
	CREATE TABLE IF NOT EXISTS test_table (
		id INTEGER PRIMARY KEY
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
