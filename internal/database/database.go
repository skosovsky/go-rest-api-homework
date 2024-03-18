package database

import (
	"github.com/skosovsky/go-rest-api-homework/internal/models"
)

type BDer interface {
	GetTasks() map[int]models.Task
	GetTasksList() []models.Task
	GetTask(id int) (*models.Task, bool)
	AddTask(task *models.Task) bool
	DeleteTask(id int) bool
}
