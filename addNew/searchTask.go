package addNew

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {
	tasks := make(map[string][]Task)
	var rows *sql.Rows
	var out Output
	var id Output

	params := r.URL.Query() //chi.URLParam(r, "param")
	param := params.Get("search")

	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	defer db.Close()
	limit := 10

	param1, err := time.Parse("02.01.2006", param)
	t := param1.Format("20060102")

	if err != nil {

		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%' ORDER BY date LIMIT :limit ", sql.Named("search", param), sql.Named("search", param), sql.Named("limit", limit))
		if err != nil {
			out.Error = err.Error()
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}
		defer rows.Close()
	} else {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit ", sql.Named("date", t), sql.Named("limit", limit))
		if err != nil {
			out.Error = err.Error()
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}
		defer rows.Close()
	}

	var task Task
	var task1 []Task
	for rows.Next() {

		err := rows.Scan(&id.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return
		}
		task1 = append(task1, task)

	}
	if task1 == nil {
		task1 = []Task{}
	}
	tasks["tasks"] = task1

	if err := rows.Err(); err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
