package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/skosovsky/go-rest-api-homework/internal/database"
	"github.com/skosovsky/go-rest-api-homework/internal/errs"
	"github.com/skosovsky/go-rest-api-homework/internal/models"
)

type ctxKey int

const (
	KeyPrincipalBD ctxKey = iota
)

// GetTasks
//
//	@Description	get all tasks in map struct
//	@ID				get-tasks
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string	"ok"
//	@Failure		500		{string}	string	"InternalServerError"
//	@Router			/tasks [get]
//
// GetTasks возвращает список задач.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(KeyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(errs.ErrNoBD)
		return
	}

	response, err := json.Marshal(tasks.GetTasks())
	if err != nil {
		err = fmt.Errorf(errs.ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(errs.ErrFiledWrite)
	}
}

// GetTasksList
//
//	@Description	get all tasks in slice struct
//	@ID				get-tasks-list
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string	"ok"
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/tasks-list [get]
//
// GetTasksList возвращает список задач слайсом.
func GetTasksList(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(KeyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(errs.ErrNoBD)
		return
	}

	response, err := json.Marshal(tasks.GetTasksList())
	if err != nil {
		err = fmt.Errorf(errs.ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(errs.ErrFiledWrite)
		return
	}
}

// PostTask
//
//	@Description	post task
//	@ID				post-task
//	@Accept			json
//	@Produce		json
//	@Success		201		{string}	string	"ok"
//	@Failure		400		{string}	string	"bad request"
//	@Router			/tasks [post]
//
// PostTask добавляет задачу в список.
func PostTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(KeyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(errs.ErrNoBD)
		return
	}

	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(errs.ErrFiledRead)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf(errs.ErrUnmarshal, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if ok = tasks.AddTask(&task); !ok {
		log.Println(errs.ErrNoAdded)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTask
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
// GetTask возвращает одну задачу из списка.
func GetTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(KeyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(errs.ErrNoBD)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(errs.ErrNotFound)
		return
	}

	task, ok := tasks.GetTask(id)
	if !ok {
		http.Error(w, string(errs.ErrNotFound), http.StatusBadRequest)
		log.Println(errs.ErrNotFound)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		err = fmt.Errorf(errs.ErrMarshal, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println(errs.ErrFiledWrite)
		return
	}
}

// DeleteTask
//
//	@Description	delete task by ID
//	@ID				delete-task-by-id
//	@Accept			json
//	@Produce		json
//	@Success		201		{string}	string	"ok"
//	@Failure		400		{string}	string	"bad request"
//	@Router			/tasks/{id} [delete]
//
// DeleteTask удаляет одну задачу из списка.
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	tasks, ok := r.Context().Value(KeyPrincipalBD).(database.BDer)
	if !ok {
		log.Println(errs.ErrNoBD)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(errs.ErrNotFound)
		return
	}

	if ok = tasks.DeleteTask(id); !ok {
		http.Error(w, string(errs.ErrNotFound), http.StatusBadRequest)
		log.Println(errs.ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
