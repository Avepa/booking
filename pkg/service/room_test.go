package service

import (
	"testing"

	"github.com/Avepa/booking/pkg"
	mock_repository "github.com/Avepa/booking/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestRoomService_Add(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRoom, room *pkg.Room)

	tests := []struct {
		name          string
		input         pkg.Room
		mock          mockBehavior
		expectedID    int64
		expectedError error
	}{
		{
			name: "OK",
			input: pkg.Room{
				ID:          54,
				Description: "Good",
				Price:       5.14,
			},
			mock: func(r *mock_repository.MockRoom, room *pkg.Room) {
				r.EXPECT().Add(room).Return(nil)
			},
			expectedID: 54,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRoom(c)
			tt.mock(repo, &tt.input)

			services := NewRoomService(repo)
			id, err := services.Add(&tt.input)
			if id != tt.expectedID {
				t.Error("incorrect id received: ", id)
			}
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			}
		})
	}
}

func TestRoomService_Get(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRoom, room []pkg.Room)

	tests := []struct {
		name          string
		input         string
		mock          mockBehavior
		expected      []pkg.Room
		expectedError error
	}{
		{
			name:  "OK date",
			input: "date",
			mock: func(r *mock_repository.MockRoom, room []pkg.Room) {
				r.EXPECT().GetByDate().Return(room, nil)
			},
			expected: []pkg.Room{
				{
					ID:          1,
					Description: "VIP",
					Price:       10.0,
					Date:        "2018.01.01",
				},
				{
					ID:          2,
					Description: "Good",
					Price:       7.0,
					Date:        "2018.01.02",
				},
			},
		},
		{
			name:  "OK price",
			input: "price",
			mock: func(r *mock_repository.MockRoom, room []pkg.Room) {
				r.EXPECT().GetByPrice().Return(room, nil)
			},
			expected: []pkg.Room{
				{
					ID:          1,
					Description: "VIP",
					Price:       10.0,
					Date:        "2018.01.01",
				},
				{
					ID:          2,
					Description: "Good",
					Price:       7.0,
					Date:        "2018.01.02",
				},
			},
		},
		{
			name:  "OK price desc",
			input: "price_desc",
			mock: func(r *mock_repository.MockRoom, room []pkg.Room) {
				r.EXPECT().GetByPriceDESC().Return(room, nil)
			},
			expected: []pkg.Room{
				{
					ID:          1,
					Description: "VIP",
					Price:       10.0,
					Date:        "2018.01.01",
				},
				{
					ID:          2,
					Description: "Good",
					Price:       7.0,
					Date:        "2018.01.02",
				},
			},
		},
		{
			name:  "OK date desc",
			input: "date_desc",
			mock: func(r *mock_repository.MockRoom, room []pkg.Room) {
				r.EXPECT().GetByDateDESC().Return(room, nil)
			},
			expected: []pkg.Room{
				{
					ID:          1,
					Description: "VIP",
					Price:       10.0,
					Date:        "2018.01.01",
				},
				{
					ID:          2,
					Description: "Good",
					Price:       7.0,
					Date:        "2018.01.02",
				},
			},
		},
		{
			name:  "OK date desc",
			input: "dsg",
			mock: func(r *mock_repository.MockRoom, room []pkg.Room) {
				r.EXPECT().GetByDateDESC().Return(room, nil)
			},
			expected: []pkg.Room{
				{
					ID:          1,
					Description: "VIP",
					Price:       10.0,
					Date:        "2018.01.01",
				},
				{
					ID:          2,
					Description: "Good",
					Price:       7.0,
					Date:        "2018.01.02",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRoom(c)
			tt.mock(repo, tt.expected)

			services := NewRoomService(repo)
			room, err := services.Get(tt.input)
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			} else if err != nil {
				return
			}

			for k := range room {
				if room[k] != tt.expected[k] {
					t.Error("incorrect data received: ", room)
					return
				}
			}
		})
	}
}

func TestRoomService_Delete(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRoom, room int64)

	tests := []struct {
		name          string
		input         int64
		mock          mockBehavior
		expectedError error
	}{
		{
			name:  "OK",
			input: 1,
			mock: func(r *mock_repository.MockRoom, room int64) {
				r.EXPECT().Delete(room).Return(nil)
			},
		},
		{
			name:  "Failed delete",
			input: 1,
			mock: func(r *mock_repository.MockRoom, room int64) {
				r.EXPECT().Delete(room).Return(pkg.ErrFailedDelete)
			},
			expectedError: pkg.ErrFailedDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRoom(c)
			tt.mock(repo, tt.input)

			services := NewRoomService(repo)
			err := services.Delete(tt.input)
			if err != tt.expectedError {
				t.Error("incorrect error received: ", err)
			}
		})
	}
}
