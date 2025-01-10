package postgres

import (
	"fmt"
	"todo_list_service/internal/storage"
)

func (s *Storage) GetMaxPriority(userID int) (int, error) {
	const op = "storage.postgres.GetMaxPriority"

	tasks, err := s.GetTasks(userID, 1)
	if err != nil {
		return 0, fmt.Errorf(`'%s: failed to get max_priority task for user [%d]: %w'`, op, userID, err)
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	return tasks[len(tasks)-1].Priority, nil
}

func (s *Storage) CreateTask(newTask *storage.Task) (task *storage.Task, err error) {
	const op = "storage.postgres.CreateTask"

	maxPriority, err := s.GetMaxPriority(newTask.UserID)
	if err != nil {
		return nil, err
	}

	tx, _ := s.db.Begin()

	createTaskStmt, err := tx.Prepare(`INSERT INTO tasks (title, description, status, priority, user_id) VALUES ($1, $2, $3, $4, $5)
							           RETURNING id, title, description, status, priority, user_id, creation_ts`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer createTaskStmt.Close()

	task = &storage.Task{}

	queryRes := createTaskStmt.QueryRow(newTask.Title, newTask.Description, storage.TaskStatusOpened, maxPriority+storage.TaskPriorityCreatedDelta, newTask.UserID)
	err = queryRes.Scan(&task.ID, &task.Title, &task.Status, &task.UserID, &task.CreationTs)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	insertActionStmt, err := tx.Prepare(`INSERT INTO task_actions (action_type, user_id, task_id) VALUES ($1, $2, $3)`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer insertActionStmt.Close()

	if _, err = insertActionStmt.Exec(storage.CreateTaskType, task.UserID, task.ID); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`'%s: failed to commit transaction: %w'`, op, err)
	}

	return
}

func (s *Storage) UpdateTask(updatedTask *storage.Task) (*storage.Task, error) {
	const op = "storage.postgres.UpdateTask"

	priority := updatedTask.Priority
	if updatedTask.Status == storage.TaskStatusClosed {
		priority = storage.TaskPriorityClosed
	}

	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4 WHERE user_id = $5 AND id = $6
			  RETURNING id, title, description, status, priority, user_id, creation_ts`

	tx, _ := s.db.Begin()

	updateStmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer updateStmt.Close()

	task := &storage.Task{}

	queryRes := updateStmt.QueryRow(updatedTask.Title, updatedTask.Description, updatedTask.Status, priority, updatedTask.UserID, updatedTask.ID)
	err = queryRes.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.UserID, &task.CreationTs)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	insertStmt, err := tx.Prepare(`INSERT INTO task_actions (action_type, user_id, task_id) VALUES ($1, $2, $3)`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(storage.UpdateTaskType, task.UserID, task.ID); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`'%s: failed to commit transaction: %w'`, op, err)
	}

	return task, nil
}

func (s *Storage) GetTask(taskID, userID int) (task *storage.Task, err error) {
	const op = "storage.postgres.GetTasks"

	rows, err := s.db.Query(`SELECT title, status, creation_ts FROM tasks
		WHERE user_id = $1 AND id = $2`, userID, taskID)

	if err != nil {
		return nil, fmt.Errorf(`'%s: failed to get task for user: %w'`, op, err)
	} else {
		defer rows.Close()

		rows.Next()
		task = &storage.Task{
			ID: taskID,
		}
		if err := rows.Scan(&task.Title, &task.Status, &task.CreationTs); err != nil {
			return nil, fmt.Errorf(`'%s: failed to read task: %w'`, op, err)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf(`'%s: got error while selecting task from db: %w'`, op, err)
		}
		return
	}
}

func (s *Storage) GetTasks(userID, limit int) (tasks []storage.Task, err error) {
	const op = "storage.postgres.GetTasks"

	tasks = []storage.Task{}

	rows, err := s.db.Query("SELECT id, title, description, status, priority, creation_ts FROM tasks WHERE user_id = $1 ORDER BY priority DESC LIMIT $2", userID, limit)

	if err != nil {
		return nil, fmt.Errorf(`'%s: failed to get tasks for user [%d]: %w'`, op, userID, err)
	} else {
		defer rows.Close()

		for rows.Next() {
			task := storage.Task{
				UserID: userID,
			}
			if err := rows.Scan(&task.ID, &task.Title, &task.Priority, &task.Status, &task.Priority, &task.CreationTs); err != nil {
				return nil, fmt.Errorf(`'%s: failed to read task: %w'`, op, err)
			}
			tasks = append(tasks, task)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf(`'%s: failed to get tasks for user [%d]: %w'`, op, userID, err)
		}
		return
	}
}
