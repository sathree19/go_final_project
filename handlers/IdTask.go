package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/str"
	"net/http"
	"strconv"
)

func (h *Handler) GetId(w http.ResponseWriter, r *http.Request) {
	var task str.Task
	var out str.Output

	params := r.URL.Query()

	param := params.Get("id")

	if param == "" {
		out.Error = errors.New("Задача не найдена")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	param1, err := strconv.Atoi(param)
	if err != nil {
		out.Error = errors.New("Задача не найдена")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}

	task, err = h.Store.SelectId(param1)
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
	w.Write(resp)

}
