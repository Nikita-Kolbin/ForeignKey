package sqlite

func (s *Storage) initFooterStyles() error {
	q := `
	CREATE TABLE IF NOT EXISTS footer_styles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER UNIQUE,
		background_color TEXT,
		content TEXT,
		font_size TEXT,
		height TEXT,
		width TEXT,
	    FOREIGN KEY (website_id) REFERENCES websites (id) ON DELETE CASCADE
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
