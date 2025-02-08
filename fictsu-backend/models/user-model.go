package models

import (
	"time"
)

type UserModel struct {
	ID        	int			`json:"id"`
	User_ID    	string		`json:"user_id"`
	Super_User 	int			`json:"super_user"`
	Name      	string		`json:"name"`
	Email     	string		`json:"email"`
	Avatar_URL 	string		`json:"avatar_url"`
	Joined 		time.Time	`json:"joined"`
}
