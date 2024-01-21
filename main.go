package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type status int

const (
	Incomplete status = iota // 0
	Completed                // 1
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status status `json:"status"`
}

var cache sync.Map
var idCounter atomic.Int64

func getTasks(c *gin.Context) {
	var tasks []Task
	cache.Range(func(key, value any) bool {
		task, ok := value.(Task)
		if ok {
			tasks = append(tasks, task)
		}
		return true
	})
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = int(idCounter.Add(1))
	cache.Store(fmt.Sprint(task.ID), task)

	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	var updatedTask Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, ok := cache.Load(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task := value.(Task)
	if updatedTask.Name != "" {
		task.Name = updatedTask.Name
	}
	if updatedTask.Status == Incomplete || updatedTask.Status == Completed {
		task.Status = updatedTask.Status
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
	cache.Store(taskID, task)

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	if _, ok := cache.LoadAndDelete(taskID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func main() {
	r := gin.Default()

	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
