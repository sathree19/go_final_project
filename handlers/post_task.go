package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go_final_project/model"
	"go_final_project/repeatTask"
)

// Добавляет новую задачу
func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {

	var task model.Task
	var out model.Output

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Date == "" || task.Date == "today" {
		task.Date = time.Now().Format("20060102")
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		out.Error = errors.New("неверный формат даты")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		out.Error = errors.New("поле 'Заголовок' обязательное")
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
				out.Error = err
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
				return
			}

		}

	}

	out = h.Store.Add(task, out)

	fmt.Fprintf(w, `{"id": "%d"}`, out.ID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}
