package models

import "time"

type TransferLog struct {
	ID         uint       `json:"id" gorm:"primary_key"`
	UserID     uint       `json:"userId"`
	PlaylistID uint       `json:"playlistID"`
	Status     StatusType `json:"status"`
	Message    string     `json:"message"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type StatusType string

const (
	InProgress StatusType = "in_progress"
	Completed  StatusType = "completed"
	Failed     StatusType = "failed"
)
