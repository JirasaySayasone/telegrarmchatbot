package model

import "time"

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type User struct {
	UserID     int       `json:"user_id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FullName   string    `json:"fullname"` // Changed: fullname not full_name
	CreateAt   time.Time `json:"create_at"`
}

type Room struct {
	RoomID   int       `json:"room_id"`
	RoomName string    `json:"room_name"`
	Capacity int       `json:"capacity"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`

	// Booked_by       string    `json:"booked_by"`

	// Day             int       `json:"day"`
	// Month           int       `json:"month"`
	// Year            int       `json:"year"`
	// Timestamp_start time.Time `json:"timestamp_start"`
	// Timestamp_end   time.Time `json:"timestamp_end"`
	// Booking_no      string    `json:"booking_no"`

}

type Booking struct {
	BookingID int       `json:"booking_id"`
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	Topic     string    `json:"topic"`
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	CreateAt  time.Time `json:"create_at"`

	//join fields
	RoomName     string   `json:"room_name,omitempty"`
	Username     string   `json:"username,omitempty"`
	FullName     string   `json:"fullname,omitempty"`
	Participants []string `json:"participants,omitempty"`
}

type Participants struct {
	ParticipantID int    `json:"participant_id"`
	BookingID     int    `json:"booking_id"`
	Name          string `json:"name"`
}

type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsFree    bool      `json:"is_free"`
	Booking   *Booking  `json:"booking,omitempty"`
}

type RoomSchedule struct {
	RoomID    int        `json:"room_id"`
	RoomName  string     `json:"room_name"`
	Date      time.Time  `json:"date"`
	TimeSlots []TimeSlot `json:"time_slots"`
}
