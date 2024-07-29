package addNew

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s ParcelStore) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := make(map[string][]Task)

	var out Output

	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("Другой метод"))
		return
	}

	limit := 20

	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	var task Task
	var task1 []Task

	for rows.Next() {

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return
		}
		task1 = append(task1, task)

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
