package postgres

import (
	"fmt"
	"todo_list_service/pkg/storage"
)

func (s *Storage) CreateTask(title string, userID int) (task *storage.Task, err error) {
	const op = "storage.postgres.CreateTask"

	tx, _ := s.db.Begin()

	stmt, err := tx.Prepare(`INSERT INTO tasks (title, status, user_id) VALUES ($1, $2, $3)
							 RETURNING id, title, status, user_id, creation_ts`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer stmt.Close()

	task = &storage.Task{}

	err = stmt.QueryRow(title, 0, userID).Scan(&task.Id, &task.Title, &task.Status, &task.UserID, &task.CreationTs)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`'%s: failed to commit transaction: %w'`, op, err)
	}

	return
}

func (s *Storage) GetTask(taskID, userID int) (task *storage.Task, err error) {
	const op = "storage.postgres.GetTasks"

	rows, err := s.db.Query(`SELECT title, status, user_id, creation_ts FROM tasks
		WHERE id = $1 AND user_id = $2`, taskID, userID)

	if err != nil {
		return nil, fmt.Errorf(`'%s: failed to get task for user: %w'`, op, err)
	} else {
		defer rows.Close()

		rows.Next()
		task = &storage.Task{
			Id: taskID,
		}
		if err := rows.Scan(&task.Id, &task.Title, &task.Status, &task.CreationTs); err != nil {
			return nil, fmt.Errorf(`'%s: failed to read task: %w'`, op, err)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf(`'%s: got error while selecting task from db: %w'`, op, err)
		}
		return
	}
}

func (s *Storage) GetTasks(userID int) (tasks []storage.Task, err error) {
	const op = "storage.postgres.GetTasks"

	tasks = []storage.Task{}

	rows, err := s.db.Query("SELECT id, title, status, creation_ts FROM tasks WHERE user_id = $1", userID)

	if err != nil {
		return nil, fmt.Errorf(`'%s: failed to get tasks for user: %w'`, op, err)
	} else {
		defer rows.Close()

		for rows.Next() {
			task := storage.Task{
				UserID: userID,
			}
			if err := rows.Scan(&task.Id, &task.Title, &task.Status, &task.CreationTs); err != nil {
				return nil, fmt.Errorf(`'%s: failed to read task: %w'`, op, err)
			}
			tasks = append(tasks, task)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf(`'%s: failed to get tasks for user: %w'`, op, err)
		}
		return
	}
}
