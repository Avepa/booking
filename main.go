package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Query("SELECT `id`, `date`, `price`, `description` FROM room" +
		" ORDER BY date")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("err")
	var id int64
	var data string
	var price float64
	var desc string
	for rows.Next() {
		rows.Scan(
			&id,
			&data,
			&price,
			&desc,
		)
		fmt.Println(
			id,
			data,
			price,
			desc,
		)
	}

	fmt.Println("end")

	res, err := db.Exec(
		"INSERT INTO bookings (id, room_id, date_start, date_end) VALUES (12, ?, ?, ?)",
		23,
		"2018.01.10",
		"2019.01.17",
	)
	fmt.Println(err, err == mysql.ErrBusyBuffer)
	if err != nil {
		fmt.Println(errors.Is(err, mysql.ErrBusyBuffer))
		fmt.Println(errors.Is(err, mysql.ErrCleartextPassword))
		fmt.Println(errors.Is(err, mysql.ErrMalformPkt))
		fmt.Println(errors.Is(err, mysql.ErrNativePassword))
		fmt.Println(errors.Is(err, mysql.ErrNoTLS))
		fmt.Println(errors.Is(err, mysql.ErrOldPassword))
		fmt.Println(errors.Is(err, mysql.ErrOldProtocol))
		fmt.Println(errors.Is(err, mysql.ErrPktSync))
		fmt.Println(errors.Is(err, mysql.ErrPktSyncMul))
		fmt.Println(errors.Is(err, mysql.ErrPktTooLarge))
		fmt.Println(errors.Is(err, mysql.ErrUnknownPlugin))
		fmt.Println(errors.Is(err, sql.ErrConnDone))
		fmt.Println(errors.Is(err, sql.ErrNoRows))
		fmt.Println(errors.Is(err, sql.ErrTxDone))

	}
	fmt.Println(res.LastInsertId, res.RowsAffected)
	/*
			http.HandleFunc("/", A)
			http.ListenAndServe(":80", nil)
		}

		func A(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Header.Get("sdg"))
			fmt.Println(r.URL.Query().Get("sdg"))*/
}
