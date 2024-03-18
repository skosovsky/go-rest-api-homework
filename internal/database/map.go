package database

import (
	"strconv"

	"github.com/skosovsky/go-rest-api-homework/internal/models"
)

type BDMap map[int]models.Task

func NewBDMap() *BDMap {
	return &BDMap{}
}

func (b *BDMap) GetTasks() map[int]models.Task {
	return *b
}

func (b *BDMap) GetTasksList() []models.Task {
	bd := make([]models.Task, 0, len(*b))

	for _, task := range *b {
		bd = append(bd, task)
	}

	return bd
}

func (b *BDMap) GetTask(taskID int) (*models.Task, bool) {
	if _, ok := (*b)[taskID]; !ok {
		return nil, false
	}

	task := (*b)[taskID]
	return &task, true
}

func (b *BDMap) AddTask(task *models.Task) bool {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return false
	}

	(*b)[id] = *task
	return true
}

func (b *BDMap) DeleteTask(taskID int) bool {
	if _, ok := (*b)[taskID]; !ok {
		return false
	}

	delete(*b, taskID)
	return true
}
