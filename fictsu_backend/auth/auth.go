package auth

import (
	"context"
	"fmt"
	"net/mail"

	env "fictsu_backend/config"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

type TokenRequest struct {
	IDToken string `json:"id_token"`
}

func VerifyTokenHandler(ctx *gin.Context) string {
	idToken := ctx.PostForm("id_token")
	payload, err := idtoken.Validate(context.Background(), idToken, env.GoogleClientID)
	if err != nil {
		fmt.Println(err.Error())
		return "Invalid request"
	}
	var email string = payload.Claims["email"].(string)
	return email
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
