package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"go_final_project/handlers"
	"go_final_project/middleware"
	"go_final_project/repeatTask"
	"go_final_project/storage"
)

func main() {

	pass := "12345"
	envPass := os.Getenv("TODO_PASSWORD")
	if len(envPass) > 0 {

		pass = envPass
	}

	addr := "storage/scheduler.db"
	envAddr := os.Getenv("TODO_DBFILE")
	if len(envAddr) > 0 {

		addr = envAddr
	}

	storage.TaskDB(addr)

	db, err := sql.Open("sqlite3", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.NewParcelStore(db)
	service := handlers.NewHandler(store)

	handlerTasks := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL

		query := url.RawQuery
		param := strings.Split(query, `=`)

		if param[0] == "search" {
			service.GetSearch(w, r)
			return
		}
		service.GetTasks(w, r)

	}

	handlerTask := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			service.PutTask(w, r)
		} else if r.Method == http.MethodGet {
			service.GetId(w, r)
		} else if r.Method == http.MethodDelete {
			service.DeleteTask(w, r)
		} else {
			service.PostTask(w, r)
		}
	}

	handlerTaskDone := func(w http.ResponseWriter, r *http.Request) {
		service.DoneTask(w, r)
	}

	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))
	http.HandleFunc("/api/nextdate", repeatTask.MainHandle)
	http.HandleFunc("/api/signin", middleware.PostSign)

	if os.Getenv("TODO_PASSWORD") != "" {

		http.HandleFunc("/api/task", middleware.Auth(handlerTask, pass))
		http.HandleFunc("/api/tasks", middleware.Auth(handlerTasks, pass))
		http.HandleFunc("/api/task/done", middleware.Auth(handlerTaskDone, pass))

	} else {

		http.HandleFunc("/api/tasks", handlerTasks)
		http.HandleFunc("/api/task", handlerTask)
		http.HandleFunc("/api/task/done", handlerTaskDone)

	}

	port := "7540"
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {

		port = envPort
	}

	toDoPort := strings.Join([]string{":", port}, "")
	fmt.Println("Server", toDoPort, "is listening...")

	err = http.ListenAndServe(toDoPort, nil)

	if err != nil {
		panic(err)
	}

}
