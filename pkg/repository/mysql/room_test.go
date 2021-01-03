package mysql

import (
	"database/sql"
	"testing"

	"github.com/Avepa/booking/pkg"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoomMySQL_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewRoomMySQL(db)

	tests := []struct {
		name    string
		input   *pkg.Room
		mock    func()
		want    int
		wantErr error
	}{
		{
			name: "OK_1",
			input: &pkg.Room{
				Description: "GOOD",
				Price:       54.54,
			},
			mock: func() {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO room").
					WithArgs("GOOD", 54.54).WillReturnResult(result)
			},
			want: 1,
		},
		{
			name: "OK_2",
			input: &pkg.Room{
				Price: 127,
			},
			mock: func() {
				result := sqlmock.NewResult(2, 1)
				mock.ExpectExec("INSERT INTO room").
					WithArgs("", 127.0).WillReturnResult(result)
			},
			want: 2,
		},
		{
			name: "Failed Save",
			input: &pkg.Room{
				Description: "",
				Price:       240.4,
			},
			mock: func() {
				result := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO room").
					WithArgs("", 240).WillReturnResult(result)
			},
			wantErr: pkg.ErrFailedSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err = r.Add(tt.input)
			if err != tt.wantErr {
				t.Error(err, tt.wantErr)
			}
		})
	}
}

func TestRoomMySQL_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewRoomMySQL(db)

	tests := []struct {
		name    string
		input   int64
		mock    func()
		wantErr error
	}{
		{
			name:  "OK",
			input: 1,
			mock: func() {
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM `room` WHERE (.+)").
					WithArgs(1).WillReturnResult(result)
			},
		},
		{
			name:  "Not Found",
			input: 2,
			mock: func() {
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("DELETE FROM `room` WHERE (.+)").
					WithArgs(2).WillReturnResult(result)
			},
			wantErr: pkg.ErrIDNotFound,
		},
		{
			name:  "Failed Delete 1",
			input: 3,
			mock: func() {
				mock.ExpectExec("DELETE FROM `room` WHERE (.+)").
					WithArgs(3).WillReturnError(sql.ErrConnDone)
			},
			wantErr: pkg.ErrFailedDelete,
		},
		{
			name:  "Failed Delete 2",
			input: 4,
			mock: func() {
				result := sqlmock.NewErrorResult(sql.ErrConnDone)
				mock.ExpectExec("DELETE FROM `room` WHERE (.+)").
					WithArgs(4).WillReturnResult(result)
			},
			wantErr: pkg.ErrFailedDelete,
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

func TestRoomMySQL_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewRoomMySQL(db)

	tests := []struct {
		name    string
		sort    func() ([]pkg.Room, error)
		mock    func()
		want    []pkg.Room
		wantErr error
	}{
		{
			name: "OK Func GetByDate()",
			sort: r.GetByDate,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "price", "description"}).
					AddRow(1, "2018.01.03", 3.54, "Good room").
					AddRow(2, "2018.03.06", 5.03, "VIP ROOM").
					AddRow(3, "2019.10.03", 10, "")

				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY date").
					WillReturnRows(rows)
			},
			want: []pkg.Room{
				{
					ID:          1,
					Price:       3.54,
					Date:        "2018.01.03",
					Description: "Good room",
				},
				{
					ID:          2,
					Price:       5.03,
					Date:        "2018.03.06",
					Description: "VIP ROOM",
				},
				{
					ID:          3,
					Price:       10.0,
					Date:        "2019.10.03",
					Description: "",
				},
			},
		},
		{
			name: "Conn done Func GetByDate()",
			sort: r.GetByDate,
			mock: func() {
				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY date").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: sql.ErrConnDone,
		},
		{
			name: "OK Func GetByDateDESC()",
			sort: r.GetByDateDESC,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "price", "description"}).
					AddRow(3, "2019.10.03", 10, "").
					AddRow(2, "2018.03.06", 5.03, "VIP ROOM").
					AddRow(1, "2018.01.03", 3.54, "Good room")

				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY date DESC").
					WillReturnRows(rows)
			},
			want: []pkg.Room{
				{
					ID:          3,
					Price:       10.0,
					Date:        "2019.10.03",
					Description: "",
				},
				{
					ID:          2,
					Price:       5.03,
					Date:        "2018.03.06",
					Description: "VIP ROOM",
				},
				{
					ID:          1,
					Price:       3.54,
					Date:        "2018.01.03",
					Description: "Good room",
				},
			},
		},
		{
			name: "Conn done Func GetByDate()",
			sort: r.GetByDateDESC,
			mock: func() {
				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY date DESC").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: sql.ErrConnDone,
		},

		{
			name: "OK Func GetByPrice()",
			sort: r.GetByPrice,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "price", "description"}).
					AddRow(1, "2018.01.03", 3.54, "Good room").
					AddRow(2, "2018.03.06", 5.03, "VIP ROOM").
					AddRow(3, "2019.10.03", 10, "")

				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY price").
					WillReturnRows(rows)
			},
			want: []pkg.Room{
				{
					ID:          1,
					Price:       3.54,
					Date:        "2018.01.03",
					Description: "Good room",
				},
				{
					ID:          2,
					Price:       5.03,
					Date:        "2018.03.06",
					Description: "VIP ROOM",
				},
				{
					ID:          3,
					Price:       10.0,
					Date:        "2019.10.03",
					Description: "",
				},
			},
		},
		{
			name: "Conn done Func GetByPrice()",
			sort: r.GetByPrice,
			mock: func() {
				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY price").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: sql.ErrConnDone,
		},
		{
			name: "OK Func GetByPriceDESC()",
			sort: r.GetByPriceDESC,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "price", "description"}).
					AddRow(3, "2019.10.03", 10, "").
					AddRow(2, "2018.03.06", 5.03, "VIP ROOM").
					AddRow(1, "2018.01.03", 3.54, "Good room")

				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY price DESC").
					WillReturnRows(rows)
			},
			want: []pkg.Room{
				{
					ID:          3,
					Price:       10.0,
					Date:        "2019.10.03",
					Description: "",
				},
				{
					ID:          2,
					Price:       5.03,
					Date:        "2018.03.06",
					Description: "VIP ROOM",
				},
				{
					ID:          1,
					Price:       3.54,
					Date:        "2018.01.03",
					Description: "Good room",
				},
			},
		},
		{
			name: "Conn done Func GetByPriceDESC()",
			sort: r.GetByPriceDESC,
			mock: func() {
				mock.ExpectQuery("SELECT `id`, `date`, `price`, `description` FROM room" +
					" ORDER BY price DESC").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			room, err := tt.sort()
			if err != tt.wantErr {
				t.Error(err)
			} else if err == nil {
				for i := range tt.want {
					if tt.want[i] != room[i] {
						t.Error("array sorted incorrectly")
					}
				}
			}
		})
	}
}
