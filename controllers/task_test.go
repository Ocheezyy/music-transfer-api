package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/Ocheezyy/music-transfer-api/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TaskResponse struct {
	Data models.Task `json:"data"`
}

func TestGetTask_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	taskType := models.PlaylistTransfer
	status := models.Queued

	rows := sqlmock.NewRows([]string{"id", "task_type", "status", "message"}).
		AddRow(1, taskType, status, nil)

	mock.ExpectQuery(`SELECT \* FROM "tasks"`).
		WithArgs(1).
		WillReturnRows(rows)

	controller := NewTaskController(db)

	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", 1)}}
	c.Set("currentUser", models.User{ID: 1})

	controller.GetTask(c)

	var taskRes TaskResponse
	err := json.Unmarshal(res.Body.Bytes(), &taskRes)
	if err != nil {
		log.Printf("Failed to unmarshal response %s", err)
	}
	assert.Nil(t, err)
	task := taskRes.Data

	assert.Equal(t, status, task.Status)
	assert.Equal(t, taskType, task.TaskType)
	assert.Equal(t, http.StatusOK, res.Code)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestGetTask_NotFound(t *testing.T) {
	db, mock := test.NewMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "tasks"`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	controller := NewTaskController(db)

	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("currentUser", models.User{ID: 1})

	controller.GetTask(c)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Contains(t, res.Body.String(), `"error":"task not found"`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestCreateTask_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	taskType := models.PlaylistTransfer
	status := models.Queued
	body := models.CreateTaskInput{
		TaskType: taskType,
		Status:   status,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery(`INSERT INTO "tasks"`).
		WithArgs(
			sqlmock.AnyArg(),
			body.TaskType,
			body.Status,
		).
		WillReturnRows(rows)

	controller := NewTaskController(db)

	c, res := getTestContext()
	c.Set("currentUser", models.User{ID: 1})
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req
	controller.CreateTask(c)

	var taskRes TaskResponse
	err := json.Unmarshal(res.Body.Bytes(), &taskRes)
	if err != nil {
		log.Printf("Failed to unmarshal response: %s", err)
	}
	assert.Nil(t, err)
	task := taskRes.Data

	assert.Equal(t, status, task.Status)
	assert.Equal(t, taskType, task.TaskType)
	assert.Equal(t, http.StatusCreated, res.Code)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

// func TestUpdateTask_Success(t *testing.T) {

// }

// func TestUpdateTask_NotFound(t *testing.T) {

// }
