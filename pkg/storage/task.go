package storage

import "time"

type Task struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Status     int8      `json:"status"`
	UserID     int       `json:"user_id"`
	CreationTs time.Time `json:"creation_ts"`
}
