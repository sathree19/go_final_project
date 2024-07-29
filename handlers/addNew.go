package addNew

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/repeatTask"
	"net/http"
	"time"
)

type Task struct {
	Id      int64  `json:"id,string,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Output struct {
	ID    int64  `json:"id,string,omitempty"`
	Error string `json:"error"`
}

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) PostTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("GET-Метод запрещен!"))
		return
	}

	var task Task
	var out Output

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format("20060102")
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		out.Error = "Неверный формат даты"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		out.Error = "Поле 'Заголовок' обязательное"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if taskDate.Sub(time.Now()) < 0 {
		switch {
		case task.Repeat == "":
			task.Date = time.Now().Format("20060102")
		case time.Now().Format("20060102") == task.Date:
			task.Date = time.Now().Format("20060102")

		default:
			task.Date, err = repeatTask.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				out.Error = err.Error()
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
				return
			}

		}

	}

	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	out.ID, err = res.LastInsertId()

	if err != nil {
		out.Error = err.Error()
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, `{"id": "%d"}`, out.ID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}
