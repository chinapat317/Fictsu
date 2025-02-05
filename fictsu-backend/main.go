package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	handlers "github.com/Fictsu/Fictsu/handlers"
)

func main() {
	db.Connection()
	defer db.CloseConnection()

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/", handlers.GetAllFictions)
	api.GET("/:fiction_id", handlers.GetFiction)
	api.GET("/:fiction_id/:chapter_id", handlers.GetChapter)

	api.POST("/", handlers.CreateFiction)
	api.POST("/:fiction_id", handlers.CreateChapter)

	api.PUT("/:fiction_id", handlers.EditFiction)
	api.PUT("/:fiction_id/:chapter_id", handlers.EditChapter)

	api.DELETE("/:fiction_id", handlers.DeleteFiction)
	api.DELETE("/:fiction_id/:chapter_id", handlers.DeleteChapter)

	router.Run()
}
