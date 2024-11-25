package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) GetCardByPersonID(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	c, err := h.Service.GetCardByPersonID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
