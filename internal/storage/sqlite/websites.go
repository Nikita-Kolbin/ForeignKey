package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

// TODO: сделать ограниченяе на один сайт для админа

func (s *Storage) initWebsites() error {
	q := `
	CREATE TABLE IF NOT EXISTS websites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id INTEGER,
		alias TEXT UNIQUE,
	    FOREIGN KEY (admin_id) REFERENCES admins (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateWebsite(alias string, adminId int) error {
	const op = "storage.sqlite.CreateWebsite"

	q := `INSERT INTO websites (alias, admin_id) VALUES (?, ?)`

	_, err := s.db.Exec(q, alias, adminId)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasTaken)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWebsitesAliases(adminId int) ([]string, error) {
	const op = "storage.sqlite.GetWebsitesAliases"

	q := `SELECT alias FROM websites WHERE admin_id=?`

	rows, err := s.db.Query(q, adminId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]string, 0)

	var alias string
	for rows.Next() {
		if err = rows.Scan(&alias); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, alias)
	}

	return res, nil
}
