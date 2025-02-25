package handlers

import "net/http"

type UserHandler struct {
	*Handler
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
}
