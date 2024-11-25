package handlers

import (
	"Tasks/model"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var card model.Card
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(data, &card)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", card.PersonID))
	createdCard, err := h.Service.AddCard(&card, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	data, err = json.MarshalIndent(createdCard, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
