// db/room.go

package db

import (
	"database/sql"
	"telegrarmchatbot/internal/model"
)

func GetAllActiveRooms(db *sql.DB) ([]model.Room, error) {
	query := `SELECT room_id, room_name, capacity, COALESCE(status, 'ACTIVE'), create_at 
	          FROM rooms WHERE COALESCE(status, 'ACTIVE') = 'ACTIVE' ORDER BY room_name`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err := rows.Scan(&room.RoomID, &room.RoomName, &room.Capacity, &room.Status, &room.CreateAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	
	return rooms, nil
}

func GetRoomByID(db *sql.DB, roomID int) (*model.Room, error) {
	var room model.Room
	query := `SELECT room_id, room_name, capacity, COALESCE(status, 'ACTIVE'), create_at 
	          FROM rooms WHERE room_id = $1`
	
	err := db.QueryRow(query, roomID).Scan(
		&room.RoomID, &room.RoomName, &room.Capacity, &room.Status, &room.CreateAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &room, nil
}

func GetRoomByName(db *sql.DB, roomName string) (*model.Room, error) {
	var room model.Room
	query := `SELECT room_id, room_name, capacity, COALESCE(status, 'ACTIVE'), create_at 
	          FROM rooms WHERE room_name = $1`
	
	err := db.QueryRow(query, roomName).Scan(
		&room.RoomID, &room.RoomName, &room.Capacity, &room.Status, &room.CreateAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &room, nil
}