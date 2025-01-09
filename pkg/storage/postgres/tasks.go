package postgres

import (
	"fmt"
	"todo_list_service/pkg/storage"
)

func (s *Storage) GetTasks(userID int) (tasks []storage.Task, err error) {
	const op = "storage.postgres.GetTasks"

	rows, err := s.db.Query("SELECT id, title, status, creation_ts FROM tasks WHERE user_id = $1", userID)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get tasks for user: %w", op, err)
	} else {
		defer rows.Close()

		for rows.Next() {
			task := storage.Task{
				UserId: userID,
			}
			if err := rows.Scan(&task.Id, &task.Title, &task.Status, &task.CreationTs); err != nil {
				return nil, fmt.Errorf("%s: failed to read task: %w", op, err)
			}
			tasks = append(tasks, task)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: failed to get tasks for user: %w", op, err)
		}
		return
	}
}
