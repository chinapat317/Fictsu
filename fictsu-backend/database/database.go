package database

import (
	"log"
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connection() {
	var err error
	DB, err = sql.Open("postgres", "postgres://kwang:fictsu@db:5432/fictsu?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Ping to verify connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	log.Println("Connected to the database successfully")
}

func CloseConnection() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
