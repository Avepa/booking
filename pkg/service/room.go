package service

import (
	"github.com/Avepa/booking/pkg"
	"github.com/Avepa/booking/pkg/repository"
)

type RoomService struct {
	repo repository.Room
}

func NewRoomService(repo repository.Room) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Add(room *pkg.Room) (int64, error) {
	if room.Price < 0.0 {
		return 0, pkg.ErrPriceNotValid
	}

	err := s.repo.Add(room)
	return room.ID, err
}

func (s *RoomService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *RoomService) Get(sort string) ([]pkg.Room, error) {
	switch sort {
	case "date":
		return s.repo.GetByDate()
	case "price":
		return s.repo.GetByPrice()
	case "price_desc":
		return s.repo.GetByPriceDESC()
	default:
		return s.repo.GetByDateDESC()
	}
}
