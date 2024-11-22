package handlers

import (
	"Tasks/model"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) AddPeopleFromFile(w http.ResponseWriter, r *http.Request) {
	file, parameters, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	values := strings.Split(parameters.Filename, ".")
	if values[len(values)-1] != "json" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("wrong file format"))
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var persons []model.Person
	err = json.Unmarshal(data, &persons)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = h.Service.AddPeople(persons, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("people created"))
}
