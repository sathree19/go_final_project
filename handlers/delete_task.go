package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go_final_project/model"
)

// Удаление задачи
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var out model.Output

	params := r.URL.Query()

	id := params.Get("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		out.Error = errors.New("неверный id")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, out.Error), http.StatusBadRequest)
		return
	}
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
