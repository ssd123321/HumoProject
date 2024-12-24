package handlers

import (
	"Tasks/model"
	"Tasks/smtp"
	"Tasks/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var person model.SigningRequest
	err = json.Unmarshal(data, &person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	refreshJWT, accessJWT, err := h.Service.Login(&person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("accessJWT", accessJWT)
	w.Header().Add("refreshJWT", refreshJWT)
	w.WriteHeader(http.StatusOK)
	w.Write(utils.ConvertToJson("successfully", "logged in successfully"))
}
func (h *Handler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("Здесь")
	contextID := r.Context().Value("id").(int)
	if id != contextID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("not permissible"))
		return
	}
	fmt.Println(223)
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
func (h *Handler) GetPersonInFile(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	contextID := r.Context().Value("id").(int)
	if id != contextID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "person_id", id))
	person, err := h.Service.GetPersonByID(r.Context())
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
}
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	fmt.Println(123333)
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
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	value := values["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server crush..."))
		return
	}
	contextID := r.Context().Value("id").(int)
	if id != contextID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	id, err = h.Service.DeletePerson(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(utils.ConvertToJson("successfully", "person deleted"))
}
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	var Person *model.Person
	err = json.Unmarshal(data, &Person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	log.Printf("+%v", Person)
	Person, err = h.Service.AddPerson(Person, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	data, err = json.MarshalIndent(Person, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.ConvertToJson("successfully", "user created")
	defer r.Body.Close()
}
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
	w.Write(utils.ConvertToJson("successfully", "people from file created"))
}
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
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
	contextID := r.Context().Value("id").(int)
	if person.ID != contextID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", person.ID))
	UpdatedPerson, err := h.Service.UpdatePerson(&person, r.Context())
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
