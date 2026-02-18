// db/booking.go - NEW FILE

package db

import (
	"database/sql"
	"fmt"

	"telegrarmchatbot/internal/model"
	"time"
)

// GetBookingsByDate retrieves all bookings for a specific date
func GetBookingsByDate(db *sql.DB, date time.Time) ([]model.Booking, error) {
	query := `
	SELECT b.booking_id, b.room_id, b.user_id, b.topic, b.date, 
	       b.start_time, b.end_time, b.status, b.create_at,
	       r.room_name, u.username, u.fullname
	FROM bookings b
	JOIN rooms r ON b.room_id = r.room_id
	JOIN users u ON b.user_id = u.user_id
	WHERE b.date = $1 AND b.status = 'SUCCESS'
	ORDER BY r.room_name, b.start_time`

	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.BookingID, &booking.RoomID, &booking.UserID, &booking.Topic,
			&booking.Date, &booking.StartTime, &booking.EndTime, &booking.Status,
			&booking.CreateAt, &booking.RoomName, &booking.Username, &booking.FullName,
		)
		if err != nil {
			return nil, err
		}

		// Get participants for this booking
		participants, _ := GetParticipantsByBookingID(db, booking.BookingID)
		booking.Participants = participants

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// CheckTimeConflict checks if a time slot is already booked
func CheckTimeConflict(db *sql.DB, roomID int, date time.Time, startTime, endTime string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM bookings 
	WHERE room_id = $1 
	AND date = $2 
	AND status = 'SUCCESS'
	AND (
		(start_time < $4 AND end_time > $3) OR
		(start_time >= $3 AND start_time < $4)
	)`

	var count int
	err := db.QueryRow(query, roomID, date, startTime, endTime).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateBooking creates a new booking with participants
func CreateBooking(db *sql.DB, booking *model.Booking, participants []string) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert booking
	query := `
	INSERT INTO bookings (room_id, user_id, topic, date, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5, $6, 'SUCCESS')
	RETURNING booking_id, create_at`

	err = tx.QueryRow(
		query,
		booking.RoomID, booking.UserID, booking.Topic,
		booking.Date, booking.StartTime, booking.EndTime,
	).Scan(&booking.BookingID, &booking.CreateAt)

	if err != nil {
		return err
	}

	// Insert participants
	if len(participants) > 0 {
		participantQuery := `INSERT INTO participants (booking_id, name) VALUES ($1, $2)`
		for _, name := range participants {
			_, err = tx.Exec(participantQuery, booking.BookingID, name)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	return tx.Commit()
}

// GetUserBookings retrieves all active bookings for a user
func GetUserBookings(db *sql.DB, userID int) ([]model.Booking, error) {
	query := `
	SELECT b.booking_id, b.room_id, b.user_id, b.topic, b.date, 
	       b.start_time, b.end_time, b.status, b.create_at,
	       r.room_name
	FROM bookings b
	JOIN rooms r ON b.room_id = r.room_id
	WHERE b.user_id = $1 
	AND b.status = 'SUCCESS'
	AND b.date >= CURRENT_DATE
	ORDER BY b.date, b.start_time`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.BookingID, &booking.RoomID, &booking.UserID, &booking.Topic,
			&booking.Date, &booking.StartTime, &booking.EndTime, &booking.Status,
			&booking.CreateAt, &booking.RoomName,
		)
		if err != nil {
			return nil, err
		}

		// Get participants
		participants, _ := GetParticipantsByBookingID(db, booking.BookingID)
		booking.Participants = participants

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// CancelBooking marks a booking as cancelled
func CancelBooking(db *sql.DB, bookingID int, userID int) error {
	query := `
	UPDATE bookings 
	SET status = 'CANCELLED' 
	WHERE booking_id = $1 AND user_id = $2 AND status = 'SUCCESS'`

	result, err := db.Exec(query, bookingID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found or already cancelled")
	}

	return nil
}

// GetParticipantsByBookingID retrieves all participants for a booking
func GetParticipantsByBookingID(db *sql.DB, bookingID int) ([]string, error) {
	query := `SELECT name FROM participants WHERE booking_id = $1`

	rows, err := db.Query(query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		participants = append(participants, name)
	}

	return participants, nil
}

// GetBookingByID retrieves a single booking with all details
func GetBookingByID(db *sql.DB, bookingID int) (*model.Booking, error) {
	query := `
	SELECT b.booking_id, b.room_id, b.user_id, b.topic, b.date, 
	       b.start_time, b.end_time, b.status, b.create_at,
	       r.room_name, u.username, u.fullname
	FROM bookings b
	JOIN rooms r ON b.room_id = r.room_id
	JOIN users u ON b.user_id = u.user_id
	WHERE b.booking_id = $1`

	var booking model.Booking
	err := db.QueryRow(query, bookingID).Scan(
		&booking.BookingID, &booking.RoomID, &booking.UserID, &booking.Topic,
		&booking.Date, &booking.StartTime, &booking.EndTime, &booking.Status,
		&booking.CreateAt, &booking.RoomName, &booking.Username, &booking.FullName,
	)

	if err != nil {
		return nil, err
	}
 
	// Get participants
	participants, _ := GetParticipantsByBookingID(db, booking.BookingID)
	booking.Participants = participants

	return &booking, nil
}
