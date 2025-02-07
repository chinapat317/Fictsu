package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func Authen(ctx *gin.Context) {
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

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	session, err := store.Get(ctx.Request, "fictsu-session")
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create session"})
		return
	}

	session.Values["userID"] = user.UserID
	session.Values["name"] = user.Name
	session.Values["email"] = user.Email
	session.Values["avatar"] = user.AvatarURL
	err = session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save session"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"ID": user.UserID,
		"Name": user.Name,
		"Email": user.Email,
		"Avatar": user.AvatarURL,
	})
}
