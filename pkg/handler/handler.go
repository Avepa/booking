package handler

import (
	"github.com/gorilla/mux"

	"github.com/Avepa/booking/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/room/add", h.addRoom).Methods("POST")
	router.HandleFunc("/room/list", h.getRoom).Methods("GET")
	router.HandleFunc("/room/delete", h.deleteRoom).Methods("DELETE")

	router.HandleFunc("/bookings/create", h.createBooking).Methods("POST")
	router.HandleFunc("/bookings/list", h.getBookings).Methods("GET")
	router.HandleFunc("/bookings/delete", h.deleteBookings).Methods("DELETE")

	return router
}
