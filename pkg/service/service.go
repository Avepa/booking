package service

import (
	"github.com/Avepa/booking/pkg"
	"github.com/Avepa/booking/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Room interface {
	Add(room *pkg.Room) (int64, error)
	Delete(id int64) error
	Get(sort string) ([]pkg.Room, error)
}

type Bookings interface {
	Add(room int64, booking *pkg.Booking) (int64, error)
	Delete(id int64) error
	Get(roomID int64) ([]pkg.Booking, error)
}

type Service struct {
	Room
	Bookings
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Room:     NewRoomService(repos.Room),
		Bookings: NewBookingsService(repos.Bookings),
	}
}
