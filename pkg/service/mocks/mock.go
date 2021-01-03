// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	pkg "github.com/Avepa/booking/pkg"
	gomock "github.com/golang/mock/gomock"
)

// MockRoom is a mock of Room interface.
type MockRoom struct {
	ctrl     *gomock.Controller
	recorder *MockRoomMockRecorder
}

// MockRoomMockRecorder is the mock recorder for MockRoom.
type MockRoomMockRecorder struct {
	mock *MockRoom
}

// NewMockRoom creates a new mock instance.
func NewMockRoom(ctrl *gomock.Controller) *MockRoom {
	mock := &MockRoom{ctrl: ctrl}
	mock.recorder = &MockRoomMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoom) EXPECT() *MockRoomMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockRoom) Add(room *pkg.Room) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", room)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockRoomMockRecorder) Add(room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRoom)(nil).Add), room)
}

// Delete mocks base method.
func (m *MockRoom) Delete(id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRoomMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoom)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockRoom) Get(sort string) ([]pkg.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", sort)
	ret0, _ := ret[0].([]pkg.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRoomMockRecorder) Get(sort interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRoom)(nil).Get), sort)
}

// MockBookings is a mock of Bookings interface.
type MockBookings struct {
	ctrl     *gomock.Controller
	recorder *MockBookingsMockRecorder
}

// MockBookingsMockRecorder is the mock recorder for MockBookings.
type MockBookingsMockRecorder struct {
	mock *MockBookings
}

// NewMockBookings creates a new mock instance.
func NewMockBookings(ctrl *gomock.Controller) *MockBookings {
	mock := &MockBookings{ctrl: ctrl}
	mock.recorder = &MockBookingsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookings) EXPECT() *MockBookingsMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBookings) Add(room int64, booking *pkg.Booking) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", room, booking)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockBookingsMockRecorder) Add(room, booking interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockBookings)(nil).Add), room, booking)
}

// Delete mocks base method.
func (m *MockBookings) Delete(id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBookingsMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBookings)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockBookings) Get(roomID int64) ([]pkg.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", roomID)
	ret0, _ := ret[0].([]pkg.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBookingsMockRecorder) Get(roomID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBookings)(nil).Get), roomID)
}
