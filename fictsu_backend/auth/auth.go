package auth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/mail"

	env "fictsu_backend/config"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

type TokenRequest struct {
	IDToken string `json:"id_token"`
}

func VerifyTokenHandler(ctx *gin.Context) string {
	body, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var req TokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(err.Error())
		return "Invalid request"
	}
	payload, err := idtoken.Validate(context.Background(), req.IDToken, env.GoogleClientID)
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
