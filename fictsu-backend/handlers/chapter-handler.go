package handlers

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"

	models "github.com/Fictsu/Fictsu/models"
)

func CreateChapter(ctx *gin.Context) {
	chapterCreateRequest := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapterCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	chapterCreateRequest.Created = time.Now()

	ctx.IndentedJSON(http.StatusCreated, chapterCreateRequest)
}
