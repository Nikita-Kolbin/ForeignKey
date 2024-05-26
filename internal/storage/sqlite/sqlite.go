package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{db: db}

	err = s.initTables()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (s *Storage) initTables() error {
	if err := s.initAdmins(); err != nil {
		return fmt.Errorf("can't init admins: %w", err)
	}

	if err := s.initWebsites(); err != nil {
		return fmt.Errorf("can't init websites: %w", err)
	}

	if err := s.initCustomers(); err != nil {
		return fmt.Errorf("can't init customers: %w", err)
	}

	if err := s.initImages(); err != nil {
		return fmt.Errorf("can't init images: %w", err)
	}

	if err := s.initProducts(); err != nil {
		return fmt.Errorf("can't init products: %w", err)
	}

	if err := s.initCarts(); err != nil {
		return fmt.Errorf("can't init carts: %w", err)
	}

	if err := s.initCartItems(); err != nil {
		return fmt.Errorf("can't init cart_items: %w", err)
	}

	if err := s.initOrders(); err != nil {
		return fmt.Errorf("can't init orders: %w", err)
	}

	if err := s.initOrderItems(); err != nil {
		return fmt.Errorf("can't init order_items: %w", err)
	}

	if err := s.initAboutUsStyles(); err != nil {
		return fmt.Errorf("can't init about_us_styles: %w", err)
	}

	if err := s.initFooterStyles(); err != nil {
		return fmt.Errorf("can't init footer_styles: %w", err)
	}

	if err := s.initHeaderStyles(); err != nil {
		return fmt.Errorf("can't init header_styles: %w", err)
	}

	if err := s.initProductsStyles(); err != nil {
		return fmt.Errorf("can't init products_styles: %w", err)
	}

	return nil
}
