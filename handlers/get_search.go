package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go_final_project/model"
)

// Поиск задачи
func (h *Handler) GetSearch(w http.ResponseWriter, r *http.Request) {

	var out model.Output
	var tasks []model.Task

	params := r.URL.Query()
	param := params.Get("search")

	tasks, out.Error = h.Store.Search(param, model.LimitShowTasks)
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	tasksMap := make(map[string][]model.Task)

	tasksMap["tasks"] = tasks

	resp, err := json.Marshal(tasksMap)
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
