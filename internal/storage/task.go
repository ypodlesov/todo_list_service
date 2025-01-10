package storage

import "time"

const (
	MinInt = -2147483648
	MaxInt = 2147483647

	TaskStatusOpened = 1
	TaskStatusClosed = 2

	TaskPriorityDelta  = 1000000
	TaskPriorityClosed = MinInt
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
