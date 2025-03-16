package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/httpmiddleware"
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

func (p *PhraseHandler) GetAllPhrases(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}
	paginationData, exists := httpmiddleware.PaginationFromContext(ctx)
	if !exists {
		http.Error(w, "no pagination data found", http.StatusInternalServerError)
		return
	}
	phraseModel := database.PhraseModel{DB: p.DB}
	responseData, err := phraseModel.All(ctx, paginationData.Page, paginationData.PageSize, paginationData.Sorting.SortBy, paginationData.Sorting.Order, paginationData.Sorting.GroupBy, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, responseData)
}

func (p *PhraseHandler) SearchPhrases(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}
	searchingData, exists := httpmiddleware.SearchingFromContext(ctx)
	if !exists {
		http.Error(w, "no searching data found", http.StatusInternalServerError)
		return
	}
	phraseModel := database.PhraseModel{DB: p.DB}
	responseData, err := phraseModel.Search(ctx, searchingData.Term, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, responseData)
}

type updatePhraseRequest struct {
	Phrase models.Phrase `json:"phrase"`
}

func (p *PhraseHandler) UpdatePhrase(w http.ResponseWriter, r *http.Request) {
	phraseID := chi.URLParam(r, "id")
	if phraseID == "" {
		http.Error(w, "no phrase id found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	var data updatePhraseRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	phraseModel := database.PhraseModel{DB: p.DB}
	data.Phrase.ID = phraseID

	err := phraseModel.UpdatePhrase(ctx, &data.Phrase, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, data.Phrase)
}

type updateTagRequest struct {
	Tag models.PhraseTag `json:"tag"`
}

func (p *PhraseHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	phraseID := chi.URLParam(r, "phraseID")
	tagID := chi.URLParam(r, "tagID")
	if phraseID == "" {
		http.Error(w, "no phrase id found", http.StatusBadRequest)
		return
	}
	if tagID == "" {
		http.Error(w, "no tag id found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	var data updateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	phraseModel := database.PhraseModel{DB: p.DB}
	data.Tag.ID = tagID
	data.Tag.PhraseID = phraseID

	err := phraseModel.UpdateTag(ctx, &data.Tag, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, data.Tag)
}

func (p *PhraseHandler) DeletePhrase(w http.ResponseWriter, r *http.Request) {
	phraseID := chi.URLParam(r, "id")
	if phraseID == "" {
		http.Error(w, "no phrase id found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	phraseModel := database.PhraseModel{DB: p.DB}

	err := phraseModel.DeletePhrase(ctx, phraseID, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PhraseHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	phraseID := chi.URLParam(r, "phraseID")
	tagID := chi.URLParam(r, "tagID")
	if phraseID == "" {
		http.Error(w, "no phrase id found", http.StatusBadRequest)
		return
	}
	if tagID == "" {
		http.Error(w, "no tag id found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	phraseModel := database.PhraseModel{DB: p.DB}

	err := phraseModel.DeleteTag(ctx, phraseID, tagID, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
