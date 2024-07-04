package sqlite

import (
	"ForeignKey/internal/storage"
	"database/sql"
	"fmt"
)

func (s *Storage) initCartItems() error {
	q := `
	CREATE TABLE IF NOT EXISTS cart_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cart_id INTEGER,
		product_id INTEGER,
		count INTEGER,
	    FOREIGN KEY (cart_id) REFERENCES carts (id),
	    FOREIGN KEY (product_id) REFERENCES products (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateCartItem(cartId, productId, count int) error {
	const op = "storage.sqlite.CreateCartItem"

	inCart, err := s.ItemAlreadyInCart(cartId, productId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if inCart {
		if s.AddCartItemCount(cartId, productId, count) != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	q := `INSERT INTO cart_items (cart_id, product_id, count) VALUES (?, ?, ?)`

	_, err = s.db.Exec(q, cartId, productId, count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ItemAlreadyInCart(cartId, productId int) (bool, error) {
	const op = "storage.sqlite.ItemAlreadyInCart"

	q := `SELECT COUNT(*) FROM cart_items WHERE cart_id=? AND product_id=?`

	row := s.db.QueryRow(q, cartId, productId)

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return cnt > 0, nil
}

func (s *Storage) AddCartItemCount(cartId, productId, count int) error {
	const op = "storage.sqlite.AddCartItemCount"

	q := `UPDATE cart_items SET count = count + ? WHERE cart_id=? AND product_id=?`

	_, err := s.db.Exec(q, count, cartId, productId)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ChangeCartItemCount(cartId, productId, newCount int) error {
	const op = "storage.sqlite.AddCartItemCount"

	q := `UPDATE cart_items SET count = ? WHERE cart_id=? AND product_id=?`

	_, err := s.db.Exec(q, newCount, cartId, productId)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCartItems(cartId int) ([]storage.CartItem, error) {
	const op = "storage.sqlite.GetCartItems"

	q := `SELECT id, product_id, count FROM cart_items WHERE cart_id=?`

	rows, err := s.db.Query(q, cartId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.CartItem, 0)

	var id, productId, count int
	for rows.Next() {
		if err = rows.Scan(&id, &productId, &count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		product, err := s.GetProduct(productId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		p := storage.CartItem{
			Id:      id,
			CartId:  cartId,
			Product: *product,
			Count:   count,
		}
		res = append(res, p)
	}

	return res, nil
}

func (s *Storage) DeleteCartItem(id int) error {
	const op = "storage.sqlite.DeleteCartItem"

	q := `DELETE FROM cart_items WHERE id=?`

	_, err := s.db.Exec(q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteCartItemWithTx(tx *sql.Tx, id int) error {
	const op = "storage.sqlite.DeleteCartItemWithTx"

	q := `DELETE FROM cart_items WHERE id=?`

	_, err := tx.Exec(q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
