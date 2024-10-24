package models

type CreateTransferLogInput struct {
	UserID     uint       `json:"userId" binding:"required"`
	PlaylistID uint       `json:"playlistId" binding:"required"`
	Status     StatusType `json:"status" binding:"required"`
	Message    string     `json:"message"`
}

type UpdateTransferLogInput struct {
	ID      uint       `json:"id" binding:"required"`
	Status  StatusType `json:"status" binding:"required"`
	Message string     `json:"message"`
}
