package main

import (
	"fmt"
	"go_final_project/dbS"
	"go_final_project/repeatTask"
	"os"

	"net/http"
)

func main() {

	//fmt.Println(repeatTask.NextDate(time.Date(2025, time.June, 29, 0, 0, 0, 0, time.Local), "20240622", "w 5,6,3,2"))
	// fmt.Println(int(time.Date(2024, 3, -1, 0, 0, 0, 0, time.Local).Day()))
	dbS.TackDB()

	walkDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(walkDir)))
	http.HandleFunc("/api/nextdate", repeatTask.MainHandle)

	toDoPort := os.Getenv("TODO_PORT")
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(toDoPort, nil)

	if err != nil {
		panic(err)
	}

}
