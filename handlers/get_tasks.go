package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go_final_project/model"
)

// Вывод задач
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.Store.SelectAll(model.LimitShowTasks)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	tasksMap := make(map[string][]model.Task)

	tasksMap["tasks"] = tasks

	resp, err := json.Marshal(tasksMap)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
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
