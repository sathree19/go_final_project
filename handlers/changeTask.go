package addNew

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/repeatTask"
)

func (s ParcelStore) PutTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	var out Output
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		out.Error = "Задача не найдена"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	var ids []int
	var ID int

	rows, err := s.db.Query("SELECT id FROM scheduler")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&ID)
		if err != nil {
			log.Println(err)
			return
		}

		ids = append(ids, ID)
	}

	if !repeatTask.Contains(ids, int(task.Id)) {
		out.Error = "Задача не найдена"
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

	if taskDate.Sub(time.Now()) <= 0 {
		switch {
		case task.Repeat == "":
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

	_, err = s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		out.Error = "Задача не найдена"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	task = Task{Id: task.Id, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat}

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
