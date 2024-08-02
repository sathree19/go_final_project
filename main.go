package main

import (
	"database/sql"
	"fmt"
	"go_final_project/handlers"
	"go_final_project/middleware"
	"go_final_project/repeatTask"
	"go_final_project/storage"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"net/http"
)

func main() {

	storage.TaskDB()

	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
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

		http.HandleFunc("/api/task", middleware.Auth(handlerTask))
		http.HandleFunc("/api/tasks", middleware.Auth(handlerTasks))
		http.HandleFunc("/api/task/done", middleware.Auth(handlerTaskDone))

	} else {

		http.HandleFunc("/api/tasks", handlerTasks)
		http.HandleFunc("/api/task", handlerTask)
		http.HandleFunc("/api/task/done", handlerTaskDone)

	}

	toDoPort := strings.Join([]string{":", os.Getenv("TODO_PORT")}, "")
	fmt.Println("Server", toDoPort, "is listening...")

	err = http.ListenAndServe(toDoPort, nil)

	if err != nil {
		panic(err)
	}

}
