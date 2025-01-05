package postgres

import "fmt"

func (s *Storage) CreateUser(username, hashedPassword, email string) (userID int, err error) {
	const op = "storage.postgres.CreateUser"

	tx, _ := s.db.Begin()

	stmt, err := tx.Prepare("INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: failed to prepare query: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(username, hashedPassword, email).Scan(&userID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return

}

func (s *Storage) GetUser(username string) (userID int, hashedPassword string, err error) {

	const op = "storage.postgres.GetUser"

	row := s.db.QueryRow("SELECT id, password FROM users WHERE username = $1", username)
	if err := row.Scan(&userID, &hashedPassword); err != nil {
		return 0, "", fmt.Errorf("%s: failed to scan user data: %w", op, err)
	}

	return
}
