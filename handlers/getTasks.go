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

	_, task1, err := h.Store.SelectAll(model.Limit)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	tasks := make(map[string][]model.Task)

	tasks["tasks"] = task1

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(err)

}
