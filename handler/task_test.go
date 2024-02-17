package handler

import (
	"SimpleApi/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	r := gin.Default()
	r.POST("/tasks", CreateTask)

	newTask := model.Task{
		Name:   "Test Task",
		Status: new(model.Status),
	}

	jsonTask, err := json.Marshal(newTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTask model.Task
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)

	assert.Equal(t, newTask.Name, createdTask.Name)
	assert.Equal(t, newTask.Status, createdTask.Status)
}

func TestGetTasks(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	_cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.GET("/tasks", GetTasks)

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []model.Task
	err = json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(tasks))

}
func TestUpdateTask(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	*existingTask.Status = model.Incomplete
	_cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.PUT("/tasks/:id", UpdateTask)

	updatedTask := model.Task{
		ID:     existingTask.ID,
		Name:   existingTask.Name,
		Status: new(model.Status),
	}
	*updatedTask.Status = model.Completed
	jsonTask, err := json.Marshal(updatedTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", updatedTask.ID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err)

	assert.Equal(t, existingTask.ID, updatedTask.ID)
	assert.Equal(t, existingTask.Name, updatedTask.Name)
	assert.NotEqual(t, existingTask.Status, updatedTask.Status)
}

func TestDeleteTask(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	*existingTask.Status = model.Incomplete
	_cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.DELETE("/tasks/:id", DeleteTask)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", existingTask.ID), nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "Task deleted successfully")
}

func TestCreateTaskInvalidInput(t *testing.T) {
	r := gin.Default()
	r.POST("/tasks", CreateTask)

	// Task with missing name
	invalidTask := model.Task{
		Status: new(model.Status),
	}

	jsonTask, err := json.Marshal(invalidTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "name should not be null")

	// Task with missing status
	invalidTask = model.Task{
		Name: "Invalid Task",
	}

	jsonTask, err = json.Marshal(invalidTask)
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "status should not be null")
}

func TestCreateTaskInvalidStatus(t *testing.T) {
	r := gin.Default()
	r.POST("/tasks", CreateTask)

	invalidTask := model.Task{
		Name:   "Invalid Task",
		Status: new(model.Status),
	}
	*invalidTask.Status = 999 // Invalid status

	jsonTask, err := json.Marshal(invalidTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid status value")
}

func TestUpdateTaskNotFound(t *testing.T) {
	r := gin.Default()
	r.PUT("/tasks/:id", UpdateTask)

	// Task not found
	nonExistingID := "100"
	updatedTask := model.Task{
		ID:     nonExistingID,
		Name:   "Updated Task",
		Status: new(model.Status),
	}
	*updatedTask.Status = model.Completed

	jsonTask, err := json.Marshal(updatedTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", nonExistingID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestUpdateTaskInvalidStatus(t *testing.T) {
	r := gin.Default()
	r.PUT("/tasks/:id", UpdateTask)

	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	*existingTask.Status = model.Incomplete
	_cache.Store(existingTask.ID, existingTask)

	// Try to update the task with an invalid status
	updatedTask := model.Task{
		ID:     existingTask.ID,
		Name:   "Updated Task",
		Status: new(model.Status),
	}
	*updatedTask.Status = 999 // Invalid status

	jsonTask, err := json.Marshal(updatedTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", updatedTask.ID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid status value")
}

func TestDeleteTaskNotFound(t *testing.T) {
	r := gin.Default()
	r.DELETE("/tasks/:id", DeleteTask)

	// Task not found
	nonExistingID := "100"
	deletedTask := model.Task{
		ID:     nonExistingID,
		Name:   "Updated Task",
		Status: new(model.Status),
	}

	jsonTask, err := json.Marshal(deletedTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", nonExistingID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}
