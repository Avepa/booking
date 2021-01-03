package mysql

import (
	"database/sql"
	"testing"

	"github.com/Avepa/booking/pkg"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestBookingsMySQL_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewBookingsMySQL(db)

	tests := []struct {
		name          string
		mock          func()
		inputID       int64
		inputBookings *pkg.Booking
		want          int64
		wantErr       error
	}{
		{
			name: "OK",
			mock: func() {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO bookings").
					WithArgs(3, "2018-02-03", "2018-02-10").
					WillReturnResult(result)
			},
			inputID: 3,
			inputBookings: &pkg.Booking{
				Start: "2018-02-03",
				End:   "2018-02-10",
			},
			want: 1,
		},
		{
			name: "Conn Done",
			mock: func() {
				mock.ExpectExec("INSERT INTO bookings").
					WithArgs(3, "2018-02-03", "2018-02-10").
					WillReturnError(sql.ErrConnDone)
			},
			inputID: 3,
			inputBookings: &pkg.Booking{
				Start: "2018-02-03",
				End:   "2018-02-10",
			},
			wantErr: sql.ErrConnDone,
		},
		{
			name: "No Foreign Key",
			mock: func() {
				mock.ExpectExec("INSERT INTO bookings").
					WithArgs(3, "2018-02-03", "2018-02-10").
					WillReturnError(pkg.ErrNoForeignKey)
			},
			inputID: 3,
			inputBookings: &pkg.Booking{
				Start: "2018-02-03",
				End:   "2018-02-10",
			},
			wantErr: pkg.ErrNoForeignKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err = r.Add(tt.inputID, tt.inputBookings)
			if err != tt.wantErr {
				t.Error(err)
			} else if err == nil {
				if tt.want != tt.inputBookings.ID {
					t.Error("wrong id received")
				}
			}
		})
	}
}

func TestBookingsMySQL_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewBookingsMySQL(db)

	tests := []struct {
		name    string
		mock    func()
		input   int64
		wantErr error
	}{
		{
			name: "OK",
			mock: func() {
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM bookings WHERE (.+)").
					WithArgs(1).WillReturnResult(result)
			},
			input: 1,
		},
		{
			name: "Failed Delete 1",
			mock: func() {
				mock.ExpectExec("DELETE FROM bookings WHERE (.+)").
					WithArgs(2).WillReturnError(sql.ErrConnDone)
			},
			input:   2,
			wantErr: pkg.ErrFailedDelete,
		},
		{
			name: "Failed Delete 2",
			mock: func() {
				result := sqlmock.NewErrorResult(sql.ErrConnDone)
				mock.ExpectExec("DELETE FROM bookings WHERE (.+)").
					WithArgs(3).WillReturnResult(result)
			},
			input:   3,
			wantErr: pkg.ErrFailedDelete,
		},
		{
			name: "Not Found",
			mock: func() {
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("DELETE FROM bookings WHERE (.+)").
					WithArgs(4).WillReturnResult(result)
			},
			input:   4,
			wantErr: pkg.ErrIDNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err = r.Delete(tt.input)
			if err != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestBookingsMySQL_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewBookingsMySQL(db)

	tests := []struct {
		name    string
		input   int64
		mock    func()
		want    []pkg.Booking
		wantErr error
	}{
		{
			name:  "OK",
			input: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "date_start", "date_end"}).
					AddRow(4, "2018-03-06", "2018-03-08").
					AddRow(10, "2018-10-01", "2018-11-06").
					AddRow(1, "2019-02-20", "2019-03-06")

				mock.ExpectQuery(
					"SELECT `id`, `date_start`, `date_end`" +
						"	FOR `booking` WHERE `room_id = (.+)" +
						"	ORDER BY `date_start`",
				).WithArgs(1).WillReturnRows(rows)
			},
			want: []pkg.Booking{
				{
					ID:    4,
					Start: "2018-03-06",
					End:   "2018-03-08",
				},
				{
					ID:    10,
					Start: "2018-10-01",
					End:   "2018-11-06",
				},
				{
					ID:    1,
					Start: "2019-02-20",
					End:   "2019-03-06",
				},
			},
		},
		{
			name:  "Failed",
			input: 2,
			mock: func() {
				mock.ExpectQuery(
					"SELECT `id`, `date_start`, `date_end`" +
						"	FOR `booking` WHERE `room_id = (.+)" +
						"	ORDER BY `date_start`",
				).WithArgs(2).WillReturnError(sql.ErrConnDone)
			},
			wantErr: pkg.ErrFailedGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			booking, err := r.Get(tt.input)
			if err != tt.wantErr {
				t.Error(err)
			} else if err == nil {
				for i := range tt.want {
					if tt.want[i] != booking[i] {
						t.Error("array sorted incorrectly")
					}
				}
			}
		})
	}
}
