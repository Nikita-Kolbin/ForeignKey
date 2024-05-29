package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
)

func (s *Storage) initAboutUsStyles() error {
	q := `
	CREATE TABLE IF NOT EXISTS about_us_styles (
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

func (s *Storage) CreateAboutUsStyles(websiteId int, aboutUs storage.AboutUsStyle) error {
	const op = "storage.sqlite.CreateAboutUsStyles"

	q := `INSERT INTO about_us_styles (website_id, background_color, content, font_size, height, width) 
		  VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(q, websiteId,
		aboutUs.BackgroundColor, aboutUs.Content, aboutUs.FontSize, aboutUs.Height, aboutUs.Width,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
