package db

import (
	"database/sql"
	"telegrarmchatbot/internal/model"
	"time"
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
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func createRoomTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS rooms (
	room_id SERIAL PRIMARY KEY,
	booked_by VARCHAR(100) NOT NULL,
	day INTEGER NOT NULL,
	month INTEGER NOT NULL,
	year INTEGER NOT NULL,
	timestamp_start TIMESTAMP NOT NULL,
	timestamp_end TIMESTAMP NOT NULL,
	booking_no VARCHAR(100) NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func getUserbyName(db *sql.DB, booker string) ([]model.User, error) {
	query := `SELECT  
	id, username, participants, topic, created_at
	FROM users WHERE Booker = $1
	ORDER BY created_at DESC`
	rows, err := db.Query(query, booker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Users []model.User
	for rows.Next() {
		var User model.User
		err := rows.Scan(&User.UserID, &User.Booker, &User.Participants, &User.Topic)
		if err != nil {
			return nil, err
		}
		Users = append(Users, User)
	}
	return Users, nil
}

func GetRoomsByBooker(db *sql.DB, bookedBy string) ([]model.Room, error) {
	query := `
	SELECT room_id, booked_by, day, month, year, timestamp_start, timestamp_end, booking_no, created_at
	FROM rooms 
	WHERE booked_by = $1 
	ORDER BY created_at DESC`

	rows, err := db.Query(query, bookedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err := rows.Scan(&room.RoomID, &room.Booked_by, &room.Day, &room.Month, &room.Year,
			&room.Timestamp_start, &room.Timestamp_end, &room.Booking_no, &room.Create_at)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
func GetRoomByTime(db *sql.DB, startTime time.Time, endTime time.Time) ([]model.Room, error) {
	query := `
	SELECT room_id, booked_by, day, month, year, timestamp_start, timestamp_end, booking_no, created_at
	FROM rooms 
	WHERE timestamp_start >= $1 AND timestamp_end <= $2
	ORDER BY timestamp_start ASC`

	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err := rows.Scan(&room.RoomID, &room.Booked_by, &room.Day, &room.Month, &room.Year,
			&room.Timestamp_start, &room.Timestamp_end, &room.Booking_no, &room.Create_at)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
