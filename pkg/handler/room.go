package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Avepa/booking/pkg"
)

type roomID struct {
	ID int64 `json:"room_id"`
}

// example request:
//		http://localhost/room/add
// header that are used:
//		"description",
//		"price"
func (h *Handler) addRoom(w http.ResponseWriter, r *http.Request) {
	var err error
	var room pkg.Room

	room.Description = r.Header.Get("description")
	price := r.Header.Get("price")
	room.Price, err = strconv.ParseFloat(price, 64)
	if err != nil {
		err = pkg.ErrPriceNotValid
		log.Println(err)
		HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := roomID{}
	id.ID, err = h.services.Room.Add(&room)
	if err != nil {
		log.Println(err)
		if err == pkg.ErrFailedSave {
			HTTPError(w, err.Error(), http.StatusInternalServerError)
		} else {
			HTTPError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	json.NewEncoder(w).Encode(id)
}

// example request:
//		http://localhost/room/list?sorting=data
// default rooms are sorted by descending date
// rooms can be sorted by:
//	 price increase - price;
//	 descending price - price_desc;
//	 date - date;
//   descending date - date_desc, лиюо, любое другое значение
func (h *Handler) getRoom(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sorting")
	rooms, err := h.services.Room.Get(sort)
	if err != nil {
		log.Println(err)
		HTTPError(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

// example request:
//		http://localhost/room/delete?room_id=12
func (h *Handler) deleteRoom(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("room_id")
	room, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = pkg.ErrIdNotValid
		log.Println(err)
		HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.Room.Delete(room)
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
