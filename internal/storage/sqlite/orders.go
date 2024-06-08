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
		status INTEGER,
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

	q := `INSERT INTO orders (customer_id, date_time, status) VALUES (?, datetime('now'), 0)`

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

func (s *Storage) SetOrderStatus(orderId int, status int) error {
	const op = "storage.sqlite.SetOrderStatus"

	q := `UPDATE orders SET status=? WHERE id=?`

	_, err := s.db.Exec(q, status, orderId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrders(customerId int) ([]storage.Order, error) {
	const op = "storage.sqlite.GetOrdersId"

	q := `SELECT id, date_time, status FROM orders WHERE customer_id=?`

	rows, err := s.db.Query(q, customerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersId := make([]int, 0)
	dateTimes := make([]string, 0)
	statuses := make([]int, 0)

	var id, status int
	var dateTime string

	for rows.Next() {
		if err = rows.Scan(&id, &dateTime, &status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ordersId = append(ordersId, id)
		dateTimes = append(dateTimes, dateTime)
		statuses = append(statuses, status)
	}

	res := make([]storage.Order, 0)

	for i, orderId := range ordersId {
		orderItems, err := s.GetOrderItems(orderId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		o := storage.Order{
			Id:         orderId,
			CustomerId: customerId,
			DateTime:   dateTimes[i],
			OrderItems: orderItems,
			Status:     statuses[i],
		}

		res = append(res, o)
	}

	return res, nil
}

func (s *Storage) GetOrderById(id int) (*storage.Order, error) {
	const op = "storage.sqlite.GetOrderById"

	q := `SELECT customer_id, date_time, status FROM orders WHERE id=?;`

	row := s.db.QueryRow(q, id)

	var customerId, status int
	var dateTime string

	if err := row.Scan(&customerId, &dateTime, &status); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	orderItems, err := s.GetOrderItems(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	order := storage.Order{
		Id:         id,
		CustomerId: customerId,
		DateTime:   dateTime,
		OrderItems: orderItems,
		Status:     status,
	}

	return &order, nil
}

func (s *Storage) GetOrdersByWebsite(websiteId int) ([]storage.Order, error) {
	const op = "storage.sqlite.GetOrdersByWebsite"

	q := `
		SELECT id, customer_id, date_time, status
		FROM orders 
		WHERE (SELECT website_id FROM customers WHERE customers.id=orders.customer_id)=?;
	`

	rows, err := s.db.Query(q, websiteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.Order, 0)

	var id, customerId, status int
	var dateTime string

	for rows.Next() {
		if err = rows.Scan(&id, &customerId, &dateTime, &status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		orderItems, err := s.GetOrderItems(id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		order := storage.Order{
			Id:         id,
			CustomerId: customerId,
			DateTime:   dateTime,
			OrderItems: orderItems,
			Status:     status,
		}

		res = append(res, order)
	}

	return res, nil
}

func errWithRollback(tx *sql.Tx, op string, err error) error {
	_ = tx.Rollback()
	return fmt.Errorf("%s: %w", op, err)
}
