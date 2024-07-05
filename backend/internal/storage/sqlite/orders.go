package sqlite

import (
	"ForeignKey/internal/storage"
	"database/sql"
	"fmt"
	"sort"
)

func (s *Storage) initOrders() error {
	q := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		date_time TEXT,
		status INTEGER,
		comment TEXT,
	    FOREIGN KEY (customer_id) REFERENCES customers (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateOrder(customerId int, comment string) error {
	const op = "storage.sqlite.CreateOrder"

	q := `INSERT INTO orders (customer_id, date_time, status, comment) VALUES (?, datetime('now'), 0, ?)`

	tx, err := s.db.Begin()
	if err != nil {
		return errWithRollback(tx, op, err)
	}

	res, err := tx.Exec(q, customerId, comment)
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

	q := `SELECT id, date_time, status, comment FROM orders WHERE customer_id=?`

	rows, err := s.db.Query(q, customerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.Order, 0)
	var orderId, status int
	var dateTime, comment string

	for rows.Next() {
		if err = rows.Scan(&orderId, &dateTime, &status, &comment); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		orderItems, err := s.GetOrderItems(orderId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		order := storage.Order{
			Id:         orderId,
			CustomerId: customerId,
			DateTime:   dateTime,
			OrderItems: orderItems,
			Status:     status,
			Comment:    comment,
		}

		res = append(res, order)
	}

	return res, nil
}

func (s *Storage) GetOrderById(id int) (*storage.Order, error) {
	const op = "storage.sqlite.GetOrderById"

	q := `SELECT customer_id, date_time, status, comment FROM orders WHERE id=?;`

	row := s.db.QueryRow(q, id)

	var customerId, status int
	var dateTime, comment string

	if err := row.Scan(&customerId, &dateTime, &status, &comment); err != nil {
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
		Comment:    comment,
	}

	return &order, nil
}

func (s *Storage) GetOrdersByWebsite(websiteId int) ([]storage.Order, error) {
	const op = "storage.sqlite.GetOrdersByWebsite"

	q := `
		SELECT id, customer_id, date_time, status, comment
		FROM orders 
		WHERE (SELECT website_id FROM customers WHERE customers.id=orders.customer_id)=?;
	`

	rows, err := s.db.Query(q, websiteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.Order, 0)

	var id, customerId, status int
	var dateTime, comment string

	for rows.Next() {
		if err = rows.Scan(&id, &customerId, &dateTime, &status, &comment); err != nil {
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
			Comment:    comment,
		}

		res = append(res, order)
	}

	return res, nil
}

func (s *Storage) GetCompletedOrders(websiteId int) ([]storage.Order, error) {
	const op = "storage.sqlite.GetOrdersByWebsite"

	q := `
		SELECT id, customer_id, date_time, status, comment
		FROM orders 
		WHERE (SELECT website_id FROM customers WHERE customers.id=orders.customer_id)=? AND status=?;
	`

	rows, err := s.db.Query(q, websiteId, storage.StatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.Order, 0)

	var id, customerId, status int
	var dateTime, comment string

	for rows.Next() {
		if err = rows.Scan(&id, &customerId, &dateTime, &status, &comment); err != nil {
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
			Comment:    comment,
		}

		res = append(res, order)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Id < res[j].Id
	})

	return res, nil
}

func errWithRollback(tx *sql.Tx, op string, err error) error {
	_ = tx.Rollback()
	return fmt.Errorf("%s: %w", op, err)
}
