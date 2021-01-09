package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Avepa/booking/pkg"
	"github.com/Avepa/booking/pkg/service"
	mock_service "github.com/Avepa/booking/pkg/service/mocks"
	"github.com/golang/mock/gomock"
)

func TestHandler_createBookings(t *testing.T) {
	type mockBehavior func(r *mock_service.MockBookings, booking *pkg.Booking)

	type result struct {
		ID  int64  `json:"booking_id"`
		Err string `json:"error"`
	}

	tests := []struct {
		name                 string
		inputID              string
		inputBooking         *pkg.Booking
		mock                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody result
	}{
		{
			name:    "OK",
			inputID: "1",
			inputBooking: &pkg.Booking{
				Start: "2018.01.05",
				End:   "2018.02.01",
			},
			mock: func(r *mock_service.MockBookings, booking *pkg.Booking) {
				idRoom := int64(1)
				idBooking := int64(1)
				r.EXPECT().Add(idRoom, booking).Return(idBooking, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: result{
				ID: 1,
			},
		},
		{
			name:    "ID not valid",
			inputID: "sg",
			inputBooking: &pkg.Booking{
				Start: "2018.01.05",
				End:   "2018.02.01",
			},
			mock:               func(r *mock_service.MockBookings, booking *pkg.Booking) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: result{
				Err: pkg.ErrIdNotValid.Error(),
			},
		},
		{
			name:    "Date is incorrect",
			inputID: "1",
			inputBooking: &pkg.Booking{
				Start: "2018.21.05",
				End:   "2018.02.01",
			},
			mock: func(r *mock_service.MockBookings, booking *pkg.Booking) {
				idRoom := int64(1)
				idBooking := int64(0)
				r.EXPECT().Add(idRoom, booking).
					Return(idBooking, pkg.ErrDateIsIncorrect)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: result{
				Err: pkg.ErrDateIsIncorrect.Error(),
			},
		},
		{
			name:    "No foreign key",
			inputID: "1",
			inputBooking: &pkg.Booking{
				Start: "2018.01.05",
				End:   "2018.02.01",
			},
			mock: func(r *mock_service.MockBookings, booking *pkg.Booking) {
				idRoom := int64(1)
				idBooking := int64(0)
				r.EXPECT().Add(idRoom, booking).
					Return(idBooking, pkg.ErrNoForeignKey)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: result{
				Err: pkg.ErrNoForeignKey.Error(),
			},
		},
		{
			name:    "Failed get",
			inputID: "1",
			inputBooking: &pkg.Booking{
				Start: "2018.01.05",
				End:   "2018.02.01",
			},
			mock: func(r *mock_service.MockBookings, booking *pkg.Booking) {
				idRoom := int64(1)
				idBooking := int64(0)
				r.EXPECT().Add(idRoom, booking).
					Return(idBooking, pkg.ErrFailedGet)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: result{
				Err: pkg.ErrFailedGet.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookings(c)
			tt.mock(repo, tt.inputBooking)

			services := &service.Service{Bookings: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.createBooking)
			srv := httptest.NewServer(h)
			defer srv.Close()

			req, err := http.NewRequest("POST", srv.URL, nil)
			if err != nil {
				t.Error(err)
				return
			}

			req.Header.Add("room_id", tt.inputID)
			req.Header.Add("date_start", tt.inputBooking.Start)
			req.Header.Add("date_end", tt.inputBooking.End)

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil && err != io.EOF {
				t.Error(err)
				return
			}

			if resp.StatusCode != tt.expectedStatusCode {
				t.Error("wrong error code received: ", resp.StatusCode)
				return
			}

			body := result{}
			json.NewDecoder(resp.Body).Decode(&body)
			if body != tt.expectedResponseBody {
				t.Error("wrong body received: ", body)
			}
		})
	}
}

func TestHandler_getBookings(t *testing.T) {
	type mockBehavior func(r *mock_service.MockBookings, bookings []pkg.Booking)

	tests := []struct {
		name                  string
		input                 string
		mock                  mockBehavior
		expectedStatusCode    int
		expectedResponseBody  []pkg.Booking
		expectedResponseError Error
	}{
		{
			name:  "OK",
			input: "1",
			mock: func(r *mock_service.MockBookings, bookings []pkg.Booking) {
				id := int64(1)
				r.EXPECT().Get(id).Return(bookings, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: []pkg.Booking{
				{
					ID:    65,
					Start: "2018.02.25",
					End:   "2018.03.05",
				},
				{
					ID:    81,
					Start: "2018.04.05",
					End:   "2018.04.13",
				},
				{
					ID:    102,
					Start: "2018.04.20",
					End:   "2018.04.22",
				},
			},
		},
		{
			name:                  "Id not valid",
			input:                 "asf",
			mock:                  func(r *mock_service.MockBookings, bookings []pkg.Booking) {},
			expectedStatusCode:    http.StatusBadRequest,
			expectedResponseError: Error{pkg.ErrIdNotValid.Error()},
		},
		{
			name:  "Failed get",
			input: "1",
			mock: func(r *mock_service.MockBookings, bookings []pkg.Booking) {
				id := int64(1)
				r.EXPECT().Get(id).Return(bookings, pkg.ErrFailedGet)
			},
			expectedStatusCode:    http.StatusInternalServerError,
			expectedResponseError: Error{pkg.ErrFailedGet.Error()},
		},
		{
			name:  "ID not found",
			input: "1",
			mock: func(r *mock_service.MockBookings, bookings []pkg.Booking) {
				id := int64(1)
				r.EXPECT().Get(id).Return(bookings, pkg.ErrIDNotFound)
			},
			expectedStatusCode:    http.StatusBadRequest,
			expectedResponseError: Error{pkg.ErrIDNotFound.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookings(c)
			tt.mock(repo, tt.expectedResponseBody)

			services := &service.Service{Bookings: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.getBookings)
			srv := httptest.NewServer(h)
			defer srv.Close()

			url := fmt.Sprintf("%s/?%s=%s", srv.URL, "room_id", tt.input)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Error(err)
				return
			}

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil && err != io.EOF {
				t.Error(err)
				return
			}

			if resp.StatusCode != tt.expectedStatusCode {
				t.Error("wrong error code received: ", resp.StatusCode)
				return
			} else if resp.StatusCode != http.StatusOK {
				body := Error{}
				json.NewDecoder(resp.Body).Decode(&body)
				if body != tt.expectedResponseError {
					t.Error("wrong error code received: ", resp.StatusCode)
				}
				return
			}

			body := []pkg.Booking{}
			json.NewDecoder(resp.Body).Decode(&body)
			if len(body) != len(tt.expectedResponseBody) {
				t.Error(body)
			}
			for i := range body {
				if body[i] != tt.expectedResponseBody[i] {
					b := fmt.Sprint(body)
					err = errors.New("wrong body received: " + b)
				}
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestHandler_deleteBookings(t *testing.T) {
	type mockBehavior func(r *mock_service.MockBookings, id int64)

	tests := []struct {
		name                 string
		input                int64
		inputBody            string
		mock                 mockBehavior
		expected             Status
		expectedStatusCode   int
		expectedResponseBody Error
	}{
		{
			name:      "OK",
			input:     1,
			inputBody: "1",
			mock: func(r *mock_service.MockBookings, id int64) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expected:           Status{Status: "ok"},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:                 "ID not valid",
			inputBody:            "Adf",
			mock:                 func(r *mock_service.MockBookings, id int64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: Error{pkg.ErrIdNotValid.Error()},
		},
		{
			name:      "Failed delete",
			input:     2,
			inputBody: "2",
			mock: func(r *mock_service.MockBookings, id int64) {
				r.EXPECT().Delete(id).Return(pkg.ErrFailedDelete)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: Error{pkg.ErrFailedDelete.Error()},
		},
		{
			name:      "ID not found",
			input:     3,
			inputBody: "3",
			mock: func(r *mock_service.MockBookings, id int64) {
				r.EXPECT().Delete(id).Return(pkg.ErrIDNotFound)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: Error{pkg.ErrIDNotFound.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookings(c)
			tt.mock(repo, tt.input)

			services := &service.Service{Bookings: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.deleteBookings)
			srv := httptest.NewServer(h)
			defer srv.Close()

			url := fmt.Sprintf("%s/?%s=%s", srv.URL, "booking_id", tt.inputBody)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Error(err)
				return
			}

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil && err != io.EOF {
				t.Error(err)
				return
			}

			if resp.StatusCode != tt.expectedStatusCode {
				t.Error("wrong error code received:", resp.StatusCode)
				return
			} else if resp.StatusCode == http.StatusOK {
				s := Status{}
				json.NewDecoder(resp.Body).Decode(&s)
				if s != tt.expected {
					t.Error("wrong status received")
				}
				return
			}

			body := Error{}
			json.NewDecoder(resp.Body).Decode(&body)
			if body != tt.expectedResponseBody {
				t.Error("wrong body received: ", err)
			}
		})
	}
}
