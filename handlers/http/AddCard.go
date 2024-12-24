package http

import (
	"Tasks/utils"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELOOOOOOOO")
	values := r.URL.Query()
	strID := values.Get("person_id")
	bankName := values.Get("bank_name")
	personID, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	if r.Context().Value("id").(int) != personID {
		w.WriteHeader(http.StatusForbidden)
		w.Write(utils.ConvertToJson("status", "forbidden"))
		return
	}
	_, err = h.Service.AddCard(r.Context(), bankName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ConvertToJson("error", err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(utils.ConvertToJson("successfully", "Card created successfully"))
}
