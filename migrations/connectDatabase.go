package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "root"
)

func Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, _ = sql.Open("postgres", psqlInfo)

	err := DB.Ping()
	if err != nil {
		log.Fatal("Fatal connect to database!")
	}
}

func Disconnect() {
	DB.Close()
}
