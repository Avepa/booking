package handler

import (
	"database/sql"
	"encoding/json"
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

func TestHandler_addRoom(t *testing.T) {
	type mockBehavior func(r *mock_service.MockRoom, room *pkg.Room)

	type input struct {
		Price       string
		Description string
	}

	type result struct {
		ID  int64  `json:"room_id"`
		Err string `json:"error"`
	}

	tests := []struct {
		name                 string
		input                input
		inputRoom            *pkg.Room
		mock                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody result
	}{
		{
			name: "OK",
			input: input{
				Description: "Good room",
				Price:       "5.41",
			},
			inputRoom: &pkg.Room{
				Description: "Good room",
				Price:       5.41,
			},
			mock: func(r *mock_service.MockRoom, room *pkg.Room) {
				id := int64(1)
				r.EXPECT().Add(room).Return(id, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: result{
				ID: 1,
			},
		},
		{
			name: "Price not valid",
			input: input{
				Description: "Good room",
				Price:       "Ls",
			},
			mock:               func(r *mock_service.MockRoom, room *pkg.Room) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: result{
				Err: pkg.ErrPriceNotValid.Error(),
			},
		},
		{
			name: "Fail save",
			input: input{
				Description: "Good room",
				Price:       "3.55",
			},
			inputRoom: &pkg.Room{
				Description: "Good room",
				Price:       3.55,
			},
			mock: func(r *mock_service.MockRoom, room *pkg.Room) {
				id := int64(0)
				r.EXPECT().Add(room).Return(id, pkg.ErrFailedSave)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: result{
				Err: pkg.ErrFailedSave.Error(),
			},
		},
		{
			name: "Unknown error",
			input: input{
				Description: "Good room",
				Price:       "3.55",
			},
			inputRoom: &pkg.Room{
				Description: "Good room",
				Price:       3.55,
			},
			mock: func(r *mock_service.MockRoom, room *pkg.Room) {
				id := int64(0)
				r.EXPECT().Add(room).Return(id, pkg.ErrDateIsIncorrect)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: result{
				Err: pkg.ErrDateIsIncorrect.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRoom(c)
			tt.mock(repo, tt.inputRoom)

			services := &service.Service{Room: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.addRoom)
			srv := httptest.NewServer(h)
			defer srv.Close()

			req, err := http.NewRequest("POST", srv.URL, nil)
			if err != nil {
				t.Error(err)
				return
			}

			req.Header.Add("description", tt.input.Description)
			req.Header.Add("price", tt.input.Price)

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

func TestHandler_getRoom(t *testing.T) {
	type mockBehavior func(r *mock_service.MockRoom, room []pkg.Room, sort string)

	tests := []struct {
		name                 string
		input                string
		mock                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody []pkg.Room
	}{
		{
			name:  "OK",
			input: "date",
			mock: func(r *mock_service.MockRoom, room []pkg.Room, sort string) {
				r.EXPECT().Get(sort).Return(room, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: []pkg.Room{
				{
					ID:          3,
					Description: "Good",
					Price:       4.99,
					Date:        "2018.01.10",
				},
				{
					ID:          2,
					Description: "Luxury",
					Price:       5.99,
					Date:        "2018.01.08",
				},
				{
					ID:          1,
					Description: "VIP",
					Price:       7.99,
					Date:        "2018.01.06",
				},
			},
		},
		{
			name:  "Internal server error",
			input: "date",
			mock: func(r *mock_service.MockRoom, room []pkg.Room, sort string) {
				r.EXPECT().Get(sort).Return(nil, sql.ErrConnDone)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: []pkg.Room{
				{
					ID:          3,
					Description: "Good",
					Price:       4.99,
					Date:        "2018.01.10",
				},
				{
					ID:          2,
					Description: "Luxury",
					Price:       5.99,
					Date:        "2018.01.08",
				},
				{
					ID:          1,
					Description: "VIP",
					Price:       7.99,
					Date:        "2018.01.06",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRoom(c)
			tt.mock(repo, tt.expectedResponseBody, tt.input)

			services := &service.Service{Room: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.getRoom)
			srv := httptest.NewServer(h)
			defer srv.Close()

			url := fmt.Sprintf("%s/?%s=%s", srv.URL, "sorting", tt.input)
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
				if body.Err != http.StatusText(http.StatusInternalServerError) {
					t.Error("wrong error code received: ", resp.StatusCode)
					return
				}
				return
			}

			body := []pkg.Room{}
			json.NewDecoder(resp.Body).Decode(&body)
			if len(body) != len(tt.expectedResponseBody) {
				t.Error(body)
			}
			for i := range body {
				if body[i] != tt.expectedResponseBody[i] {
					b := fmt.Sprint(body)
					t.Error("wrong body received: ", b)
					return
				}
			}
		})
	}
}

func TestHandler_deleteRoom(t *testing.T) {
	type mockBehavior func(r *mock_service.MockRoom, id int64)

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
			mock: func(r *mock_service.MockRoom, id int64) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expected:           Status{Status: "ok"},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:                 "ID not valid",
			inputBody:            "Adf",
			mock:                 func(r *mock_service.MockRoom, id int64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: Error{pkg.ErrIdNotValid.Error()},
		},
		{
			name:      "Failed delete",
			input:     2,
			inputBody: "2",
			mock: func(r *mock_service.MockRoom, id int64) {
				r.EXPECT().Delete(id).Return(pkg.ErrFailedDelete)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: Error{pkg.ErrFailedDelete.Error()},
		},
		{
			name:      "ID not found",
			input:     3,
			inputBody: "3",
			mock: func(r *mock_service.MockRoom, id int64) {
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

			repo := mock_service.NewMockRoom(c)
			tt.mock(repo, tt.input)

			services := &service.Service{Room: repo}
			handler := Handler{services}
			h := http.HandlerFunc(handler.deleteRoom)
			srv := httptest.NewServer(h)
			defer srv.Close()

			url := fmt.Sprintf("%s/?%s=%s", srv.URL, "room_id", tt.inputBody)
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
