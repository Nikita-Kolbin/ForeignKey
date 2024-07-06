package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"sync"
)

type Storage struct {
	mu *sync.Mutex
	db *sql.DB
}

func New(storagePath, storageName string) (*Storage, error) {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		_ = os.MkdirAll(storagePath, os.ModePerm)
	}

	db, err := sql.Open("sqlite3", path.Join(storagePath, storageName))
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	s := &Storage{db: db, mu: &sync.Mutex{}}

	err = s.initUsers()
	if err != nil {
		return nil, fmt.Errorf("can't init users: %w", err)
	}

	return s, nil
}
