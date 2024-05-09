package sqlite

import (
	"fmt"
)

func (s *Storage) initProducts() error {
	q := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		name TEXT,
		description TEXT,
		price INTEGER,
		image_id INTEGER,
	    FOREIGN KEY (website_id) REFERENCES website (id),
	    FOREIGN KEY (image_id) REFERENCES image (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateProduct(name, description string, websiteId, price, imageId int) error {
	const op = "storage.sqlite.CreateProduct"

	q := `INSERT INTO products (website_id, name, description, price, image_id) VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(q, websiteId, name, description, price, imageId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
