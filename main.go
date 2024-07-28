package main

import (
	"fmt"
	"go_final_project/addNew"
	"go_final_project/dbS"
	"go_final_project/middleware"
	"go_final_project/repeatTask"
	"os"
	"strings"

	"net/http"
)

func handlerTasks(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.JoinPath())

	url := r.URL

	query := url.RawQuery
	param := strings.Split(query, `=`)

	if param[0] == "search" {
		addNew.GetSearch(w, r)
		return
	}
	addNew.GetTasks(w, r)

}

func handlerTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.JoinPath())

	if r.Method == http.MethodPut {
		addNew.PutTask(w, r)
	} else if r.Method == http.MethodGet {
		addNew.GetId(w, r)
	} else if r.Method == http.MethodDelete {
		addNew.DeleteTask(w, r)
	} else {
		addNew.PostTask(w, r)
	}
}

func main() {

	dbS.TackDB()
	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))
	http.HandleFunc("/api/nextdate", repeatTask.MainHandle)
	http.HandleFunc("/api/signin", middleware.PostSign)

	if os.Getenv("TODO_PASSWORD") != "" {

		http.HandleFunc("/api/task", middleware.Auth2(handlerTask))
		http.HandleFunc("/api/tasks", middleware.Auth2(handlerTasks))
		http.HandleFunc("/api/task/done", middleware.Auth2(addNew.DoneTask))

	} else {

		http.HandleFunc("/api/tasks", handlerTasks)
		http.HandleFunc("/api/task", handlerTask)
		http.HandleFunc("/api/task/done", addNew.DoneTask)

	}

	toDoPort := os.Getenv("TODO_PORT")
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(toDoPort, nil)

	if err != nil {
		panic(err)
	}

}
