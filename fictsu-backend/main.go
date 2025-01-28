package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/Fictsu/Fictsu/database"
	handlers "github.com/Fictsu/Fictsu/handlers"
)

func main() {
	db.Connection()
	router := gin.Default()

	router.GET("/", handlers.GetAllFictions)
	router.GET("/:fiction_id", handlers.GetFiction)
	router.GET("/:fiction_id/:chapter_id", handlers.GetChapter)

	router.Run()
}
