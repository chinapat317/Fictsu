package handlers

import (
	"fictsu_backend/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateFiction(ctx *gin.Context) {
	var is_verify string = auth.VerifyTokenHandler(ctx)
	if !auth.IsValidEmail(is_verify) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": is_verify})
		return
	}
}
