package main

import (
	"log"
	"os"

	"github.com/Avepa/booking/pkg/handler"
	"github.com/Avepa/booking/pkg/repository"
	"github.com/Avepa/booking/pkg/repository/mysql"
	"github.com/Avepa/booking/pkg/server"
	"github.com/Avepa/booking/pkg/service"
)

func main() {
	DBcfg := &mysql.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Username: os.Getenv("DATABASE_USERNAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_DBName"),
	}

	/*
		DBcfg := &mysql.Config{
			Host:     "127.0.0.1",
			Port:     "3306",
			Username: "root",
			Password: "12345",
			DBName:   "mydb",
		}*/

	db, err := mysql.NewMySqlDB(DBcfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	repos := repository.NewRepository(db)
	serveces := service.NewService(repos)
	handlers := handler.NewHandler(serveces)
	err = server.RunHTTPServer(
		os.Getenv("HTTTPSERVER_PORT"),
		handlers.Routes(),
	)
	/*
		err = server.RunHTTPServer(
			"80",
			handlers.Routes(),
		)*/
	if err != nil {
		log.Println(err)
		return
	}
}
