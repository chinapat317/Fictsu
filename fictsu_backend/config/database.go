package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	var dsn string = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DbHost,
		DbUser,
		DbPw,
		DbName,
		DbPort,
	)
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	DB = db
}
