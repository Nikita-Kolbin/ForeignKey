package sqlite

import (
	"ForeignKey/internal/storage"
	"crypto/sha1"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"net/mail"
)

func (s *Storage) initAdmins() error {
	q := `
	CREATE TABLE IF NOT EXISTS admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password_hash TEXT
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateAdmin(email, password string) error {
	const op = "storage.sqlite.CreateAdmin"

	q := `INSERT INTO admins (email, password_hash) VALUES (?, ?)`

	if !validEmail(email) {
		return storage.ErrInvalidEmail
	}

	hash := generatePasswordHash(password)

	_, err := s.db.Exec(q, email, hash)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrLoginTaken)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAdminId(email, password string) (int, error) {
	const op = "storage.sqlite.GetAdminId"

	q := `SELECT id FROM admins WHERE email=? AND password_hash=?`

	row := s.db.QueryRow(q, email, generatePasswordHash(password))

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

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
