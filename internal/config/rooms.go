// internal/config/rooms.go - NEW FILE

package config

var (
	RoomNames     = []string{"Room A", "Room B", "Room C"}
	WorkdayStart  = "09:00"
	WorkdayEnd    = "17:00"
	SlotDuration  = 60 // minutes
)

// GenerateTimeSlots creates hourly slots from 09:00 to 17:00
func GenerateTimeSlots() []string {
	slots := []string{
		"09:00", "10:00", "11:00", "12:00", 
		"13:00", "14:00", "15:00", "16:00",
	}
	return slots
}