package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Avepa/booking/pkg"
)

type bookingID struct {
	ID int64 `json:"booking_id"`
}

// example request:
//		http://localhost/bookings/create
//
// headers that are used:
//		room_id
//		date_start
//		date_end
//
// date format: 2006-01-02
func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	booking := pkg.Booking{
		Start: r.Header.Get("date_start"),
		End:   r.Header.Get("date_end"),
	}

	room := r.Header.Get("room_id")
	idRoom, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		err = pkg.ErrIdNotValid
		log.Println(err)
		HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := bookingID{}
	id.ID, err = h.services.Bookings.Add(idRoom, &booking)
	if err != nil {
		log.Println(err)
		if err == pkg.ErrNoForeignKey || err == pkg.ErrDateIsIncorrect {
			HTTPError(w, err.Error(), http.StatusBadRequest)
		} else {
			HTTPError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(id)
}

// example request:
//		http://localhost/bookings/list?room_id=12
func (h *Handler) getBookings(w http.ResponseWriter, r *http.Request) {
	idRoom := r.URL.Query().Get("room_id")
	id, err := strconv.ParseInt(idRoom, 10, 64)
	if err != nil {
		err = pkg.ErrIdNotValid
		log.Println(err)
		HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookings, err := h.services.Bookings.Get(id)
	if err != nil {
		log.Println(err)
		if err == pkg.ErrFailedGet {
			HTTPError(w, err.Error(), http.StatusInternalServerError)
		} else {
			HTTPError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	json.NewEncoder(w).Encode(bookings)
}

// example request:
//		http://localhost/bookings/delete?booking_id=245
func (h *Handler) deleteBookings(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("booking_id")
	booking, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = pkg.ErrIdNotValid
		log.Println(err)
		HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.Bookings.Delete(booking)
	if err != nil {
		log.Println(err)
		if err == pkg.ErrFailedDelete {
			HTTPError(w, err.Error(), http.StatusInternalServerError)
		}
		if err == pkg.ErrIDNotFound {
			HTTPError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
}
