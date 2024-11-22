package handlers

import (
	"Tasks/smtp"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func (h *Handler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	person, err := h.Service.GetPersonByID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	v := r.URL.Query()
	email := v.Get("email")
	if re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`); !re.MatchString(email) && email != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong email format"))
		return
	} else if email != "" {
		err = smtp.SendMessage(&smtp.Sen, email, &smtp.Ser, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("message sent to email"))
			return
		}
	}
	log.Printf("sent data:%s", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
