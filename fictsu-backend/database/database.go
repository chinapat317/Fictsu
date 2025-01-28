package database

import (
	"log"
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connection() {
	var err error
	DB, err = sql.Open("postgres", "postgres://kwang:fictsu@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}
