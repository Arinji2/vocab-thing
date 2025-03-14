package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/go-chi/chi/v5"
)

type PhraseHandler struct {
	*Handler
}
type createPhraseRequest struct {
	Phrase     string `json:"phrase"`
	Definition string `json:"definition"`
	FoundIn    string `json:"foundIn"`
	Public     bool   `json:"public"`
}

func (p *PhraseHandler) CreatePhrase(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	var data createPhraseRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	phraseModel := database.PhraseModel{DB: p.DB}
	phraseData := models.Phrase{
		UserID:           userSession.UserID,
		Phrase:           data.Phrase,
		PhraseDefinition: data.Definition,
		FoundIn:          data.FoundIn,
		Public:           data.Public,
		CreatedAt:        time.Now().UTC(),
	}
	err := phraseModel.CreatePhrase(ctx, &phraseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, phraseData)
}

type createPhraseTagRequest struct {
	PhraseID string `json:"phraseID"`
	TagName  string `json:"tagName"`
	TagColor string `json:"tagColor"`
}

func (p *PhraseHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	var data createPhraseTagRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	phraseModel := database.PhraseModel{DB: p.DB}

	verifiedData, err := phraseModel.ByID(ctx, data.PhraseID, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tagData := models.PhraseTag{
		PhraseID:  verifiedData.Phrase.ID,
		TagName:   data.TagName,
		TagColor:  data.TagColor,
		CreatedAt: time.Now().UTC(),
	}
	err = phraseModel.CreateTag(ctx, &tagData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, tagData)
}

type getPhraseByIDRequest struct {
	ID string `json:"id"`
}

func (p *PhraseHandler) GetPhraseByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}
	phraseID := chi.URLParam(r, "id")

	phraseModel := database.PhraseModel{DB: p.DB}
	responseData, err := phraseModel.ByID(ctx, phraseID, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, responseData)
}
