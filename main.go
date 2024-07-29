package main

import (
	"database/sql"
	"fmt"
	addNew "go_final_project/handlers"
	"go_final_project/middleware"
	"go_final_project/repeatTask"
	"go_final_project/storage"
	"log"
	"os"
	"strings"

	"net/http"
)

type ParcelService struct {
	store addNew.ParcelStore
}

func NewParcelService(store addNew.ParcelStore) ParcelService {
	return ParcelService{store: store}
}

func (s ParcelService) handlerTasks(w http.ResponseWriter, r *http.Request) {
	url := r.URL

	query := url.RawQuery
	param := strings.Split(query, `=`)

	if param[0] == "search" {
		s.store.GetSearch(w, r)
		return
	}
	s.store.GetTasks(w, r)

}

func (s ParcelService) handlerTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		s.store.PutTask(w, r)
	} else if r.Method == http.MethodGet {
		s.store.GetId(w, r)
	} else if r.Method == http.MethodDelete {
		s.store.DeleteTask(w, r)
	} else {
		s.store.PostTask(w, r)
	}
}

func (s ParcelService) handlerTaskDone(w http.ResponseWriter, r *http.Request) {
	s.store.DoneTask(w, r)
}

func main() {

	storage.TackDB()
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := addNew.NewParcelStore(db)
	service := NewParcelService(store)

	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))
	http.HandleFunc("/api/nextdate", repeatTask.MainHandle)
	http.HandleFunc("/api/signin", middleware.PostSign)

	if os.Getenv("TODO_PASSWORD") != "" {

		http.HandleFunc("/api/task", middleware.Auth(service.handlerTask))
		http.HandleFunc("/api/tasks", middleware.Auth(service.handlerTasks))
		http.HandleFunc("/api/task/done", middleware.Auth(service.handlerTaskDone))

	} else {

		http.HandleFunc("/api/tasks", service.handlerTasks)
		http.HandleFunc("/api/task", service.handlerTask)
		http.HandleFunc("/api/task/done", service.handlerTaskDone)

	}

	toDoPort := strings.Join([]string{":", os.Getenv("TODO_PORT")}, "")
	fmt.Println("Server", toDoPort, "is listening...")

	err = http.ListenAndServe(toDoPort, nil)

	if err != nil {
		panic(err)
	}

}
