package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        user_id SERIAL PRIMARY KEY,
        telegram_id BIGINT UNIQUE NOT NULL,
        username VARCHAR(50),
        fullname VARCHAR(100),
        create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS rooms (
        room_id SERIAL PRIMARY KEY,
        room_name VARCHAR(50) UNIQUE NOT NULL,
        capacity INT DEFAULT 0,
        status VARCHAR(20) DEFAULT 'ACTIVE',
        create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS bookings (
        booking_id SERIAL PRIMARY KEY,
        room_id INT REFERENCES rooms(room_id) ON DELETE CASCADE,
        user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
        topic VARCHAR(200),
        date DATE NOT NULL,
        start_time TIME NOT NULL,
        end_time TIME NOT NULL,
        status VARCHAR(20) DEFAULT 'SUCCESS',
        create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT unique_booking UNIQUE (room_id, date, start_time, end_time)
    );
    CREATE TABLE IF NOT EXISTS participants (
        participant_id SERIAL PRIMARY KEY,
        booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
