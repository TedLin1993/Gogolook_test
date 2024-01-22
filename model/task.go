package model

type Status int

const (
	Incomplete Status = iota // 0
	Completed                // 1
)

type Task struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Status *Status `json:"status"`
}
