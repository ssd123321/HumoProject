package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	strID := values.Get("person_id")
	bankName := values.Get("bankname")
	personID, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(personID)
	r = r.WithContext(context.WithValue(r.Context(), "person_id", personID))
	_, err = h.Service.AddCard(r.Context(), bankName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Card created successfully"))
}
