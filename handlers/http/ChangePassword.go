package http

import (
	"Tasks/model"
	"Tasks/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ConvertToJson("unsuccessfully", err.Error()))
		return
	}
	var n model.NewPassword
	err = json.Unmarshal(data, &n)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ConvertToJson("unsuccessfully", err.Error()))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "id", n.Id))
	if n.Id != r.Context().Value("id").(int) {
		w.WriteHeader(http.StatusForbidden)
		w.Write(utils.ConvertToJson("unsuccessfully", "access is prohibited"))
		return
	}
	err = h.Service.ChangePassword(n.NewPass, n.CurrentPassword, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ConvertToJson("unsuccessfully", err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(utils.ConvertToJson("successfully", "password changed successfully"))
}
