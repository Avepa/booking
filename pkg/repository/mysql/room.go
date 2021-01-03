package mysql

import (
	"database/sql"

	"github.com/Avepa/booking/pkg"
)

type RoomMySQL struct {
	db *sql.DB
}

func NewRoomMySQL(db *sql.DB) *RoomMySQL {
	return &RoomMySQL{db: db}
}

// uses fields: Description, Price.
// in the id field records the room id.
func (r *RoomMySQL) Add(room *pkg.Room) error {
	res, err := r.db.Exec(
		"INSERT INTO room (description, price, date) VALUES (?, ?, NOW())",
		room.Description,
		room.Price,
	)
	if err != nil {
		return pkg.ErrFailedSave
	}

	room.ID, err = res.LastInsertId()
	return err
}

func (r *RoomMySQL) Delete(id int64) error {
	res, err := r.db.Exec(
		"DELETE FROM `room` WHERE `id` = ?",
		id,
	)
	if err != nil {
		return pkg.ErrFailedDelete
	}

	n, err := res.RowsAffected()
	if err != nil {
		return pkg.ErrFailedDelete
	}
	if n == 0 {
		return pkg.ErrIDNotFound
	}

	return nil
}

func (r *RoomMySQL) get(query string) ([]pkg.Room, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	rooms := make([]pkg.Room, 0, 2)
	for rows.Next() {
		room := pkg.Room{}
		rows.Scan(
			&room.ID,
			&room.Date,
			&room.Price,
			&room.Description,
		)
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomMySQL) GetByDate() ([]pkg.Room, error) {
	query := "SELECT `id`, `date`, `price`, `description` FROM room" +
		" ORDER BY date"
	return r.get(query)
}

func (r *RoomMySQL) GetByDateDESC() ([]pkg.Room, error) {
	query := "SELECT `id`, `date`, `price`, `description` FROM room" +
		" ORDER BY date DESC"
	return r.get(query)
}

func (r *RoomMySQL) GetByPrice() ([]pkg.Room, error) {
	query := "SELECT `id`, `date`, `price`, `description` FROM room" +
		" ORDER BY price"
	return r.get(query)
}

func (r *RoomMySQL) GetByPriceDESC() ([]pkg.Room, error) {
	query := "SELECT `id`, `date`, `price`, `description` FROM room" +
		" ORDER BY price DESC"
	return r.get(query)
}
