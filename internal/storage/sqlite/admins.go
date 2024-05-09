package sqlite

import (
	"ForeignKey/internal/storage"
	"crypto/sha1"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) initAdmins() error {
	q := `
	CREATE TABLE IF NOT EXISTS admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE,
		password_hash TEXT
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateAdmin(login, password string) error {
	const op = "storage.sqlite.CreateAdmin"

	q := `INSERT INTO admins (login, password_hash) VALUES (?, ?)`

	hash := generatePasswordHash(password)

	_, err := s.db.Exec(q, login, hash)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrLoginTaken)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAdminId(login, password string) (int, error) {
	const op = "storage.sqlite.GetAdminId"

	q := `SELECT id FROM admins WHERE login=? AND password_hash=?`

	row := s.db.QueryRow(q, login, generatePasswordHash(password))

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func generatePasswordHash(p string) string {
	hash := sha1.New()
	hash.Write([]byte(p))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
