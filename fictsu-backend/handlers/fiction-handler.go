package handlers

import (
	"time"
	"net/http"
	// "database/sql"
	"github.com/gin-gonic/gin"

	models "github.com/Fictsu/Fictsu/models"
	// db "github.com/Fictsu/Fictsu/database"
)


func GetAllFictions(ctx *gin.Context) {
// 	rows, err := db.DB.Query(

// 	)
}

func GetFiction(ctx *gin.Context) {
	
}

func CreateFiction(ctx *gin.Context) {
	fictionCreateRequest := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fictionCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	fictionCreateRequest.Created = time.Now()

	ctx.IndentedJSON(http.StatusCreated, fictionCreateRequest)
}
