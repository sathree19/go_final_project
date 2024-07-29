package addNew

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (s ParcelStore) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var out Output

	params := r.URL.Query()

	param := params.Get("id")

	param1, err := strconv.Atoi(param)
	if err != nil {
		out.Error = "неверный id"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", param1))

	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", param1))
	if err != nil {
		fmt.Println(err)
		return
	}

	task = Task{}
	resp, err := json.Marshal(task)
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}