package models

type CreateTaskInput struct {
	TaskType TaskType   `json:"taskType" binding:"required"`
	Status   StatusType `json:"statusType" binding:"required"`
}

type UpdateTaskInput struct {
	ID      uint       `json:"id" binding:"required"`
	Status  StatusType `json:"status" binding:"required"`
	Message string     `json:"message" binding:"required"`
}
