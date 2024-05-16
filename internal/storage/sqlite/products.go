package sqlite

import (
	"ForeignKey/internal/storage"
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
	    FOREIGN KEY (website_id) REFERENCES websites (id),
	    FOREIGN KEY (image_id) REFERENCES images (id)
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

func (s *Storage) GetProducts(websiteId int) ([]storage.ProductInfo, error) {
	const op = "storage.sqlite.GetProducts"

	q := `SELECT id, name, description, price, image_id  FROM products WHERE website_id=?`

	rows, err := s.db.Query(q, websiteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.ProductInfo, 0)

	var name, description string
	var id, price, imageId int
	for rows.Next() {
		if err = rows.Scan(&id, &name, &description, &price, &imageId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		p := storage.ProductInfo{
			Id:          id,
			WebsiteId:   websiteId,
			Name:        name,
			Description: description,
			Price:       price,
			ImageId:     imageId,
		}
		res = append(res, p)
	}

	return res, nil
}

func (s *Storage) GetProduct(productId int) (*storage.ProductInfo, error) {
	const op = "storage.sqlite.GetProduct"

	q := `SELECT id, website_id, name, description, price, image_id  FROM products WHERE id=?`

	row := s.db.QueryRow(q, productId)

	var name, description string
	var id, websiteId, price, imageId int
	if err := row.Scan(&id, &websiteId, &name, &description, &price, &imageId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p := storage.ProductInfo{
		Id:          id,
		WebsiteId:   websiteId,
		Name:        name,
		Description: description,
		Price:       price,
		ImageId:     imageId,
	}

	return &p, nil
}
