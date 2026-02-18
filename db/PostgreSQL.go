// db/postgres.go

package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// Update with your actual credentials
	connStr := "host=localhost port=5432 user=acia password=room_bot1234 dbname=room_booking_bot sslmode=disable"
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	
	// Test connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

// SeedRooms inserts the 3 rooms if they don't exist
func SeedRooms(db *sql.DB) error {
	// First, add status column if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE rooms 
		ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'ACTIVE'
	`)
	if err != nil {
		return err
	}
	
	// Insert rooms
	rooms := []string{"Room A", "Room B", "Room C"}
	
	for _, roomName := range rooms {
		query := `
		INSERT INTO rooms (room_name, capacity, status) 
		VALUES ($1, 10, 'ACTIVE')
		ON CONFLICT (room_name) DO NOTHING`
		
		_, err := db.Exec(query, roomName)
		if err != nil {
			return err
		}
	}
	
	return nil
}