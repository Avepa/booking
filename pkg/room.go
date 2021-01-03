package pkg

type Room struct {
	ID          int64   `json:"room_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Date        string  `json:"date"`
}
