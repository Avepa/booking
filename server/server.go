package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const form = "2006-01-02"

type Room struct {
	ID          int64   `json:"room_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Date        string  `json:"date"`
}

type Booking struct {
	ID    int64  `json:"booking_id"`
	Start string `json:"date_start"`
	End   string `json:"date_end"`
}

type Handler struct {
	DB *sql.DB
}

func (h Handler) AddRoom(w http.ResponseWriter, r *http.Request) {
	desc := r.Header.Get("description")
	p := r.Header.Get("price")
	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return
	}

	res, err := h.DB.Exec(
		"INSERT INTO room (description, price, date) VALUES (?, ?, NOW())",
		desc,
		price,
	)
	if err != nil {
		return
	}

	id := struct {
		ID int64 `json:"room_id"`
	}{}
	id.ID, err = res.LastInsertId()
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		return
	}
}

// Проверить дату
func (h Handler) AddBookings(w http.ResponseWriter, r *http.Request) {
	room := r.Header.Get("room_id")
	start := r.Header.Get("date_start")
	_, err := time.Parse(form, start)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	end := r.Header.Get("date_end")
	_, err = time.Parse(form, end)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.DB.Exec(
		"INSERT INTO bookings (room_id, date_start, date_end) VALUES (?, ?, ?)",
		room,
		start,
		end,
	)
	if err != nil {
		return
	}

	id := struct {
		ID int64 `json:"booking_id"`
	}{}
	id.ID, err = res.LastInsertId()
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		return
	}
}

func (h Handler) DelRoom(w http.ResponseWriter, r *http.Request) {
	idRoom := r.Header.Get("room_id")
	id, err := strconv.ParseInt(idRoom, 10, 64)
	if err != nil {
		return
	}

	res, err := h.DB.Exec(
		"DELETE FROM room LEFT JOIN bookings ON bookings.id = room.id WHERE room.id = ?",
		id,
	)
	if err != nil {
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		return
	}
	if n == 0 {
		return
	}
}

func (h Handler) DelBookings(w http.ResponseWriter, r *http.Request) {
	idBookings := r.Header.Get("booking_id")
	id, err := strconv.ParseInt(idBookings, 10, 64)
	if err != nil {
		return
	}

	res, err := h.DB.Exec(
		"DELETE FROM bookings WHERE id = ?",
		id,
	)
	if err != nil {
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		return
	}
	if n == 0 {
		return
	}
}

//date
//date_desc
//price
//price_desc
func (h Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sorting")

	query := "SELECT * FROM room"
	switch sort {
	case "price":
		query += " ORDER BY price"
	case "price_desc":
		query += " ORDER BY price DESC"
	case "date":
		query += " ORDER BY date"
	case "date_desc":
		query += " PRDER BY date DESC"
	}
	rows, err := h.DB.Query(query)
	if err != nil {
		return
	}

	data := make([]Room, 0, 100)
	for rows.Next() {
		d := Room{}
		rows.Scan(
			&d.ID,
			&d.Description,
			&d.Price,
			&d.Date,
		)
		data = append(data, d)
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func (h Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	idRoom := r.URL.Query().Get("room_id")
	id, err := strconv.ParseInt(idRoom, 10, 64)
	if err != nil {
		return
	}
	rows, err := h.DB.Query(
		"SELECT booking.id, booking.date_start, booking.date_end "+
			"	FROM room LEFT JOIN booking ON booking.room_id = room.id "+
			"	WHERE room.id = ? ORDER BY booking.date_start",
		id,
	)
	if err != nil {
		return
	}

	data := make([]Booking, 0)
	for rows.Next() {
		d := Booking{}
		rows.Scan(
			&d.ID,
			&d.Start,
			&d.End,
		)
		data = append(data, d)
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
	return
}
