package repository

import (
	"database/sql"

	"github.com/Avepa/booking/pkg"
	"github.com/Avepa/booking/pkg/repository/mysql"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Room interface {
	Add(room *pkg.Room) error
	Delete(id int64) error
	GetByDate() ([]pkg.Room, error)
	GetByPrice() ([]pkg.Room, error)
	GetByDateDESC() ([]pkg.Room, error)
	GetByPriceDESC() ([]pkg.Room, error)
}
type Bookings interface {
	Add(room int64, bookings *pkg.Booking) error
	Delete(id int64) error
	Get(id int64) ([]pkg.Booking, error)
}

type Repository struct {
	Room
	Bookings
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Room:     mysql.NewRoomMySQL(db),
		Bookings: mysql.NewBookingsMySQL(db),
	}
}
