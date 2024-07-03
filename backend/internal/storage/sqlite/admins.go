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
		password_hash TEXT,
		
		first_name TEXT DEFAULT '',
		last_name TEXT DEFAULT '',
		father_name TEXT DEFAULT '',
		city TEXT DEFAULT '',
		image_id INTEGER DEFAULT 0
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateAdmin(email, password string) (int, error) {
	const op = "storage.sqlite.CreateAdmin"

	q := `INSERT INTO admins (email, password_hash) VALUES (?, ?)`

	if !validEmail(email) {
		return 0, storage.ErrInvalidEmail
	}

	hash := generatePasswordHash(password)

	r, err := s.db.Exec(q, email, hash)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrEmailRegistered)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return int(id), nil
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

func (s *Storage) GetAdminById(id int) (*storage.Admin, error) {
	const op = "storage.sqlite.GetAdminById"

	q := `SELECT email, first_name, last_name, father_name, city, image_id
		  FROM admins WHERE id=?`

	row := s.db.QueryRow(q, id)

	var ii int
	var e, fn, ln, fan, c string
	if err := row.Scan(&e, &fn, &ln, &fan, &c, &ii); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage.Admin{
		Id: id, Email: e, FirstName: fn, LastName: ln, FatherName: fan, City: c, ImageId: ii,
	}, nil
}

func (s *Storage) UpdateAdminProfile(fin, ln, fan, city string, id, ii int) error {
	const op = "storage.sqlite.UpdateAdminProfile"

	q := `UPDATE admins
		  SET first_name=?, last_name=?, father_name=?, city=?, image_id=?
    	  WHERE id=?`

	_, err := s.db.Exec(q, fin, ln, fan, city, ii, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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
