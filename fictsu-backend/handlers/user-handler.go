package handlers

import (
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

func GetUserProfile(ctx *gin.Context, store *sessions.CookieStore) {
	session, err_sess := store.Get(ctx.Request, "fictsu-session")
	if err_sess != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to get session"})
		return
	}

	ID_from_session := session.Values["ID"]
	if ID_from_session == nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}

	ID_to_DB := ID_from_session.(int)
	user := models.UserModel{}
	err := db.DB.QueryRow(
		`
		SELECT
			ID, User_ID, Super_User, Name, Email, Avatar_URL, Joined
		FROM
			Users
		WHERE
			ID = $1
		`,
		ID_to_DB,
	).Scan(
		&user.ID,
		&user.User_ID,
		&user.Super_User,
		&user.Name,
		&user.Email,
		&user.Avatar_URL,
		&user.Joined,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve user details"})
		}

		return
	}

	fav_fictions, err := GetFavFictions(ID_to_DB)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve favorite fictions"})
		return
	}

	contri_fictions, err := GetContributedFictions(ID_to_DB)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve contributed fictions"})
		return
	}

	user.ID = ID_to_DB
	user.Fav_Fictions = fav_fictions
	user.Contributed_Fic = contri_fictions
	ctx.IndentedJSON(http.StatusOK, gin.H{"User_Profile": user})
}

func GetUser(user_id string) (*models.UserModel, error) {
	user := models.UserModel{}
	err := db.DB.QueryRow(
		`
		SELECT
			ID, User_ID, Name, Email, Avatar_URL, Joined
		FROM
			Users
		WHERE
			User_ID = $1
		`,
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
	var new_user_ID int
	var new_user_Google_ID string
	var new_user_joined time.Time
	err := db.DB.QueryRow(
		`
		INSERT INTO Users (User_ID, Name, Email, Avatar_URL)
		VALUES ($1, $2, $3, $4)
		RETURNING ID, User_ID, Joined
		`,
		user.User_ID,
		user.Name,
		user.Email,
		user.Avatar_URL,
	).Scan(
		&new_user_ID,
		&new_user_Google_ID,
		&new_user_joined,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account: %v", err)
	}

	user.ID = new_user_ID
	user.User_ID = new_user_Google_ID
	user.Joined = new_user_joined
	return user, nil
}
