package sqlite

import (
	"database/sql"
	"fmt"
)

func (s *Storage) initCarts() error {
	q := `
	CREATE TABLE IF NOT EXISTS carts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER UNIQUE,
	    FOREIGN KEY (customer_id) REFERENCES customers (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateCart(tx *sql.Tx, customerId int) error {
	const op = "storage.sqlite.CreateCart"

	q := `INSERT INTO carts (customer_id) VALUES (?)`

	_, err := tx.Exec(q, customerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCartId(customerId int) (int, error) {
	const op = "storage.sqlite.GetCartId"

	q := `SELECT id FROM carts WHERE customer_id=?`

	row := s.db.QueryRow(q, customerId)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
