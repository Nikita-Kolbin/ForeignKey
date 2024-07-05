package sqlite

import (
	"ForeignKey/internal/storage"
	"database/sql"
	"fmt"
)

func (s *Storage) initSavedProducts() error {
	q := `
	CREATE TABLE IF NOT EXISTS seved_products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		name TEXT,
		description TEXT,
		price INTEGER,
		images_id TEXT,
		active INTEGER,
		tags TEXT DEFAULT '',
	    FOREIGN KEY (website_id) REFERENCES websites (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateSavedProduct(p *storage.ProductInfo, tx *sql.Tx) (int, error) {
	const op = "storage.sqlite.CreateSavedProduct"

	q := `INSERT INTO seved_products (website_id, name, description, price, images_id, active, tags) 
		  VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := tx.Exec(q, p.WebsiteId, p.Name, p.Description, p.Price, p.ImagesId, p.Active, p.Tags)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

func (s *Storage) GetSavedProduct(productId int) (*storage.ProductInfo, error) {
	const op = "storage.sqlite.GetSavedProduct"

	q := `SELECT id, website_id, name, description, price, images_id, active, tags 
		  FROM seved_products WHERE id=?`

	row := s.db.QueryRow(q, productId)

	var name, description, imagesId, tags string
	var id, websiteId, price, active int
	if err := row.Scan(&id, &websiteId, &name, &description, &price, &imagesId, &active, &tags); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p := storage.ProductInfo{
		Id:          id,
		WebsiteId:   websiteId,
		Name:        name,
		Description: description,
		Price:       price,
		ImagesId:    imagesId,
		Active:      active,
		Tags:        tags,
	}

	return &p, nil
}
