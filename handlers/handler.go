package handlers

import (
	"go_final_project/storage"
)

type Handler struct {
	Store storage.ParcelStore
}

func NewHandler(store storage.ParcelStore) Handler {
	return Handler{Store: store}
}
