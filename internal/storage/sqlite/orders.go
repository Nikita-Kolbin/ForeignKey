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
		date_time TEXT,
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

	q := `INSERT INTO orders (customer_id, date_time) VALUES (?, datetime('now'))`

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

	if len(cartItems) == 0 {
		return errWithRollback(tx, op, storage.ErrEmptyOrder)
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

	q := `SELECT id, date_time FROM orders WHERE customer_id=?`

	rows, err := s.db.Query(q, customerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersId := make([]int, 0)
	dateTimes := make([]string, 0)

	var id int
	var dateTime string

	for rows.Next() {
		if err = rows.Scan(&id, &dateTime); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ordersId = append(ordersId, id)
		dateTimes = append(dateTimes, dateTime)
	}

	res := make([]storage.Order, 0)

	for i, orderId := range ordersId {
		orderItems, err := s.GetOrderItems(orderId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		o := storage.Order{
			DateTime:   dateTimes[i],
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
