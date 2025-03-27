package handlers

import (
	"fictsu_backend/auth"
	"fmt"
	"net/http"

	env "fictsu_backend/config"

	"github.com/gin-gonic/gin"
)

func CreateFiction(ctx *gin.Context) {
	var email string = auth.VerifyTokenHandler(ctx)
	if !auth.IsValidEmail(email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": email})
		return
	}
	ctx.Request.ParseMultipartForm(10 << 20)
	file, header, err := ctx.Request.FormFile("cover")
	if err == nil {
		url, err := UploadImageToFirebase(file, header, env.CoverPath, env.BucketName)
		if err != nil {
			fmt.Println(err)
		}
		ctx.JSON(http.StatusOK, gin.H{"status": url})
		return
	} else {
		fmt.Println("Err is: ", err)
		ctx.JSON(http.StatusOK, gin.H{"status": err.Error()})
		return
	}
}
