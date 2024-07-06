package sqlite

func (s *Storage) initUsers() error {
	q := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		chat_id INTEGER
	);
	`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateUser(userId string, chatId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	exists, err := s.userIsExists(userId)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	q := `INSERT INTO users (username, chat_id) VALUES (?, ?)`

	_, err = s.db.Exec(q, userId, chatId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetChatIdByUsername(username string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	q := `SELECT chat_id FROM users WHERE username=?`

	row := s.db.QueryRow(q, username)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) userIsExists(username string) (bool, error) {
	q := `SELECT COUNT(*) FROM users WHERE username=?`

	row := s.db.QueryRow(q, username)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return false, err
	}

	return cnt > 0, nil
}
