package sqlite

func (s *Storage) initCustomers() error {
	// username должен быть уникальным для каждого сайта
	// отдельно, а в этой таблице они могут поторяться.
	q := `
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		username TEXT,
		password_hash TEXT,
	    FOREIGN KEY (website_id) REFERENCES website (id)                 
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
