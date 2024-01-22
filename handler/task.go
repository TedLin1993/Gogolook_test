package handler

import (
	"Gogolook_test/model"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var cache sync.Map
var idCounter atomic.Int64

func GetTasks(c *gin.Context) {
	var tasks []model.Task
	cache.Range(func(key, value any) bool {
		task, ok := value.(model.Task)
		if ok {
			tasks = append(tasks, task)
		}
		return true
	})
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task model.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if task.Name == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "name should not be null"})
		return
	}
	if task.Status == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "status should not be null"})
		return
	}
	if *task.Status != model.Incomplete && *task.Status != model.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
	task.ID = fmt.Sprint(idCounter.Add(1))
	cache.Store(task.ID, task)
	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	value, ok := cache.Load(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	task := value.(model.Task)

	var updatedTask model.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedTask.Name != "" {
		task.Name = updatedTask.Name
	}
	if updatedTask.Status != nil {
		if *updatedTask.Status == model.Incomplete || *updatedTask.Status == model.Completed {
			task.Status = updatedTask.Status
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}
	cache.Store(taskID, task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if _, ok := cache.LoadAndDelete(taskID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
