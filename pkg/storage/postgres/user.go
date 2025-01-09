package postgres

import (
	"database/sql"
	"fmt"
)

func (s *Storage) CreateUser(username, hashedPassword, email string) (userID int, err error) {
	const op = "storage.postgres.CreateUser"

	var cnt int
	row := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username)
	if err := row.Scan(&cnt); err != nil && err != sql.ErrNoRows {
		return -1, fmt.Errorf("%s: failed to scan user data: %w", op, err)
	}
	if cnt != 0 {
		return 0, fmt.Errorf("%s: user with name [%s] already exists", op, username)
	}

	tx, _ := s.db.Begin()

	stmt, err := tx.Prepare("INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		_ = tx.Rollback()
		return -1, fmt.Errorf("%s: failed to prepare query: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(username, hashedPassword, email).Scan(&userID)
	if err != nil {
		_ = tx.Rollback()
		return -1, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return
}

func (s *Storage) GetUser(username string) (userID int, hashedPassword string, err error) {
	const op = "storage.postgres.GetUser"

	row := s.db.QueryRow("SELECT id, password FROM users WHERE username = $1", username)
	if err := row.Scan(&userID, &hashedPassword); err != nil {
		return -1, "", fmt.Errorf("%s: failed to scan user data: %w", op, err)
	}

	return
}
