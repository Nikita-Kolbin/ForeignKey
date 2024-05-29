package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) initWebsites() error {
	q := `
	CREATE TABLE IF NOT EXISTS websites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id INTEGER,
		alias TEXT UNIQUE,
		background_color TEXT,
		font TEXT,
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

	q := `INSERT INTO websites (alias, admin_id, background_color, font) VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(q, alias, adminId, storage.DefaultBackgroundColor, storage.DefaultFont)
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

func (s *Storage) GetWebsite(alias string) (websiteId, adminId int, err error) {
	const op = "storage.sqlite.GetWebsite"

	q := `SELECT id, admin_id FROM websites WHERE alias=?`

	row := s.db.QueryRow(q, alias)

	if err = row.Scan(&websiteId, &adminId); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return websiteId, adminId, nil
}

func (s *Storage) DeleteWebsite(alias string) error {
	const op = "storage.sqlite.DeleteWebsite"

	q := `DELETE FROM websites WHERE alias=?`

	_, err := s.db.Exec(q, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateStyle(alias, backgroundColor, font string) error {
	const op = "storage.sqlite.UpdateStyle"

	q := `UPDATE websites SET background_color = ?, font = ? WHERE alias = ?`

	_, err := s.db.Exec(q, backgroundColor, font, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWebsiteStyle(alias string) (backgroundColor, font string, err error) {
	const op = "storage.sqlite.GetWebsiteStyle"

	q := `SELECT background_color, font FROM websites WHERE alias=?`

	row := s.db.QueryRow(q, alias)

	if err = row.Scan(&backgroundColor, &font); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return backgroundColor, font, nil
}
