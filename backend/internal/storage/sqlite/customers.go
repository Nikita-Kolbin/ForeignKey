package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
)

func (s *Storage) initCustomers() error {
	// username должен быть уникальным для каждого сайта
	// отдельно, а в этой таблице они могут поторяться.
	q := `
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		email TEXT,
		password_hash TEXT,
		
		first_name TEXT DEFAULT '',
		last_name TEXT DEFAULT '',
		father_name TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		telegram TEXT DEFAULT '',
		delivery_type TEXT DEFAULT '',
		payment_type TEXT DEFAULT '',
		
		email_notification INTEGER DEFAULT 0,
		telegram_notification INTEGER DEFAULT 0,
		
	    FOREIGN KEY (website_id) REFERENCES websites (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateCustomers(websiteId int, email, password string) (int, error) {
	const op = "storage.sqlite.CreateCustomers"

	q := `INSERT INTO customers (website_id, email, password_hash) VALUES (?, ?, ?);`

	if !validEmail(email) {
		return 0, storage.ErrInvalidEmail
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	e, err := s.CustomerIsExists(websiteId, email)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if e {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrEmailRegistered)
	}

	hash := generatePasswordHash(password)

	res, err := tx.Exec(q, websiteId, email, hash)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	customerId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.CreateCart(tx, int(customerId)); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(customerId), nil
}

func (s *Storage) GetCustomerId(websiteId int, email, password string) (int, error) {
	const op = "storage.sqlite.GetCustomerId"

	q := `SELECT id FROM customers WHERE website_id=? AND email=? AND password_hash=?`

	row := s.db.QueryRow(q, websiteId, email, generatePasswordHash(password))

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetCustomer(id int) (*storage.Customer, error) {
	const op = "storage.sqlite.GetCustomer"

	q := `SELECT website_id, email, first_name, last_name, father_name, phone, 
		  telegram, delivery_type, payment_type, email_notification, telegram_notification
		  FROM customers WHERE id=?`

	row := s.db.QueryRow(q, id)

	var websiteId, en, tgn int
	var email, fin, ln, fan, ph, tg, dt, pt string
	if err := row.Scan(&websiteId, &email, &fin, &ln, &fan, &ph, &tg, &dt, &pt, &en, &tgn); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	customer := &storage.Customer{
		Id:                   id,
		WebsiteId:            websiteId,
		Email:                email,
		FirstName:            fin,
		LastName:             ln,
		FatherName:           fan,
		Phone:                ph,
		Telegram:             tg,
		DeliveryType:         dt,
		PaymentType:          pt,
		EmailNotification:    en,
		TelegramNotification: tgn,
	}

	return customer, nil
}

func (s *Storage) GetCustomersByWebsite(websiteId int) ([]storage.Customer, error) {
	const op = "storage.sqlite.GetCustomersByWebsite"

	q := `SELECT id, email, first_name, last_name, father_name, phone, 
		  telegram, delivery_type, payment_type,  email_notification, telegram_notification
		  FROM customers WHERE website_id=?;`

	rows, err := s.db.Query(q, websiteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]storage.Customer, 0)

	var id, en, tgn int
	var email, fin, ln, fan, ph, tg, dt, pt string

	for rows.Next() {
		if err = rows.Scan(&id, &email, &fin, &ln, &fan, &ph, &tg, &dt, &pt, &en, &tgn); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		c := storage.Customer{
			Id:                   id,
			WebsiteId:            websiteId,
			Email:                email,
			FirstName:            fin,
			LastName:             ln,
			FatherName:           fan,
			Phone:                ph,
			Telegram:             tg,
			DeliveryType:         dt,
			PaymentType:          pt,
			EmailNotification:    en,
			TelegramNotification: tgn,
		}

		res = append(res, c)
	}

	return res, nil
}

func (s *Storage) UpdateCustomerProfile(fin, ln, fan, ph, tg, dt, pt string, id int) error {
	const op = "storage.sqlite.UpdateAdminProfile"

	q := `UPDATE customers
		  SET first_name=?, last_name=?, father_name=?, phone=?, 
		      telegram=?, delivery_type=?, payment_type=?
    	  WHERE id=?`

	_, err := s.db.Exec(q, fin, ln, fan, ph, tg, dt, pt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CustomerIsExists(websiteId int, email string) (bool, error) {
	const op = "storage.sqlite.CustomerIsExists"

	q := `SELECT COUNT(*) FROM customers WHERE website_id=? AND email=?`

	row := s.db.QueryRow(q, websiteId, email)

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return cnt > 0, nil
}

func (s *Storage) SetCustomerEmailNotification(customerId, notification int) error {
	const op = "storage.sqlite.SetCustomerEmailNotification"

	if notification != 0 && notification != 1 {
		return storage.ErrInvalidNotification
	}

	q := `UPDATE customers SET email_notification=? WHERE id=?`

	_, err := s.db.Exec(q, notification, customerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SetCustomerTelegramNotification(customerId, notification int) error {
	const op = "storage.sqlite.SetCustomerTelegramNotification"

	if notification != 0 && notification != 1 {
		return storage.ErrInvalidNotification
	}

	q := `UPDATE customers SET telegram_notification=? WHERE id=?`

	_, err := s.db.Exec(q, notification, customerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
