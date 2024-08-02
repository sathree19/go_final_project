package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("Другой метод"))
		return
	}

	limit := "20"
	_, tasks, err := h.Store.SelectAll(limit)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
