package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/database"
)

type UserHandler struct {
	*Handler
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	userModel := database.UserModel{DB: h.DB}
	users, err := userModel.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, user := range users {
		w.Write([]byte(user.Username))
	}
}
