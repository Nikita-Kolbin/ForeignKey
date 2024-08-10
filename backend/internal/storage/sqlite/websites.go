package sqlite

import (
	"ForeignKey/internal/storage"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) initWebsites() error {
	q := `
	CREATE TABLE IF NOT EXISTS websites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id INTEGER,
		alias TEXT UNIQUE,
		
		background_color TEXT DEFAULT 'white',
		text_color TEXT DEFAULT 'black',
		font TEXT DEFAULT 'Arial',
		
		main_one TEXT DEFAULT '',
		main_two TEXT DEFAULT '',
		
		about_one TEXT DEFAULT '',
		about_two TEXT DEFAULT '',
		about_three TEXT DEFAULT '',
		about_four TEXT DEFAULT '',
		about_five TEXT DEFAULT '',
		about_six TEXT DEFAULT '',
		about_image_one INTEGER DEFAULT 0,
		about_image_two INTEGER DEFAULT 0,
		about_image_three INTEGER DEFAULT 0,
		about_image_four INTEGER DEFAULT 0,
		
		new_product_one TEXT DEFAULT '',
		product_one TEXT DEFAULT '',
		
		contact_one TEXT DEFAULT '',
		contact_two TEXT DEFAULT '',
		contact_three TEXT DEFAULT '',
		contact_four TEXT DEFAULT '',
		contact_five TEXT DEFAULT '',
		
	    FOREIGN KEY (admin_id) REFERENCES admins (id)
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateWebsite(alias string, adminId int) error {
	const op = "storage.sqlite.CreateWebsite"

	a, err := s.GetWebsitesAliases(adminId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(a) > 0 {
		return storage.ErrAdminHaveWebsite
	}

	q := `INSERT INTO websites (alias, admin_id) VALUES (?, ?)`

	_, err = s.db.Exec(q, alias, adminId)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasTaken)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWebsitesAliases(adminId int) ([]string, error) {
	const op = "storage.sqlite.GetWebsitesAliases"

	q := `SELECT alias FROM websites WHERE admin_id=?`

	rows, err := s.db.Query(q, adminId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]string, 0)

	var alias string
	for rows.Next() {
		if err = rows.Scan(&alias); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, alias)
	}

	return res, nil
}

func (s *Storage) GetWebsite(alias string) (websiteId, adminId int, err error) {
	const op = "storage.sqlite.GetWebsite"

	q := `SELECT id, admin_id FROM websites WHERE alias=?`

	row := s.db.QueryRow(q, alias)

	if err = row.Scan(&websiteId, &adminId); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return websiteId, adminId, nil
}

func (s *Storage) GetWebsiteById(id int) (adminId int, alias string, err error) {
	const op = "storage.sqlite.GetWebsiteById"

	q := `SELECT admin_id, alias FROM websites WHERE id=?`

	row := s.db.QueryRow(q, id)

	if err = row.Scan(&adminId, &alias); err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return adminId, alias, nil
}

func (s *Storage) DeleteWebsite(alias string) error {
	const op = "storage.sqlite.DeleteWebsite"

	q := `DELETE FROM websites WHERE alias=?`

	_, err := s.db.Exec(q, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateStyle(alias string, style storage.WebsiteStyle) error {
	const op = "storage.sqlite.UpdateStyle"

	q := `UPDATE websites SET main_one = ?, main_two = ?,
        background_color = ?, text_color = ?, font = ?,
    	about_one = ?, about_two = ?, about_three = ?, about_four = ?, about_five = ?, about_six = ?,
    	about_image_one = ?, about_image_two = ?, about_image_three = ?, about_image_four = ?,
    	new_product_one = ?, product_one = ?,
    	contact_one = ?, contact_two = ?, contact_three = ?, contact_four = ?, contact_five = ?
    WHERE alias = ?`

	_, err := s.db.Exec(
		q,
		style.MainOne, style.MainTwo,
		style.BackgroundColor, style.TextColor, style.Font,
		style.AboutOne, style.AboutTwo, style.AboutThree,
		style.AboutFour, style.AboutFive, style.AboutSix,
		style.AboutImageOne, style.AboutImageTwo,
		style.AboutImageThree, style.AboutImageFour,
		style.NewProductOne, style.ProductOne,
		style.ContactOne, style.ContactTwo, style.ContactThree,
		style.ContactFour, style.ContactFive,
		alias,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWebsiteStyle(alias string) (style *storage.WebsiteStyle, err error) {
	const op = "storage.sqlite.GetWebsiteStyle"

	q := `SELECT background_color, text_color, font, main_one, main_two,
        about_one, about_two, about_three, about_four, about_five, about_six, 
        about_image_one, about_image_two, about_image_three, about_image_four, 
        new_product_one, product_one, 
        contact_one, contact_two, contact_three, contact_four, contact_five 
	FROM websites WHERE alias=?`

	row := s.db.QueryRow(q, alias)

	style = &storage.WebsiteStyle{}
	err = row.Scan(
		&style.BackgroundColor, &style.TextColor, &style.Font,
		&style.MainOne, &style.MainTwo,
		&style.AboutOne, &style.AboutTwo, &style.AboutThree,
		&style.AboutFour, &style.AboutFive, &style.AboutSix,
		&style.AboutImageOne, &style.AboutImageTwo,
		&style.AboutImageThree, &style.AboutImageFour,
		&style.NewProductOne, &style.ProductOne,
		&style.ContactOne, &style.ContactTwo, &style.ContactThree,
		&style.ContactFour, &style.ContactFive,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return style, nil
}
