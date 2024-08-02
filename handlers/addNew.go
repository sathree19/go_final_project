package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/repeatTask"
	"go_final_project/storage"
	"go_final_project/str"
	"net/http"
	"time"
)

type Handler struct {
	Store storage.ParcelStore
}

func NewHandler(store storage.ParcelStore) Handler {
	return Handler{Store: store}
}

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("GET-Метод запрещен!"))
		return
	}

	var task str.Task
	var out str.Output

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
		out.Error = errors.New("Неверный формат даты")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		out.Error = errors.New("Поле 'Заголовок' обязательное")
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
