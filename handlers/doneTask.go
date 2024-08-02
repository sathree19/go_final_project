package handlers

import (
	"encoding/json"
	"fmt"
	"go_final_project/repeatTask"
	"go_final_project/str"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	var task str.Task
	var out str.Output

	params := r.URL.Query()

	param := params.Get("id")

	param1, err := strconv.Atoi(param)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	task, out.Error = h.Store.SelectId(param1)
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {

		_, err = h.Store.SelectId(param1)
		if err != nil {
			out.Error = err
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		out.Error = h.Store.Delete(param1)
		if out.Error != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		task, _ = h.Store.SelectId(param1)
		resp, err := json.Marshal(task)

		if err != nil {
			out.Error = err
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
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
