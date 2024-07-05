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
		telegram TEXT DEFAULT '',
		image_id INTEGER DEFAULT 0,
		
		email_notification INTEGER DEFAULT 0,
		telegram_notification INTEGER DEFAULT 0
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

	q := `SELECT email, first_name, last_name, father_name, city, 
          telegram, image_id, email_notification, telegram_notification
		  FROM admins WHERE id=?`

	row := s.db.QueryRow(q, id)

	var ii, en, tgn int
	var e, fn, ln, fan, tg, c string
	if err := row.Scan(&e, &fn, &ln, &fan, &c, &tg, &ii, &en, &tgn); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage.Admin{
		Id:                   id,
		Email:                e,
		FirstName:            fn,
		LastName:             ln,
		FatherName:           fan,
		City:                 c,
		Telegram:             tg,
		ImageId:              ii,
		EmailNotification:    en,
		TelegramNotification: tgn,
	}, nil
}

func (s *Storage) UpdateAdminProfile(fin, ln, fan, city, tg string, id, ii int) error {
	const op = "storage.sqlite.UpdateAdminProfile"

	q := `UPDATE admins
		  SET first_name=?, last_name=?, father_name=?, city=?, telegram=?, image_id=?
    	  WHERE id=?`

	_, err := s.db.Exec(q, fin, ln, fan, city, tg, ii, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SetAdminEmailNotification(adminId, notification int) error {
	const op = "storage.sqlite.SetAdminEmailNotification"

	if notification != 0 && notification != 1 {
		return storage.ErrInvalidNotification
	}

	q := `UPDATE admins SET email_notification=? WHERE id=?`

	_, err := s.db.Exec(q, notification, adminId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SetAdminTelegramNotification(adminId, notification int) error {
	const op = "storage.sqlite.SetAdminTelegramNotification"

	if notification != 0 && notification != 1 {
		return storage.ErrInvalidNotification
	}

	q := `UPDATE admins SET telegram_notification=? WHERE id=?`

	_, err := s.db.Exec(q, notification, adminId)
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
