package sqlite

import (
	"ForeignKey/internal/storage"
	"database/sql"
	"fmt"
)

func (s *Storage) initOrderItems() error {
	q := `
	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER,
		saved_product_id INTEGER,
		count INTEGER,
	    FOREIGN KEY (order_id) REFERENCES orders (id),
	    FOREIGN KEY (saved_product_id) REFERENCES saved_products (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateOrderItem(tx *sql.Tx, orderId, productId, count int) error {
	const op = "storage.sqlite.CreateOrderItem"

	p, err := s.GetProduct(productId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	savedProductId, err := s.CreateSavedProduct(p, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	q := `INSERT INTO order_items (order_id, saved_product_id, count) VALUES (?, ?, ?)`

	_, err = tx.Exec(q, orderId, savedProductId, count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrderItems(orderId int) ([]storage.OrderItem, error) {
	const op = "storage.sqlite.GetOrderItems"

	q := `SELECT id, saved_product_id, count FROM order_items WHERE order_id=?`

	rows, err := s.db.Query(q, orderId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.OrderItem, 0)

	var id, productId, count int
	for rows.Next() {
		if err = rows.Scan(&id, &productId, &count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		product, err := s.GetSavedProduct(productId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		p := storage.OrderItem{
			Id:      id,
			OrderId: orderId,
			Product: *product,
			Count:   count,
		}
		res = append(res, p)
	}

	return res, nil
}
