package main

import (
	"go_final_project/dbS"
	"net/http"
	"os"
)

func main() {

	dbS.TackDB()

	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))

	toDoPort := os.Getenv("TODO_PORT")

	err := http.ListenAndServe(toDoPort, nil)
	if err != nil {
		panic(err)
	}

}
