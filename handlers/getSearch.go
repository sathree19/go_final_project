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
	var task1 []model.Task

	params := r.URL.Query()
	param := params.Get("search")

	out.Error, task1 = h.Store.Search(param, model.Limit)
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
	tasks := make(map[string][]model.Task)

	tasks["tasks"] = task1

	resp, err := json.Marshal(tasks)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(err)

}
