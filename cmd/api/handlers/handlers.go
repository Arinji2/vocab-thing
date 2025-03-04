package handlers

import (
	"database/sql"
)

type Handler struct {
	DB *sql.DB
}

// NewHandler creates a new base Handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}
