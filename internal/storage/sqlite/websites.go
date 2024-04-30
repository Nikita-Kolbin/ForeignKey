package sqlite

func (s *Storage) initWebsites() error {
	q := `
	CREATE TABLE IF NOT EXISTS websites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id INTEGER UNIQUE,
		alias TEXT,
	    FOREIGN KEY (admin_id) REFERENCES admins (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
