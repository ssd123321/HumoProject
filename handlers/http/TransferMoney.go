package http

import (
	"Tasks/utils"
	"net/http"
	"strconv"
)

func (h *Handler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	sender := values.Get("senderNum")
	receiver := values.Get("receiverNum")
	money := values.Get("sum")
	sum, err := strconv.ParseFloat(money, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	senderNum, err := strconv.Atoi(sender)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	receiverNum, err := strconv.Atoi(receiver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = h.Service.ExecuteTransaction(senderNum, receiverNum, sum, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(utils.ConvertToJson("successfully", "transaction executed successfully"))
}
