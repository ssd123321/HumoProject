package http

import (
	"Tasks/Service"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service *Service.Service
	Router  *mux.Router
}

func NewHandler(service *Service.Service, router *mux.Router) *Handler {
	return &Handler{Service: service, Router: router}
}
