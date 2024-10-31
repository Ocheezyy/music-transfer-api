package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{DB: db}
}

func (tc *TaskController) GetTask(c *gin.Context) {
	taskId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("GetTask: id argument is not a uint")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var task models.Task
	tc.DB.Where("id=?", taskId).Find(&task)

	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var createTaskInput models.CreateTaskInput

	if err := c.ShouldBindJSON(&createTaskInput); err != nil {
		log.Printf("CreateTask 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask := models.Task{
		TaskType: createTaskInput.TaskType,
		Status:   createTaskInput.Status,
	}
	tc.DB.Create(&newTask)

	c.JSON(http.StatusCreated, gin.H{"data": newTask})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	var updateTaskInput models.UpdateTaskInput
	if err := c.ShouldBindJSON(&updateTaskInput); err != nil {
		log.Printf("UpdateTask 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	tc.DB.Where("id=?", updateTaskInput.ID).Find(&task)

	if task.ID == 0 {
		log.Printf("UpdateTask, id: %d not found", updateTaskInput.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	task.Status = updateTaskInput.Status
	task.Message = updateTaskInput.Message
	tc.DB.Save(&task)
	c.JSON(http.StatusOK, gin.H{})
}
