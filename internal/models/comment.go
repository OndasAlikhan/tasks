package models

import "time"

type Comment struct {
	Id        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	Text      string    `json:"string"`
	CreatedAt time.Time `json:"created_at"`
}
