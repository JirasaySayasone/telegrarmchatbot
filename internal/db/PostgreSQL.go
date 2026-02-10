package db

import (
	"database/sql"
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://user:password@localhost:8080/dbname")
}

func createUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255),
		participants TEXT,
		topic TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createRoomTable(db *sql.DB) error {
	query := `
	room_id SERIAL PRIMARY KEY,
	booked_by VARCHAR(100) NOT NULL,
	day INTEGER NOT NULL,
	month INTEGER NOT NULL,
	year INTEGER NOT NULL,
	timestamp_start TIMESTAMP NOT NULL,
	timestamp_end TIMESTAMP NOT NULL,
	booking_no VARCHAR(100) NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
