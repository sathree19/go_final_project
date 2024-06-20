package main

import (
	"net/http"
	"os"
)

func main() {
	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))

	toDoPort := os.Getenv("TODO_PORT")

	err := http.ListenAndServe(toDoPort, nil)
	if err != nil {
		panic(err)
	}
}
