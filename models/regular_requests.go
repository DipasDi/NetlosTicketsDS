package models

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func Connectdb() *sql.DB {
	DB, err := sql.Open("mysql", "")
	if err != nil {
		log.Println(err)
	}
	return DB
}
