package handlers

import (
	"time"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)


func GetAllFictions(ctx *gin.Context) {
	rows, err := db.DB.Query(
		"SELECT ID, Title, Subtitle, Author, Artist, Status, Synopsis, Genre, Created FROM Fictions",
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	defer rows.Close()
	fictions := []models.FictionModel{}
	for rows.Next() {
		fiction := models.FictionModel{}
		if err := rows.Scan(
			&fiction.ID,
			&fiction.Title,
			&fiction.Subtitle,
			&fiction.Author,
			&fiction.Artist,
			&fiction.Status,
			&fiction.Synopsis,
			&fiction.Genre,
			&fiction.Created,
		); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		fictions = append(fictions, fiction)
	}

	ctx.IndentedJSON(http.StatusOK, fictions)
}

func GetFiction(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	fiction := models.FictionModel{}
	err := db.DB.QueryRow(
		"SELECT ID, Title, Subtitle, Author, Artist, Status, Synopsis, Genre, Created FROM Fictions WHERE ID = $1", fiction_id,
	).Scan(
		&fiction.ID,
		&fiction.Title,
		&fiction.Subtitle,
		&fiction.Author,
		&fiction.Artist,
		&fiction.Status,
		&fiction.Synopsis,
		&fiction.Genre,
		&fiction.Created,
	)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}

		return
	}

	// Get chapters of the fiction
	chapters, err := GetAllChapters(fiction_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"Fiction": fiction,
		"Chapters": chapters,
	})
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
