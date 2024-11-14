package handlers

import (
	"Tasks/smtp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

func (p *Handler) GetPersonInFile(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	person, err := p.Service.GetPersonByID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	queries := r.URL.Query()
	email := queries.Get("email")
	if email != "" {
		f, err := os.OpenFile("Smtp_File.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = f.WriteString(string(data))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = smtp.SendFile(&smtp.Sen, "frosta123456@gmail.com", &smtp.Ser, f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		return
	}
	w.Header().Add("Content-Length", fmt.Sprint(len(string(data))))
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Disposition", "attachment;filename=data.json")
	w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
