package storage

import "time"

const (
	MinInt = ^int(^uint(0) >> 1)
	MaxInt = int(^uint(0) >> 1)

	TaskStatusOpened = 1
	TaskStatusClosed = 2

	TaskPriorityCreatedDelta = 10000
	TaskPriorityClosed       = MinInt
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int8      `json:"status"`
	UserID      int       `json:"user_id"`
	Priority    int       `json:"priority"`
	CreationTs  time.Time `json:"creation_ts"`
}
