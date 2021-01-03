package service

import (
	"time"

	"github.com/Avepa/booking/pkg"
	"github.com/Avepa/booking/pkg/repository"
)

const form = "2006-01-02"

type BookingsService struct {
	repo repository.Bookings
}

func NewBookingsService(repo repository.Bookings) *BookingsService {
	return &BookingsService{repo: repo}
}

func (s *BookingsService) Add(id int64, booking *pkg.Booking) (int64, error) {
	_, err := time.Parse(form, booking.Start)
	if err != nil {
		return 0, pkg.ErrDateIsIncorrect
	}

	_, err = time.Parse(form, booking.End)
	if err != nil {
		return 0, pkg.ErrDateIsIncorrect
	}

	err = s.repo.Add(id, booking)
	return booking.ID, err
}

func (s *BookingsService) Get(roomID int64) ([]pkg.Booking, error) {
	return s.repo.Get(roomID)
}

func (s *BookingsService) Delete(id int64) error {
	return s.repo.Delete(id)
}
