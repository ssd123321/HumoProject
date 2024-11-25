package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	people, err := h.Service.GetPeople(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.MarshalIndent(people, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
