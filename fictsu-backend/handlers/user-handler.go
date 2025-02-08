package handlers

import (
	"fmt"
	"time"
	"database/sql"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

func GetUser(user_id string) (*models.UserModel, error) {
	user := models.UserModel{}
	err := db.DB.QueryRow(
		"SELECT ID, User_ID, Name, Email, Avatar_URL, Joined FROM Users WHERE User_ID = $1",
		user_id,
	).Scan(
		&user.ID,
		&user.User_ID,
		&user.Name,
		&user.Email,
		&user.Avatar_URL,
		&user.Joined,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to retrieve user from database: %v", err)
		}
	}

	return &user, nil
}

func CreateUser(user *models.UserModel) (*models.UserModel, error) {
	var newUserID int
	var newUserGoogleID string
	var newUserJoined time.Time
	err := db.DB.QueryRow(
		"INSERT INTO Users (User_ID, Name, Email, Avatar_URL) VALUES ($1, $2, $3, $4) RETURNING ID, User_ID, Joined",
		user.User_ID,
		user.Name,
		user.Email,
		user.Avatar_URL,
	).Scan(
		&newUserID,
		&newUserGoogleID,
		&newUserJoined,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account: %v", err)
	}

	user.ID = newUserID
	user.User_ID = newUserGoogleID
	user.Joined = newUserJoined
	return user, nil
}
