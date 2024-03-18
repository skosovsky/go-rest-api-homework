package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/skosovsky/go-rest-api-homework/internal/models"
)

const (
	readWrite = 0666
	sizePart  = 64
)

type BDFile struct {
	os.File
}

func NewBDFile() (BDFile, error) {
	file, err := os.OpenFile("database.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.FileMode(readWrite))
	if err != nil {
		err = fmt.Errorf("filed to open file: %w", err)
		log.Println(err)
		return BDFile{}, err
	}

	return BDFile{*file}, nil
}

func (b *BDFile) GetTasks() map[int]models.Task {
	_, err := b.Seek(0, 0)
	if err != nil {
		return map[int]models.Task{}
	}

	var data []byte
	for {
		var countRead int
		dataPart := make([]byte, sizePart)
		countRead, err = b.Read(dataPart)

		if errors.Is(err, io.EOF) {
			break
		}

		data = append(data, dataPart[:countRead]...)

		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("filed to read file: %w", err)
			log.Println(err)
			return map[int]models.Task{}
		}
	}

	if len(data) == 0 {
		return map[int]models.Task{}
	}

	var bd []models.Task
	err = json.Unmarshal(data, &bd)
	if err != nil {
		err = fmt.Errorf("filed to unmarshal: %w", err)
		log.Println(err)
		return map[int]models.Task{}
	}

	content := map[int]models.Task{}

	for _, task := range bd {
		var id int
		id, err = strconv.Atoi(task.ID)
		if err != nil {
			continue
		}

		content[id] = task
	}

	return content
}

func (b *BDFile) GetTasksList() []models.Task {
	_, err := b.Seek(0, 0)
	if err != nil {
		return []models.Task{}
	}

	var data []byte
	for {
		var countRead int
		dataPart := make([]byte, sizePart)
		countRead, err = b.Read(dataPart)

		if errors.Is(err, io.EOF) {
			break
		}

		data = append(data, dataPart[:countRead]...)

		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("filed to read file: %w", err)
			log.Println(err)
			return []models.Task{}
		}
	}

	if len(data) == 0 {
		return []models.Task{}
	}

	var bd []models.Task
	err = json.Unmarshal(data, &bd)
	if err != nil {
		err = fmt.Errorf("filed to unmarshal: %w", err)
		log.Println(err)
		return []models.Task{}
	}

	return bd
}

func (b *BDFile) GetTask(taskID int) (*models.Task, bool) {
	_, err := b.Seek(0, 0)
	if err != nil {
		return nil, false
	}

	var data []byte
	for {
		var countRead int
		dataPart := make([]byte, sizePart)
		countRead, err = b.Read(dataPart)

		if errors.Is(err, io.EOF) {
			break
		}

		data = append(data, dataPart[:countRead]...)

		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("filed to read file: %w", err)
			log.Println(err)
			return nil, false
		}
	}

	if len(data) == 0 {
		return nil, false
	}

	var bd []models.Task
	err = json.Unmarshal(data, &bd)
	if err != nil {
		err = fmt.Errorf("filed to unmarshal: %w", err)
		log.Println(err)
		return nil, false
	}

	var task models.Task

	for i := range bd {
		if strconv.Itoa(taskID) == bd[i].ID {
			task = bd[i]
			return &task, true
		}
	}

	return nil, false
}

func (b *BDFile) AddTask(task *models.Task) bool {
	_, err := b.Seek(0, 0)
	if err != nil {
		return false
	}

	var data []byte
	for {
		var countRead int
		dataPart := make([]byte, sizePart)
		countRead, err = b.Read(dataPart)

		if errors.Is(err, io.EOF) {
			break
		}

		data = append(data, dataPart[:countRead]...)

		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("filed to read file: %w", err)
			log.Println(err)
			return false
		}
	}

	var bd []models.Task
	var isExist bool
	if len(data) != 0 {
		err = json.Unmarshal(data, &bd)
		if err != nil {
			err = fmt.Errorf("filed to unmarshal: %w", err)
			log.Println(err)
			return false
		}

		for i := range bd {
			if task.ID == bd[i].ID {
				bd[i] = *task
				isExist = true
				break
			}
		}
	}

	if !isExist {
		bd = append(bd, *task)
	}

	bdJSON, err := json.Marshal(bd)
	if err != nil {
		err = fmt.Errorf("filed to marshal: %w", err)
		log.Println(err)
		return false
	}

	err = b.Truncate(0)
	if err != nil {
		return false
	}
	_, err = b.Seek(0, 0)
	if err != nil {
		return false
	}

	_, err = b.Write(bdJSON)
	if err != nil {
		err = fmt.Errorf("filed to write file: %w", err)
		log.Println(err)
		return false
	}

	return true
}

func (b *BDFile) DeleteTask(taskID int) bool {
	_, err := b.Seek(0, 0)
	if err != nil {
		return false
	}

	var data []byte
	for {
		var countRead int
		dataPart := make([]byte, sizePart)
		countRead, err = b.Read(dataPart)

		if errors.Is(err, io.EOF) {
			break
		}

		data = append(data, dataPart[:countRead]...)

		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("filed to read file: %w", err)
			log.Println(err)
			return false
		}
	}

	if len(data) == 0 {
		return false
	}

	var bd []models.Task
	err = json.Unmarshal(data, &bd)
	if err != nil {
		err = fmt.Errorf("filed to unmarshal: %w", err)
		log.Println(err)
		return false
	}

	var isExist bool
	for i := range bd {
		if strconv.Itoa(taskID) == bd[i].ID {
			copy(bd[i:], bd[i+1:])
			bd[len(bd)-1] = models.Task{} //nolint:exhaustruct // it's clean field
			bd = bd[:len(bd)-1]
			isExist = true
			break
		}
	}

	if !isExist {
		return false
	}

	bdJSON, err := json.Marshal(bd)
	if err != nil {
		err = fmt.Errorf("filed to marshal: %w", err)
		log.Println(err)
		return false
	}

	err = b.Truncate(0)
	if err != nil {
		return false
	}

	_, err = b.Write(bdJSON)
	if err != nil {
		err = fmt.Errorf("filed to write file: %w", err)
		log.Println(err)
		return false
	}

	return true
}
