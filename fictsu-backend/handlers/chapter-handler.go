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

func GetAllChapters(fiction_id int) ([]models.ChapterModel, error) {
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
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	chapter_id, err := strconv.Atoi(ctx.Param("chapter_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid chapter ID"})
		return
	}

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
		if err_select == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve chapter"})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, chapter)
}

func CreateChapter(ctx *gin.Context) {
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	chapterCreateRequest := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapterCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for chapter creation"})
		return
	}

	var nextChapterID int
	err_nextChapterID := db.DB.QueryRow(
		"SELECT COALESCE(MAX(ID), 0) + 1 FROM Chapters WHERE Fiction_ID = $1",
		fiction_id,
	).Scan(&nextChapterID)
	if err_nextChapterID != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to calculate next chapter ID"})
        return
	}

	var newCreatedTS time.Time
	err_insert := db.DB.QueryRow(
		"INSERT INTO Chapters (Fiction_ID, ID, Title, Content) VALUES ($1, $2, $3, $4) RETURNING Created",
		fiction_id,
		nextChapterID,
		chapterCreateRequest.Title,
		chapterCreateRequest.Content,
	).Scan(&newCreatedTS)
	if err_insert != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create chapter"})
		return
	}

	chapterCreateRequest.ID = nextChapterID
	chapterCreateRequest.Created = newCreatedTS
	ctx.IndentedJSON(http.StatusCreated, chapterCreateRequest)
}

func EditChapter(ctx *gin.Context) {
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	chapter_id, err := strconv.Atoi(ctx.Param("chapter_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid chapter ID"})
		return
	}

	chapterUpdateRequest := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapterUpdateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input data"})
		return
	}

	query := "UPDATE Chapters SET "
	params := []interface{}{}
	paramIndex := 1
	if chapterUpdateRequest.Title != "" {
		query += "Title = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, chapterUpdateRequest.Title)
		paramIndex++
	}

	if chapterUpdateRequest.Content != "" {
		query += "Content = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, chapterUpdateRequest.Content)
		paramIndex++
	}

	if len(params) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "No valid fields provided for update"})
		return
	}

	query = strings.TrimSuffix(query, ", ") + " WHERE ID = $" + strconv.Itoa(paramIndex) + " AND Fiction_ID = $" + strconv.Itoa(paramIndex + 1)
	params = append(params, chapter_id, fiction_id)

	result, err := db.DB.Exec(query, params...)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update chapter"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Chapter updated successfully"})
}

func DeleteChapter(ctx *gin.Context) {
	fiction_id, err := strconv.Atoi(ctx.Param("fiction_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	chapter_id, err := strconv.Atoi(ctx.Param("chapter_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid chapter ID"})
		return
	}

	result, err := db.DB.Exec("DELETE FROM Chapters WHERE Fiction_ID = $1 AND ID = $2", fiction_id, chapter_id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete chapter"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Chapter deleted successfully"})
}
