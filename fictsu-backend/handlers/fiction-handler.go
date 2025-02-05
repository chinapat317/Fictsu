package handlers

import (
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

func CreateFiction(ctx *gin.Context) {
	fictionCreateRequest := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fictionCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for fiction creation"})
		return
	}

	var newFictionID int
	var newCreatedTS time.Time
	err := db.DB.QueryRow(
		"INSERT INTO Fictions (Cover, Title, Subtitle, Author, Artist, Status, Synopsis) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID, Created",
		fictionCreateRequest.Cover,
		fictionCreateRequest.Title,
		fictionCreateRequest.Subtitle,
		fictionCreateRequest.Author,
		fictionCreateRequest.Artist,
		fictionCreateRequest.Status,
		fictionCreateRequest.Synopsis,
	).Scan(&newFictionID, &newCreatedTS)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create fiction"})
		return
	}

	fictionCreateRequest.ID = newFictionID
	fictionCreateRequest.Created = newCreatedTS
	ctx.IndentedJSON(http.StatusCreated, fictionCreateRequest)
}

func EditFiction(ctx *gin.Context) {
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	fictionUpdateRequest := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fictionUpdateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input data"})
		return
	}

	query := "UPDATE Fictions SET "
	params := []interface{}{}
	paramIndex := 1
	if fictionUpdateRequest.Cover != "" {
		query += "Cover = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Cover)
		paramIndex++
	}

	if fictionUpdateRequest.Title != "" {
		query += "Title = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Title)
		paramIndex++
	}

	if fictionUpdateRequest.Subtitle != "" {
		query += "Subtitle = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Subtitle)
		paramIndex++
	}

	if fictionUpdateRequest.Author != "" {
		query += "Author = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Author)
		paramIndex++
	}

	if fictionUpdateRequest.Artist != "" {
		query += "Artist = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Artist)
		paramIndex++
	}

	if fictionUpdateRequest.Status != "" {
		query += "Status = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Status)
		paramIndex++
	}

	if fictionUpdateRequest.Synopsis != "" {
		query += "Synopsis = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, fictionUpdateRequest.Synopsis)
		paramIndex++
	}

	if len(params) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "No valid fields provided for update"})
		return
	}

	query = query[:len(query) - 2] + " WHERE ID = $" + strconv.Itoa(paramIndex)
	params = append(params, fiction_id)

	result, err := db.DB.Exec(query, params...)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update fiction"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Fiction not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Fiction updated successfully"})
}

func DeleteFiction(ctx *gin.Context) {
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	result, err := db.DB.Exec("DELETE FROM Fictions WHERE ID = $1", fiction_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete fiction"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Fiction not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Fiction deleted successfully"})
}
