package handlers

import (
	"encoding/json"
	"fmt"
	"go_final_project/str"
	"net/http"
)

func (h *Handler) GetSearch(w http.ResponseWriter, r *http.Request) {

	var out str.Output
	var tasks map[string][]str.Task

	params := r.URL.Query()
	param := params.Get("search")

	limit := 10

	out.Error, tasks = h.Store.Search(param, limit)
	if out.Error != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		out.Error = err
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
