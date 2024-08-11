package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_final_project/model"
	"go_final_project/repeatTask"
)

// Выполнение задачи
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var out model.Output

	params := r.URL.Query()

	id := params.Get("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	task, out.Error = h.Store.SelectId(idInt)
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {

		_, err = h.Store.SelectId(idInt)
		if err != nil {
			out.Error = err
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		out.Error = h.Store.Delete(idInt)
		if out.Error != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		task, _ = h.Store.SelectId(idInt)
		if err != nil {
			out.Error = err
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}
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
	} else {

		task.Date, err = repeatTask.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {

			out.Error = err
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		task, out = h.Store.Update(task, out)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{}")
	}

}
