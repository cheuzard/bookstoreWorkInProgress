package main

import (
	"bookstore/pkg/config"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type App struct {
	db *sql.DB
}

func main() {
	conf, err := config.LoadConfig("DB_config.json")
	if err != nil {
		_ = errors.New("Error loading config" + err.Error())
	}
	log.Printf("conf loader:%v\n", conf.DSN)

	conn, err := sql.Open("mysql", conf.DSN)
	if err != nil {
		log.Fatalf("error opening db connection :%v\n", err)
	}
	app := App{db: conn}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(app.db)

	http.Handle("/", app.OpenConnectionHandler())
	log.Printf("server starting on port: 8080 \n")
	println()
	println()
	println()
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
