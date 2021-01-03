package pkg

type Booking struct {
	ID    int64  `json:"booking_id"`
	Start string `json:"date_start"`
	End   string `json:"date_end"`
}
