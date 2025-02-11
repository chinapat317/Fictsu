package handlers

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

func GetAllChapters(fiction_id string) ([]models.ChapterModel, error) {
	rows, err := db.DB.Query(
		"SELECT Fiction_ID, ID, Title, Content, Created FROM Chapters WHERE Fiction_ID = $1 ORDER BY ID",
		fiction_id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chapters")
	}

	defer rows.Close()
	chapters := []models.ChapterModel{}
	for rows.Next() {
		chapter := models.ChapterModel{}
		if err := rows.Scan(
			&chapter.Fiction_ID,
			&chapter.ID,
			&chapter.Title,
			&chapter.Content,
			&chapter.Created,
		); err != nil {
			return nil, fmt.Errorf("failed to process chapter data")
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func GetChapter(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	chapter_id := ctx.Param("chapter_id")

	chapter := models.ChapterModel{}
	err_select := db.DB.QueryRow(
		"SELECT Fiction_ID, ID, Title, Content, Created FROM Chapters WHERE Fiction_ID = $1 AND ID = $2",
		fiction_id,
		chapter_id,
	).Scan(
		&chapter.Fiction_ID,
		&chapter.ID,
		&chapter.Title,
		&chapter.Content,
		&chapter.Created,
	)
	if err_select != nil {
		if err_select == sql.ErrNoRows {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve chapter"})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, chapter)
}

func CreateChapter(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")

	chapter_create_request := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapter_create_request); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for chapter creation"})
		return
	}

	var next_chapter_ID int
	err_next_chapter_ID := db.DB.QueryRow(
		"SELECT COALESCE(MAX(ID), 0) + 1 FROM Chapters WHERE Fiction_ID = $1",
		fiction_id,
	).Scan(&next_chapter_ID)
	if err_next_chapter_ID != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to calculate next chapter ID"})
        return
	}

	var new_created_TS time.Time
	err_insert := db.DB.QueryRow(
		"INSERT INTO Chapters (Fiction_ID, ID, Title, Content) VALUES ($1, $2, $3, $4) RETURNING Created",
		fiction_id,
		next_chapter_ID,
		chapter_create_request.Title,
		chapter_create_request.Content,
	).Scan(&new_created_TS)
	if err_insert != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create chapter"})
		return
	}

	chapter_create_request.ID = next_chapter_ID
	chapter_create_request.Created = new_created_TS
	ctx.IndentedJSON(http.StatusCreated, chapter_create_request)
}

func EditChapter(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	chapter_id := ctx.Param("chapter_id")

	chapter_update_request := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapter_update_request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input data"})
		return
	}

	query := "UPDATE Chapters SET "
	params := []interface{}{}
	param_index := 1
	if chapter_update_request.Title != "" {
		query += "Title = $" + strconv.Itoa(param_index) + ", "
		params = append(params, chapter_update_request.Title)
		param_index++
	}

	if chapter_update_request.Content != "" {
		query += "Content = $" + strconv.Itoa(param_index) + ", "
		params = append(params, chapter_update_request.Content)
		param_index++
	}

	if len(params) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "No valid fields provided for update"})
		return
	}

	query = strings.TrimSuffix(query, ", ") + " WHERE ID = $" + strconv.Itoa(param_index) + " AND Fiction_ID = $" + strconv.Itoa(param_index + 1)
	params = append(params, chapter_id, fiction_id)

	result, err := db.DB.Exec(query, params...)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update chapter"})
		return
	}

	rows_affected, _ := result.RowsAffected()
	if rows_affected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Chapter updated successfully"})
}

func DeleteChapter(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	chapter_id := ctx.Param("chapter_id")

	result, err := db.DB.Exec("DELETE FROM Chapters WHERE Fiction_ID = $1 AND ID = $2", fiction_id, chapter_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete chapter"})
		return
	}

	rows_affected, _ := result.RowsAffected()
	if rows_affected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Chapter deleted successfully"})
}
