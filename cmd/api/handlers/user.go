package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/models"
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
		writeJSON(w, http.StatusOK, user)
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var data createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userModel := database.UserModel{DB: h.DB}
	userData := models.User{
		Email:    data.Email,
		Username: data.Username,
	}
	err := userModel.Create(ctx, &userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, http.StatusOK, userData)
}
