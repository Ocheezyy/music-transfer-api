package models

import "time"

type Task struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	TaskType  TaskType   `json:"taskType"`
	Status    StatusType `json:"status"`
	Message   string     `json:"message"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskType string

const (
	PlaylistTransfer  TaskType = "playlist_transfer"
	EmailNotification TaskType = "email_notification"
)
