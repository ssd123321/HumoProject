package handlers

import (
	"Tasks/model"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (p *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person model.Person
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(data, &person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", person.ID))
	UpdatedPerson, err := p.Service.UpdatePerson(&person, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	data, err = json.MarshalIndent(UpdatedPerson, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
