package handlers

import (
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

func GetAllChapters(fiction_id string) ([]models.ChapterModel, error) {
	rows, err := db.DB.Query(
		"SELECT ID, Fiction_ID, Title, Content, Created FROM Chapters WHERE Fiction_ID = $1 ORDER BY ID", fiction_id,
	)
	if err != nil {
		return []models.ChapterModel{}, fmt.Errorf("no such a chapter")
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
			return []models.ChapterModel{}, fmt.Errorf("no such a chapter")
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func GetChapter(ctx *gin.Context) {
	fiction_id := ctx.Param("fiction_id")
	chapter_id := ctx.Param("chapter_id")
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
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, chapter)
}

func CreateChapter(ctx *gin.Context) {
	chapterCreateRequest := models.ChapterModel{}
	if err := ctx.ShouldBindJSON(&chapterCreateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	chapterCreateRequest.Created = time.Now()

	ctx.IndentedJSON(http.StatusCreated, chapterCreateRequest)
}
