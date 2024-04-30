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
		username TEXT UNIQUE,
		password_hash TEXT
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateAdmin(username, password string) error {
	const op = "storage.sqlite.CreateAdmin"

	q := `INSERT INTO admins (username, password_hash) VALUES (?, ?)`

	hash := generatePasswordHash(password)

	_, err := s.db.Exec(q, username, hash)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrUsernameTaken)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAdminId(username, password string) (int, error) {
	const op = "storage.sqlite.AdminIsExists"

	q := `SELECT id FROM admins WHERE username=? AND password_hash=?`

	row := s.db.QueryRow(q, username, generatePasswordHash(password))

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
