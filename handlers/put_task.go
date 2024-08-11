package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/model"
	"go_final_project/repeatTask"
)

// Редактирование задачи
func (h *Handler) PutTask(w http.ResponseWriter, r *http.Request) {

	var task model.Task
	var out model.Output
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		out.Error = errors.New("Задача не найдена")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	var tasks []model.Task
	var ids []int

	tasks, out.Error = h.Store.SelectAll("ALL")
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
	}

	for _, taskId := range tasks {
		ids = append(ids, int(taskId.Id))
	}

	if !repeatTask.Contains(ids, int(task.Id)) {
		out.Error = errors.New("Задача не найдена")
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

	if taskDate.Sub(time.Now()) <= 0 {
		switch {
		case task.Repeat == "":
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

	task, out = h.Store.Update(task, out)

	resp, err := json.Marshal(task)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
}
