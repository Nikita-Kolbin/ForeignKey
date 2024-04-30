package sqlite

// TODO: Add image id foreign key
func (s *Storage) initProducts() error {
	q := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER,
		name TEXT,
		price INTEGER,
		image_id INTEGER,
	    FOREIGN KEY (website_id) REFERENCES website (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
