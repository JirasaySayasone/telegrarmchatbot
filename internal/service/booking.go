// internal/service/booking.go - UPDATE

package service

import (
	"database/sql"
	"fmt"
	"strings"
	"telegrarmchatbot/db"
	"telegrarmchatbot/internal/config"
	"telegrarmchatbot/internal/model"
	"time"
)

type BookingService struct {
	DB *sql.DB
}

func NewBookingService(database *sql.DB) *BookingService {
	return &BookingService{DB: database}
}

// GenerateTodayTimetable creates a full schedule for all rooms for today
func (s *BookingService) GenerateTodayTimetable() ([]model.RoomSchedule, error) {
	today := time.Now()
	return s.GenerateTimetableForDate(today)
}

// GenerateTimetableForDate creates schedule for a specific date
func (s *BookingService) GenerateTimetableForDate(date time.Time) ([]model.RoomSchedule, error) {
	// Get all active rooms
	rooms, err := db.GetAllActiveRooms(s.DB)
	if err != nil {
		return nil, err
	}

	// Get all bookings for this date
	bookings, err := db.GetBookingsByDate(s.DB, date)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup: roomID -> startTime -> booking
	bookingMap := make(map[int]map[string]*model.Booking)
	for _, room := range rooms {
		bookingMap[room.RoomID] = make(map[string]*model.Booking)
	}

	for i := range bookings {
		booking := &bookings[i]
		bookingMap[booking.RoomID][booking.StartTime.Format("15:04")] = booking
	}

	// Generate schedules for all rooms
	var schedules []model.RoomSchedule
	timeSlots := config.GenerateTimeSlots()

	for _, room := range rooms {
		schedule := model.RoomSchedule{
			RoomID:    room.RoomID,
			RoomName:  room.RoomName,
			Date:      date,
			TimeSlots: []model.TimeSlot{},
		}

		for i := 0; i < len(timeSlots)-1; i++ {
			startStr := timeSlots[i]
			endStr := timeSlots[i+1]

			// Parse time strings (e.g. "09:00") and attach the provided date
			tStart, err := time.Parse("15:04", startStr)
			if err != nil {
				return nil, err
			}
			tEnd, err := time.Parse("15:04", endStr)
			if err != nil {
				return nil, err
			}

			startTime := time.Date(date.Year(), date.Month(), date.Day(), tStart.Hour(), tStart.Minute(), 0, 0, date.Location())
			endTime := time.Date(date.Year(), date.Month(), date.Day(), tEnd.Hour(), tEnd.Minute(), 0, 0, date.Location())

			slot := model.TimeSlot{
				StartTime: startTime,
				EndTime:   endTime,
			}

			// Check if this slot is booked (bookingMap keys are in "15:04" format)
			if booking, exists := bookingMap[room.RoomID][startTime.Format("15:04")]; exists {
				slot.IsFree = false
				slot.Booking = booking
			} else {
				slot.IsFree = true
			}

			schedule.TimeSlots = append(schedule.TimeSlots, slot)
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// FormatTimetableMessage converts schedules to Telegram message
func (s *BookingService) FormatTimetableMessage(schedules []model.RoomSchedule) string {
	if len(schedules) == 0 {
		return "No schedule available."
	}

	date := schedules[0].Date.Format("02 Jan 2006")
	message := fmt.Sprintf("📅 *Room Schedule - %s*\n\n", date)

	for _, schedule := range schedules {
		message += fmt.Sprintf("🏢 *%s*\n", schedule.RoomName)

		for _, slot := range schedule.TimeSlots {
			if slot.IsFree {
				message += fmt.Sprintf("  ✅ %s-%s FREE\n", slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
			} else {
				booking := slot.Booking
				message += fmt.Sprintf("  ❌ %s-%s BOOKED\n", slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
				message += fmt.Sprintf("     👤 By: %s\n", booking.FullName)
				message += fmt.Sprintf("     📝 %s\n", booking.Topic)
				if len(booking.Participants) > 0 {
					message += fmt.Sprintf("     👥 %s\n", strings.Join(booking.Participants, ", "))
				}
			}
		}
		message += "\n"
	}

	return message
}

// FormatUserBookings formats a user's bookings into a message
func (s *BookingService) FormatUserBookings(bookings []model.Booking) string {
	if len(bookings) == 0 {
		return "You have no active bookings."
	}

	message := "*Your Bookings:*\n\n"
	for i, booking := range bookings {
		message += fmt.Sprintf("%d. 🏢 %s\n", i+1, booking.RoomName)
		message += fmt.Sprintf("   📅 %s\n", booking.Date.Format("02 Jan 2006"))
		message += fmt.Sprintf("   ⏰ %s - %s\n", booking.StartTime.Format("15:04"), booking.EndTime.Format("15:04"))
		message += fmt.Sprintf("   📝 %s\n", booking.Topic)
		if len(booking.Participants) > 0 {
			message += fmt.Sprintf("   👥 %s\n", strings.Join(booking.Participants, ", "))
		}
		message += fmt.Sprintf("   🔖 ID: `%d`\n\n", booking.BookingID)
	}

	return message
}
