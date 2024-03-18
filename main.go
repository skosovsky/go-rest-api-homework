package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/skosovsky/go-rest-api-homework/database"
	"github.com/skosovsky/go-rest-api-homework/database/databases"
	_ "github.com/skosovsky/go-rest-api-homework/docs"
	"github.com/skosovsky/go-rest-api-homework/models"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// baseMode переключает разные варианты реализации "базы данных": map, slc, file.
const baseMode = "map"

type ctxKey int

const (
	keyPrincipalBD ctxKey = iota
)

type Error string

const (
	ErrNoBD       = Error("no BD")
	ErrNoAdded    = Error("task not added")
	ErrNotFound   = Error("task not found")
	ErrFiledWrite = Error("filed to write")
	ErrFiledRead  = Error("filed to read")
	ErrMarshal    = "filed to marshal response: %w"
	ErrUnmarshal  = "filed to unmarshal response: %w"
)

const (
	ReadTimeoutSeconds  = 5
	WriteTimeoutSeconds = 10
)

// initDB подключает выбранный тип "базы данных".
func initDB() database.BDer { //nolint:ireturn // it's check interface for NewBD
	switch baseMode {
	case "map":
		bd := databases.NewBDMap()
		addExampleData(bd)

		return bd
	case "slc":
		bd := databases.NewBDSlc()
		addExampleData(bd)

		return bd
	case "file":
		bd, err := databases.NewBDFile()
		if err != nil {
			panic(ErrNoBD)
		}

		return &bd
	default:
		panic(ErrNoBD)
	}
}

// addExampleData наполняет "базу данных" тестовыми данными.
func addExampleData(bd database.BDer) {
	exampleTasks := []models.Task{
		{
			ID:          "1",
			Description: "Сделать финальное задание темы REST API",
			Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
			Applications: []string{
				"VS Code",
				"Terminal",
				"git",
			},
		},
		{
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

	for _, v := range exampleTasks {
		if ok := bd.AddTask(&v); !ok { //nolint:gosec // fix in 1.22
			log.Println(ErrNoAdded)
			break
		}
	}
}

// getTasks
//
//	@Description	get all tasks in map struct
//	@ID				get-tasks
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string	"ok"
//	@Failure		500		{string}	string	"InternalServerError"
//	@Router			/tasks [get]
//
// getTasks возвращает список задач.
func getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(keyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(ErrNoBD)
		return
	}

	response, err := json.Marshal(tasks.GetTasks())
	if err != nil {
		err = fmt.Errorf(ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(ErrFiledWrite)
	}
}

// getTasksList
//
//	@Description	get all tasks in slice struct
//	@ID				get-tasks-list
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string	"ok"
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/tasks-list [get]
//
// getTasksList возвращает список задач слайсом.
func getTasksList(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(keyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(ErrNoBD)
		return
	}

	response, err := json.Marshal(tasks.GetTasksList())
	if err != nil {
		err = fmt.Errorf(ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(ErrFiledWrite)
		return
	}
}

// postTask
//
//	@Description	post task
//	@ID				post-task
//	@Accept			json
//	@Produce		json
//	@Success		201		{string}	string	"ok"
//	@Failure		400		{string}	string	"bad request"
//	@Router			/tasks [post]
//
// postTask добавляет задачу в список.
func postTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(keyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(ErrNoBD)
		return
	}

	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(ErrFiledRead)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf(ErrUnmarshal, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if ok = tasks.AddTask(&task); !ok {
		log.Println(ErrNoAdded)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getTask
//
//		@Description	get task by ID
//		@ID				get-task-by-id
//		@Accept			json
//		@Produce		json
//		@Success		201		{string}	string	"ok"
//		@Failure		400		{string}	string	"bad request"
//	 @Failure		500		{string}	string	"internal server error"
//		@Router			/tasks/{id} [get]
//
// getTask возвращает одну задачу из списка.
func getTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(keyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(ErrNoBD)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(ErrNotFound)
		return
	}

	task, ok := tasks.GetTask(id)
	if !ok {
		http.Error(w, string(ErrNotFound), http.StatusBadRequest)
		log.Println(ErrNotFound)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		err = fmt.Errorf(ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(ErrFiledWrite)
		return
	}
}

// deleteTask
//
//	@Description	delete task by ID
//	@ID				delete-task-by-id
//	@Accept			json
//	@Produce		json
//	@Success		201		{string}	string	"ok"
//	@Failure		400		{string}	string	"bad request"
//	@Router			/tasks/{id} [delete]
//
// deleteTask удаляет одну задачу из списка.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(keyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(ErrNoBD)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(ErrNotFound)
		return
	}

	if ok = tasks.DeleteTask(id); !ok {
		http.Error(w, string(ErrNotFound), http.StatusBadRequest)
		log.Println(ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//	@title			Task API
//	@version		0.1
//	@description	API Server for TODOList Application

//	@host		localhost:8080
//	@BasePath	/

func main() {
	tasks := initDB()
	file, ok := tasks.(*databases.BDFile)
	if ok {
		defer func(file *databases.BDFile) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)
	}

	router := chi.NewRouter()
	router.Get("/tasks", getTasks)
	router.Get("/tasks-list", getTasksList)
	router.Post("/tasks", postTask)
	router.Get("/tasks/{id}", getTask)
	router.Delete("/tasks/{id}", deleteTask)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // http://localhost:8080/swagger/index.html
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  ReadTimeoutSeconds * time.Second,
		WriteTimeout: WriteTimeoutSeconds * time.Second,
		ConnContext: func(_ context.Context, _ net.Conn) context.Context {
			return context.WithValue(context.Background(), keyPrincipalBD, tasks)
		},
	}

	err := server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("filed to runing server: %w", err)
		log.Println(err)
		return
	}
}
