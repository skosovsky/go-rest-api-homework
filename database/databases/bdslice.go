package databases

import (
	"strconv"

	"github.com/skosovsky/go-rest-api-homework/models"
)

type BDSlc []models.Task

func NewBDSlc() *BDSlc {
	return &BDSlc{}
}

func (b *BDSlc) GetTasks() map[int]models.Task {
	bd := map[int]models.Task{}

	for _, task := range *b {
		id, err := strconv.Atoi(task.ID)
		if err != nil {
			continue
		}

		bd[id] = task
	}

	return bd
}

func (b *BDSlc) GetTasksList() []models.Task {
	return *b
}

func (b *BDSlc) GetTask(taskID int) (*models.Task, bool) {
	var task models.Task

	for i := range *b {
		if strconv.Itoa(taskID) == (*b)[i].ID {
			task = (*b)[i]
			return &task, true
		}
	}

	return nil, false
}

func (b *BDSlc) AddTask(task *models.Task) bool {
	for i := range *b {
		if task.ID == (*b)[i].ID {
			(*b)[i] = *task
			return true
		}
	}

	*b = append(*b, *task)
	return true
}

func (b *BDSlc) DeleteTask(taskID int) bool {
	for i := range *b {
		if strconv.Itoa(taskID) == (*b)[i].ID {
			copy((*b)[i:], (*b)[i+1:])
			(*b)[len(*b)-1] = models.Task{} //nolint:exhaustruct // it's clean field
			*b = (*b)[:len(*b)-1]
			return true
		}
	}

	return false
}
