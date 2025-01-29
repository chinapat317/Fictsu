package handlers

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)


func GetAllFictions(ctx *gin.Context) {
	rows, err := db.DB.Query(
		"SELECT ID, Cover, Title, Subtitle, Author, Artist, Status, Synopsis, Created FROM Fictions",
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch fictions"})
		return
	}

	defer rows.Close()
	fictions := []models.FictionModel{}
	for rows.Next() {
		fiction := models.FictionModel{}
		if err := rows.Scan(
			&fiction.ID,
			&fiction.Cover,
			&fiction.Title,
			&fiction.Subtitle,
			&fiction.Author,
			&fiction.Artist,
			&fiction.Status,
			&fiction.Synopsis,
			&fiction.Created,
		); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Error processing fictions"})
			return
		}

		fictions = append(fictions, fiction)
	}

	ctx.IndentedJSON(http.StatusOK, fictions)
}

func GetFiction(ctx *gin.Context) {
	fiction_id, err_str_to_int := strconv.Atoi(ctx.Param("fiction_id"))
	if err_str_to_int != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	fiction := models.FictionModel{}
	err := db.DB.QueryRow(
		"SELECT ID, Cover, Title, Subtitle, Author, Artist, Status, Synopsis, Created FROM Fictions WHERE ID = $1", fiction_id,
	).Scan(
		&fiction.ID,
		&fiction.Cover,
		&fiction.Title,
		&fiction.Subtitle,
		&fiction.Author,
		&fiction.Artist,
		&fiction.Status,
		&fiction.Synopsis,
		&fiction.Created,
	)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Fiction not found"})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve fiction"})
		}

		return
	}

	// Get genres of the fiction
	genres, err := GetAllGenres(fiction_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Get chapters of the fiction
	chapters, err := GetAllChapters(fiction_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	fiction.Genres = genres
	fiction.Chapters = chapters
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"Fiction": fiction,
	})
}

func GetAllGenres(ficion_id int) ([]models.GenreModel, error) {
	rows, err := db.DB.Query(
		"SELECT ID, Genre_Name FROM Genres JOIN AssignGenretoFiction ON ID = Genre_ID WHERE Fiction_ID = $1 ORDER BY ID", ficion_id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genres")
	}

	defer rows.Close()
	genres := []models.GenreModel{}
	for rows.Next() {
		genre := models.GenreModel{}
		if err := rows.Scan(
			&genre.ID,
			&genre.Genre_Name,
		); err != nil {
			return nil, fmt.Errorf("failed to process genre data")
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func CreateFiction(ctx *gin.Context) {
	fictionCreateRequest := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fictionCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for fiction creation"})
		return
	}

	fictionCreateRequest.Created = time.Now()

	ctx.IndentedJSON(http.StatusCreated, fictionCreateRequest)
}
