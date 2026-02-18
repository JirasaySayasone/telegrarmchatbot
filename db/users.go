// db/user.go

package db

import (
	"database/sql"
	"telegrarmchatbot/internal/model"
)

func CreateOrGetUser(db *sql.DB, telegramID int64, username, fullName string) (*model.User, error) {
	var user model.User
	
	// Try to get existing user
	query := `SELECT user_id, telegram_id, username, fullname, create_at 
	          FROM users WHERE telegram_id = $1`
	
	err := db.QueryRow(query, telegramID).Scan(
		&user.UserID, &user.TelegramID, &user.Username, &user.FullName, &user.CreateAt,
	)
	
	if err == sql.ErrNoRows {
		// User doesn't exist, create new one
		insertQuery := `
		INSERT INTO users (telegram_id, username, fullname) 
		VALUES ($1, $2, $3) 
		RETURNING user_id, telegram_id, username, fullname, create_at`
		
		err = db.QueryRow(insertQuery, telegramID, username, fullName).Scan(
			&user.UserID, &user.TelegramID, &user.Username, &user.FullName, &user.CreateAt,
		)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func GetUserByTelegramID(db *sql.DB, telegramID int64) (*model.User, error) {
	var user model.User
	query := `SELECT user_id, telegram_id, username, fullname, create_at 
	          FROM users WHERE telegram_id = $1`
	
	err := db.QueryRow(query, telegramID).Scan(
		&user.UserID, &user.TelegramID, &user.Username, &user.FullName, &user.CreateAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}