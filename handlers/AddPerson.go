package handlers

import (
	"Tasks/model"
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var Person *model.Person
	err = json.Unmarshal(data, &Person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	Person, err = h.Service.AddPerson(Person, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server crush..."))
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	data, err = json.MarshalIndent(Person, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server crush..."))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	defer r.Body.Close()
}
