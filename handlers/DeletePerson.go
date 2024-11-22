package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server crush..."))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	id, err = h.Service.DeletePerson(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf(w, "person with id %d deleted", id)
}
