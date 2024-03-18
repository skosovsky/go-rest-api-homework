package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/skosovsky/go-rest-api-homework/docs"
	"github.com/skosovsky/go-rest-api-homework/internal/controllers"
	"github.com/skosovsky/go-rest-api-homework/internal/database"
	"github.com/skosovsky/go-rest-api-homework/internal/errs"
	"github.com/skosovsky/go-rest-api-homework/internal/models"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// baseMode переключает разные варианты реализации "базы данных": map, slc, file.
const baseMode = "map"

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
)

// initDB подключает выбранный тип "базы данных".
func initDB() database.BDer { //nolint:ireturn // it's check interface for NewBD
	switch baseMode {
	case "map":
		bd := database.NewBDMap()
		addExampleData(bd)

		return bd
	case "slc":
		bd := database.NewBDSlc()
		addExampleData(bd)

		return bd
	case "file":
		bd, err := database.NewBDFile()
		if err != nil {
			panic(errs.ErrNoBD)
		}

		return &bd
	default:
		panic(errs.ErrNoBD)
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
			log.Println(errs.ErrNoAdded)
			break
		}
	}
}

//	@title			Task API
//	@version		0.1
//	@description	API Server for TODOList Application

// @host		localhost:8080
// @BasePath	/

func main() {
	tasks := initDB()
	file, ok := tasks.(*database.BDFile)
	if ok {
		defer func(file *database.BDFile) {
			err := file.Close()
			if err != nil {
				log.Println(errs.ErrCloseFile)
				return
			}
		}(file)
	}

	router := chi.NewRouter()
	router.Use(middleware.SetHeader("Content-Type", "application/json, charset=utf-8"))

	router.Get("/tasks", controllers.GetTasks)
	router.Get("/tasks-list", controllers.GetTasksList)
	router.Post("/tasks", controllers.PostTask)
	router.Get("/tasks/{id}", controllers.GetTask)
	router.Delete("/tasks/{id}", controllers.DeleteTask)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // http://localhost:8080/swagger/index.html
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		ConnContext: func(_ context.Context, _ net.Conn) context.Context {
			return context.WithValue(context.Background(), controllers.KeyPrincipalBD, tasks)
		},
	}

	err := server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf(errs.ErrRunningServer, err)
		log.Println(err)
		return
	}
}
