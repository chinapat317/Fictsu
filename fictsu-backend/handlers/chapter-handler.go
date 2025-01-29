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

func GetAllChapters(fiction_id int) ([]models.ChapterModel, error) {
	rows, err := db.DB.Query(
		"SELECT ID, Fiction_ID, Title, Content, Created FROM Chapters WHERE Fiction_ID = $1 ORDER BY ID", fiction_id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chapters")
	}

	defer rows.Close()
	chapters := []models.ChapterModel{}
	for rows.Next() {
		chapter := models.ChapterModel{}
		if err := rows.Scan(
			&chapter.ID,
			&chapter.FictionID,
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
	fiction_id, err_str_to_int := strconv.Atoi(ctx.Param("fiction_id"))
	if err_str_to_int != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid fiction ID"})
		return
	}

	chapter_id, err_str_to_int := strconv.Atoi(ctx.Param("chapter_id"))
	if err_str_to_int != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid chapter ID"})
		return
	}

	chapter := models.ChapterModel{}
	err := db.DB.QueryRow(
		"SELECT ID, Fiction_ID, Title, Content, Created FROM Chapters WHERE Fiction_ID = $1 AND ID = $2 ORDER BY ID",
		fiction_id,
		chapter_id,
	).Scan(
		&chapter.ID,
		&chapter.FictionID,
		&chapter.Title,
		&chapter.Content,
		&chapter.Created,
	)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Chapter not found"})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve chapter"})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, chapter)
}

func CreateChapter(ctx *gin.Context) {
	chapterCreateRequest := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapterCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data provided for chapter creation"})
		return
	}

	chapterCreateRequest.Created = time.Now()

	ctx.IndentedJSON(http.StatusCreated, chapterCreateRequest)
}
