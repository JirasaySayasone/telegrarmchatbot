package model

import "time"

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
type User struct {
	UserID       int      `json:"user_id"`
	Booker       string   `json:"booker"`
	Participants []string `json:"participants"`
	Topic        string   `json:"topic"`
}

type Room struct {
	RoomID          int       `json:"room_id"`
	Booked_by       string    `json:"booked_by"`
	Day             int       `json:"day"`
	Month           int       `json:"month"`
	Year            int       `json:"year"`
	Timestamp_start time.Time `json:"timestamp_start"`
	Timestamp_end   time.Time `json:"timestamp_end"`
	Booking_no      string    `json:"booking_no"`
	Create_at       time.Time `json:"create_at"`
}
