package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
)

func (s *Storage) initCustomers() error {
	// username должен быть уникальным для каждого сайта
	// отдельно, а в этой таблице они могут поторяться.
	q := `
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		email TEXT,
		password_hash TEXT,
	    FOREIGN KEY (website_id) REFERENCES website (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateCustomers(websiteId int, email, password string) error {
	const op = "storage.sqlite.CreateCustomers"

	q := `INSERT INTO customers (website_id, email, password_hash) VALUES (?, ?, ?);`

	if !validEmail(email) {
		return storage.ErrInvalidEmail
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	e, err := s.CustomerIsExists(websiteId, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if e {
		return fmt.Errorf("%s: %w", op, storage.ErrLoginTaken)
	}

	hash := generatePasswordHash(password)

	res, err := tx.Exec(q, websiteId, email, hash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	customerId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = s.CreateCart(tx, int(customerId)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCustomerId(websiteId int, email, password string) (int, error) {
	const op = "storage.sqlite.GetCustomerId"

	q := `SELECT id FROM customers WHERE website_id=? AND email=? AND password_hash=?`

	row := s.db.QueryRow(q, websiteId, email, generatePasswordHash(password))

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetCustomer(id int) (*storage.Customer, error) {
	const op = "storage.sqlite.GetCustomer"

	q := `SELECT website_id, email FROM customers WHERE id=?`

	row := s.db.QueryRow(q, id)

	var websiteId int
	var email string
	if err := row.Scan(&websiteId, &email); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage.Customer{WebsiteId: websiteId, Email: email}, nil
}

func (s *Storage) CustomerIsExists(websiteId int, email string) (bool, error) {
	const op = "storage.sqlite.CustomerIsExists"

	q := `SELECT COUNT(*) FROM customers WHERE website_id=? AND email=?`

	row := s.db.QueryRow(q, websiteId, email)

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return cnt > 0, nil
}
