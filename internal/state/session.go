// internal/state/session.go - NEW FILE

package state

import (
	"sync"
	"time"
)

type BookingSession struct {
	UserID       int64
	Step         string // "select_room", "select_time", "enter_topic", "enter_participants"
	RoomID       int
	RoomName     string
	Date         time.Time
	StartTime    string
	EndTime      string
	Topic        string
	Participants []string
}

type SessionManager struct {
	sessions map[int64]*BookingSession
	mu       sync.RWMutex
}

var Manager = &SessionManager{
	sessions: make(map[int64]*BookingSession),
}

func (sm *SessionManager) GetSession(userID int64) *BookingSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[userID]
}

func (sm *SessionManager) SetSession(userID int64, session *BookingSession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userID] = session
}

func (sm *SessionManager) ClearSession(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, userID)
}

func (sm *SessionManager) StartBooking(userID int64) {
	sm.SetSession(userID, &BookingSession{
		UserID: userID,
		Step:   "select_room",
		Date:   time.Now(),
	})
}