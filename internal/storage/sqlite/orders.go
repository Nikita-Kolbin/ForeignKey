package sqlite

import (
	"ForeignKey/internal/storage"
	"database/sql"
	"fmt"
)

func (s *Storage) initOrders() error {
	q := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
	    FOREIGN KEY (customer_id) REFERENCES customers (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateOrder(customerId int) error {
	const op = "storage.sqlite.CreateOrder"

	q := `INSERT INTO orders (customer_id) VALUES (?)`

	tx, err := s.db.Begin()
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	res, err := tx.Exec(q, customerId)
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	orderId := int(lid)

	cartId, err := s.GetCartId(customerId)
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	cartItems, err := s.GetCartItems(cartId)
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	for _, item := range cartItems {
		err = s.CreateOrderItem(tx, orderId, item.Product.Id, item.Count)
		if err != nil {
			return errWithRollback(tx, op, err)
		}

		err = s.DeleteCartItemWithTx(tx, item.Id)
		if err != nil {
			return errWithRollback(tx, op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errWithRollback(tx, op, err)
	}

	return nil
}

func (s *Storage) GetOrders(customerId int) ([]storage.Order, error) {
	const op = "storage.sqlite.GetOrdersId"

	q := `SELECT id FROM orders WHERE customer_id=?`

	rows, err := s.db.Query(q, customerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersId := make([]int, 0)

	var id int

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ordersId = append(ordersId, id)
	}

	res := make([]storage.Order, 0)

	for _, orderId := range ordersId {
		orderItems, err := s.GetOrderItems(orderId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		o := storage.Order{
			OrderItems: orderItems,
		}

		res = append(res, o)
	}

	return res, nil
}

func errWithRollback(tx *sql.Tx, op string, err error) error {
	_ = tx.Rollback()
	return fmt.Errorf("%s: %w", op, err)
}
