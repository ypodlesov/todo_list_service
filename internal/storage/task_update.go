package storage

import "time"

const UpdateNothingType = -1
const CreateTaskType = 0
const UpdateTaskTitleType = 1
const UpdateTaskStatusType = 2
const UpdateTaskTitleStatusType = 3

type TaskUpdate struct {
	Id         int       `json:"id"`
	ActionType int8      `json:"action_type"`
	UserID     int       `json:"user_id"`
	TaskID     int       `json:"task_id"`
	Ts         time.Time `json:"ts"`
}
