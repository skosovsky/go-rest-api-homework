package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	ReadTimeoutSeconds  = 5
	WriteTimeoutSeconds = 10
)

// Task структура для хранения задач.
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// tasks база данных задач.
var tasks = map[string]Task{ //nolint:gochecknoglobals // it's BD
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks возвращает список задач.
func getTasks(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		err = fmt.Errorf("filed to marshal response: %w", err)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = fmt.Errorf("filed to write response: %w", err)
		log.Println(err)
	}
}

// postTask добавляет задачу в список.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("filed to read buffer: %w", err)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf("filed to unmarshal: %w", err)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}

// getTask возвращает одну задачу из списка.
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		err := errors.New("task not found")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		err = fmt.Errorf("filed to marshal response: %w", err)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = fmt.Errorf("filed to write response: %w", err)
		log.Println(err)
	}
}

// deleteTask удаляет одну задачу из списка.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		err := errors.New("task not found")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := chi.NewRouter()

	router.Get("/tasks", getTasks)
	router.Post("/tasks", postTask)
	router.Get("/tasks/{id}", getTask)
	router.Delete("/tasks/{id}", deleteTask)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  ReadTimeoutSeconds * time.Second,
		WriteTimeout: WriteTimeoutSeconds * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("filed to runing server: %w", err)
		log.Println(err)
		return
	}
}
