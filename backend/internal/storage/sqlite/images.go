package sqlite

import (
	"fmt"
)

func (s *Storage) initImages() error {
	q := `
	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateImage(path string) (int, error) {
	const op = "storage.sqlite.CreateImage"

	q := `INSERT INTO images (path) VALUES (?)`

	e, err := s.db.Exec(q, path)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := e.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

func (s *Storage) GetImagePath(id int) (string, error) {
	const op = "storage.sqlite.GetImagePath"

	q := `SELECT path FROM images WHERE id=?`

	row := s.db.QueryRow(q, id)

	var path string

	if err := row.Scan(&path); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return path, nil
}
