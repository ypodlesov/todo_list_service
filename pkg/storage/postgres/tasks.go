package postgres

import (
	"database/sql"
	"fmt"
	"todo_list_service/pkg/storage"
)

const UpdateNothing = 0
const UpdateTaskTitleType = 1
const UpdateTaskStatusType = 2
const UpdateTaskTitleStatusType = 3

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

	err = stmt.QueryRow(title, 1, userID).Scan(&task.Id, &task.Title, &task.Status, &task.UserID, &task.CreationTs)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to execute query: %w'`, op, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`'%s: failed to commit transaction: %w'`, op, err)
	}

	return
}

func ConstructUpdateQuery(taskID int, taskTitle string, taskStatus int8) (int, string) {
	hasTitle := len(taskTitle) > 0

	returningStmt := "RETURNING id, title, status, user_id, creation_ts"

	if hasTitle && taskStatus > 0 {
		return UpdateTaskTitleStatusType, fmt.Sprintf(`UPDATE tasks SET title = $1, status = $2 WHERE id = $3 AND user_id = $4 %s`, returningStmt)
	} else if hasTitle {
		return UpdateTaskTitleType, fmt.Sprintf(`UPDATE tasks SET title = $1 WHERE id = $2 AND user_id = $3 %s`, returningStmt)
	} else if taskStatus > 0 {
		return UpdateTaskStatusType, fmt.Sprintf(`UPDATE tasks SET status = $1 WHERE id = $2 AND user_id = $3 %s`, returningStmt)
	} else {
		return UpdateNothing, ``
	}
}

func (s *Storage) UpdateTask(taskID int, taskTitle string, taskStatus int8, userID int) (task *storage.Task, err error) {
	const op = "storage.postgres.UpdateTask"

	task = &storage.Task{}

	updateType, query := ConstructUpdateQuery(taskID, taskTitle, taskStatus)

	if updateType == UpdateNothing {
		return task, nil
	}

	tx, _ := s.db.Begin()

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf(`'%s: failed to prepare query: %w'`, op, err)
	}
	defer stmt.Close()

	var res *sql.Row
	if updateType == UpdateTaskTitleType {
		res = stmt.QueryRow(taskTitle, taskID, userID)
	} else if updateType == UpdateTaskStatusType {
		res = stmt.QueryRow(taskStatus, taskID, userID)
	} else {
		res = stmt.QueryRow(taskTitle, taskStatus, taskID, userID)
	}

	err = res.Scan(&task.Id, &task.Title, &task.Status, &task.UserID, &task.CreationTs)
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

	rows, err := s.db.Query(`SELECT title, status, creation_ts FROM tasks
		WHERE id = $1 AND user_id = $2`, taskID, userID)

	if err != nil {
		return nil, fmt.Errorf(`'%s: failed to get task for user: %w'`, op, err)
	} else {
		defer rows.Close()

		rows.Next()
		task = &storage.Task{
			Id: taskID,
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
