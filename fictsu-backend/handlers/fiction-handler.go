package handlers

import (
	"time"
	"strconv"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

func GetAllFictions(ctx *gin.Context) {
	rows, err := db.DB.Query(
		"SELECT ID, Contributor_ID, Contributor_Name, Cover, Title, Subtitle, Author, Artist, Status, Synopsis, Created FROM Fictions",
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
			&fiction.Contributor_ID,
			&fiction.Contributor_Name,
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
	fiction_id := ctx.Param("fiction_id")
	fiction := models.FictionModel{}
	err := db.DB.QueryRow(
		"SELECT ID, Contributor_ID, Contributor_Name, Cover, Title, Subtitle, Author, Artist, Status, Synopsis, Created FROM Fictions WHERE ID = $1",
		fiction_id,
	).Scan(
		&fiction.ID,
		&fiction.Contributor_ID,
		&fiction.Contributor_Name,
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
		if err == sql.ErrNoRows {
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

func CreateFiction(ctx *gin.Context, store *sessions.CookieStore) {
	session, err_sess := store.Get(ctx.Request, "fictsu-session")
	if err_sess != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to get session"})
		return
	}

	ID_from_session := session.Values["ID"]
	name_from_session := session.Values["name"]
	if ID_from_session == nil || name_from_session == nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized. Please log in to create a fiction."})
		return
	}

	ID_to_DB := ID_from_session.(int)
	name_to_DB := name_from_session.(string)
	fiction_create_request := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fiction_create_request); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for fiction creation"})
		return
	}

	fiction_create_request.Contributor_ID = ID_to_DB
	fiction_create_request.Contributor_Name = name_to_DB

	var new_fiction_ID int
	var new_created_TS time.Time
	err := db.DB.QueryRow(
		"INSERT INTO Fictions (Contributor_ID, Contributor_Name, Cover, Title, Subtitle, Author, Artist, Status, Synopsis) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ID, Created",
		fiction_create_request.Contributor_ID,
		fiction_create_request.Contributor_Name,
		fiction_create_request.Cover,
		fiction_create_request.Title,
		fiction_create_request.Subtitle,
		fiction_create_request.Author,
		fiction_create_request.Artist,
		fiction_create_request.Status,
		fiction_create_request.Synopsis,
	).Scan(&new_fiction_ID, &new_created_TS)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create fiction"})
		return
	}

	fiction_create_request.ID = new_fiction_ID
	fiction_create_request.Created = new_created_TS
	ctx.IndentedJSON(http.StatusCreated, fiction_create_request)
}

func EditFiction(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	fiction_update_request := models.FictionModel{}
	if err := ctx.ShouldBindJSON(&fiction_update_request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input data"})
		return
	}

	query := "UPDATE Fictions SET "
	params := []interface{}{}
	param_index := 1
	if fiction_update_request.Cover != "" {
		query += "Cover = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Cover)
		param_index++
	}

	if fiction_update_request.Title != "" {
		query += "Title = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Title)
		param_index++
	}

	if fiction_update_request.Subtitle != "" {
		query += "Subtitle = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Subtitle)
		param_index++
	}

	if fiction_update_request.Author != "" {
		query += "Author = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Author)
		param_index++
	}

	if fiction_update_request.Artist != "" {
		query += "Artist = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Artist)
		param_index++
	}

	if fiction_update_request.Status != "" {
		query += "Status = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Status)
		param_index++
	}

	if fiction_update_request.Synopsis != "" {
		query += "Synopsis = $" + strconv.Itoa(param_index) + ", "
		params = append(params, fiction_update_request.Synopsis)
		param_index++
	}

	if len(params) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "No valid fields provided for update"})
		return
	}

	query = query[:len(query) - 2] + " WHERE ID = $" + strconv.Itoa(param_index)
	params = append(params, fiction_id)

	result, err := db.DB.Exec(query, params...)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update fiction"})
		return
	}

	rows_affected, _ := result.RowsAffected()
	if rows_affected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Fiction not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Fiction updated successfully"})
}

func DeleteFiction(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	result, err := db.DB.Exec("DELETE FROM Fictions WHERE ID = $1", fiction_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete fiction"})
		return
	}

	rows_affected, _ := result.RowsAffected()
	if rows_affected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Fiction not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Fiction deleted successfully"})
}
