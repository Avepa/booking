package mysql

import (
	"database/sql"

	"github.com/Avepa/booking/pkg"
)

type BookingsMySQL struct {
	db *sql.DB
}

func NewBookingsMySQL(db *sql.DB) *BookingsMySQL {
	return &BookingsMySQL{db: db}
}

func (r *BookingsMySQL) Add(room int64, bookings *pkg.Booking) error {
	res, err := r.db.Exec(
		"INSERT INTO bookings (room_id, date_start, date_end) VALUES (?, ?, ?)",
		room,
		bookings.Start,
		bookings.End,
	)
	if err != nil {
		n := len(pkg.ErrNoForeignKey.Error())
		if len(err.Error()) >= n {
			if err.Error()[:n] == pkg.ErrNoForeignKey.Error() {
				err = pkg.ErrNoForeignKey
			}
		}
		return err
	}

	bookings.ID, err = res.LastInsertId()
	return err
}

func (r *BookingsMySQL) Delete(id int64) error {
	res, err := r.db.Exec(
		"DELETE FROM bookings WHERE id = ?",
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

// returns bookings by room id
// sorted by start date
func (r *BookingsMySQL) Get(id int64) ([]pkg.Booking, error) {
	rows, err := r.db.Query(
		"SELECT `id`, `date_start`, `date_end`"+
			"	FROM `bookings` WHERE `room_id` = ?"+
			"	ORDER BY `date_start`",
		id,
	)
	if err != nil {
		return nil, pkg.ErrFailedGet
	}

	bookings := make([]pkg.Booking, 0, 1)
	for rows.Next() {
		b := pkg.Booking{}
		rows.Scan(
			&b.ID,
			&b.Start,
			&b.End,
		)
		bookings = append(bookings, b)
	}

	if len(bookings) == 0 {
		check := true
		row := r.db.QueryRow(
			"SELECT EXISTS (SELECT id FROM room WHERE id = ?)",
			id,
		)

		err = row.Scan(&check)
		if err != nil {
			return nil, pkg.ErrFailedGet
		}
		if !check {
			return nil, pkg.ErrIDNotFound
		}
	}
	return bookings, err
}
