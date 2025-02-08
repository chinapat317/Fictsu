package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	models "github.com/Fictsu/Fictsu/models"
)

func GetOpenAuthorization(ctx *gin.Context) {
	provider := ctx.Param("provider")
	query := ctx.Request.URL.Query()
	query.Add("provider", provider)
	ctx.Request.URL.RawQuery = query.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func Callback(ctx *gin.Context, store *sessions.CookieStore) {
	provider := ctx.Param("provider")
	query := ctx.Request.URL.Query()
	query.Add("provider", provider)
	ctx.Request.URL.RawQuery = query.Encode()

	user, err_complete := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err_complete != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error1": err_complete.Error()})
		return
	}

	// Check if the user exists in the database
	user_in_db, err_get := GetUser(user.UserID)
	if err_get != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error2": err_get.Error()})
		return
	}

	// Create a new user if not found
	if user_in_db == nil {
		new_user := &models.UserModel{
			User_ID: user.UserID,
			Name: user.Name,
			Email: user.Email,
			Avatar_URL: user.AvatarURL,
		}

		created_user, err_create_user := CreateUser(new_user)
		if err_create_user != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error3": err_create_user.Error()})
			return
		}

		user_in_db = created_user
	}

	session, err_sess := store.Get(ctx.Request, "fictsu-session")
	if err_sess != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error4": "Failed to create session"})
		return
	}

	session.Values["userID"] = user_in_db.User_ID
	session.Values["name"] = user_in_db.Name
	session.Values["email"] = user_in_db.Email
	session.Values["avatar"] = user_in_db.Avatar_URL
	err_sess = session.Save(ctx.Request, ctx.Writer)
	if err_sess != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error5": "Failed to save session"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"ID": user_in_db.ID,
		"User_ID": user_in_db.User_ID,
		"Name": user_in_db.Name,
		"Email": user_in_db.Email,
		"Avatar_URL": user_in_db.Avatar_URL,
		"Joined": user_in_db.Joined,
	})
}
