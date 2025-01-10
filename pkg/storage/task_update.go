package storage

import "time"

type TaskUpdate struct {
	Id         int       `json:"id"`
	ActionType int8      `json:"action_type"`
	UserID     int       `json:"user_id"`
	TaskID     int       `json:"task_id"`
	Ts         time.Time `json:"ts"`
}
