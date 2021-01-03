package service

import (
	"testing"

	"github.com/Avepa/booking/pkg"
	mock_repository "github.com/Avepa/booking/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestBookingsService_Add(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockBookings, room int64, booking *pkg.Booking)

	tests := []struct {
		name          string
		inputID       int64
		inputBooking  pkg.Booking
		mock          mockBehavior
		expected      int64
		expectedError error
	}{
		{
			name:    "OK",
			inputID: 1,
			inputBooking: pkg.Booking{
				ID:    4,
				Start: "2018-02-05",
				End:   "2018-02-07",
			},
			mock: func(r *mock_repository.MockBookings, room int64, booking *pkg.Booking) {
				r.EXPECT().Add(room, booking).Return(nil)
			},
			expected: 4,
		},
		{
			name:    "Date is incorrect",
			inputID: 1,
			inputBooking: pkg.Booking{
				ID:    4,
				Start: "2018.02.05",
				End:   "2018.02.07",
			},
			mock:          func(r *mock_repository.MockBookings, room int64, booking *pkg.Booking) {},
			expected:      0,
			expectedError: pkg.ErrDateIsIncorrect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockBookings(c)
			tt.mock(repo, tt.inputID, &tt.inputBooking)

			services := NewBookingsService(repo)
			id, err := services.Add(tt.inputID, &tt.inputBooking)
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			}
			if id != tt.expected {
				t.Error("incorrect id received: ", id)
			}
		})
	}
}

func TestBookingsService_Get(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockBookings, room int64, booking []pkg.Booking)

	tests := []struct {
		name          string
		input         int64
		mock          mockBehavior
		expected      []pkg.Booking
		expectedError error
	}{
		{

			name:  "OK",
			input: 1,
			mock: func(r *mock_repository.MockBookings, room int64, booking []pkg.Booking) {
				r.EXPECT().Get(room).Return(booking, nil)
			},
			expected: []pkg.Booking{
				{
					ID:    4,
					Start: "2018-02-05",
					End:   "2018-02-07",
				},
				{
					ID:    12,
					Start: "2018-02-15",
					End:   "2018-02-27",
				},
			},
		},
		{
			name:  "Date is incorrect",
			input: 1,
			mock: func(r *mock_repository.MockBookings, room int64, booking []pkg.Booking) {
				r.EXPECT().Get(room).Return(booking, pkg.ErrIDNotFound)
			},
			expectedError: pkg.ErrIDNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockBookings(c)
			tt.mock(repo, tt.input, tt.expected)

			services := NewBookingsService(repo)
			bookings, err := services.Get(tt.input)
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			} else if err != nil {
				return
			}

			for k := range bookings {
				if bookings[k] != tt.expected[k] {
					t.Error("incorrect data received: ", bookings)
					return
				}
			}
		})
	}
}

func TestBookingsService_Delete(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockBookings, booking int64)

	tests := []struct {
		name          string
		input         int64
		mock          mockBehavior
		expectedError error
	}{
		{
			name:  "OK",
			input: 12,
			mock: func(r *mock_repository.MockBookings, booking int64) {
				r.EXPECT().Delete(booking).Return(nil)
			},
		},
		{
			name:  "Failed delete",
			input: 15,
			mock: func(r *mock_repository.MockBookings, booking int64) {
				r.EXPECT().Delete(booking).Return(pkg.ErrFailedDelete)
			},
			expectedError: pkg.ErrFailedDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockBookings(c)
			tt.mock(repo, tt.input)

			services := NewBookingsService(repo)
			err := services.Delete(tt.input)
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			}
		})
	}
}
