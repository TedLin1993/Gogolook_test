package handler

import (
	. "SimpleApi/model"
	"SimpleApi/validation"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var _cache sync.Map
var _idCounter atomic.Int64
var _validate = validation.Init()

func GetTasks(c *gin.Context) {
	var tasks []Task
	_cache.Range(func(key, value any) bool {
		task, ok := value.(Task)
		if ok {
			tasks = append(tasks, task)
		}
		return true
	})
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := _validate.Struct(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = fmt.Sprint(_idCounter.Add(1))
	_cache.Store(task.ID, task)
	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	value, ok := _cache.Load(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	task := value.(Task)

	var updatedTask Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedTask.Name != "" {
		task.Name = updatedTask.Name
	}
	if updatedTask.Status != nil {
		err := _validate.Struct(task)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if *updatedTask.Status == Incomplete || *updatedTask.Status == Completed {
			task.Status = updatedTask.Status
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}
	_cache.Store(taskID, task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if _, ok := _cache.LoadAndDelete(taskID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
