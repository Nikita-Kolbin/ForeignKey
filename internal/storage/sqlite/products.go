package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
	"strconv"
	"strings"
)

func (s *Storage) initProducts() error {
	q := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		name TEXT,
		description TEXT,
		price INTEGER,
		images_id TEXT,
	    FOREIGN KEY (website_id) REFERENCES websites (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateProduct(name, description, imagesId string, websiteId, price int) error {
	const op = "storage.sqlite.CreateProduct"

	if err := validateImagesId(imagesId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	q := `INSERT INTO products (website_id, name, description, price, images_id) VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(q, websiteId, name, description, price, imagesId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetProducts(websiteId int) ([]storage.ProductInfo, error) {
	const op = "storage.sqlite.GetProducts"

	q := `SELECT id, name, description, price, images_id  FROM products WHERE website_id=?`

	rows, err := s.db.Query(q, websiteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.ProductInfo, 0)

	var name, description, imagesId string
	var id, price int
	for rows.Next() {
		if err = rows.Scan(&id, &name, &description, &price, &imagesId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		p := storage.ProductInfo{
			Id:          id,
			WebsiteId:   websiteId,
			Name:        name,
			Description: description,
			Price:       price,
			ImagesId:    imagesId,
		}
		res = append(res, p)
	}

	return res, nil
}

func (s *Storage) GetProduct(productId int) (*storage.ProductInfo, error) {
	const op = "storage.sqlite.GetProduct"

	q := `SELECT id, website_id, name, description, price, images_id  FROM products WHERE id=?`

	row := s.db.QueryRow(q, productId)

	var name, description, imagesId string
	var id, websiteId, price int
	if err := row.Scan(&id, &websiteId, &name, &description, &price, &imagesId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p := storage.ProductInfo{
		Id:          id,
		WebsiteId:   websiteId,
		Name:        name,
		Description: description,
		Price:       price,
		ImagesId:    imagesId,
	}

	return &p, nil
}

func validateImagesId(ids string) error {
	if len(ids) == 0 {
		return nil
	}

	s := strings.Split(ids, " ")
	for _, id := range s {
		if _, err := strconv.Atoi(id); err != nil {
			return storage.ErrInvalidImagesIs
		}
	}

	return nil
}
