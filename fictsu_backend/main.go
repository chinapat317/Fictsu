package main

import (
	"fictsu_backend/handlers"
	"time"

	env "fictsu_backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func CreateRoute() {
	router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func main() {
	env.LoadEnv()
	env.ConnectDatabase()
	env.InitFirebaseApp()
	CreateRoute()
	api := router.Group("/api")
	api.POST("/f/c", handlers.CreateFiction)
	router.Run()
}
