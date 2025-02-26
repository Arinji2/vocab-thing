package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

// NewHandler creates a new base Handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// parseRequestBody parses the JSON body of a request into the provided struct.
func parseRequestBody(r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	return nil
}
