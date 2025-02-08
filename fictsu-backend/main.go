package main

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/providers/google"

	db "github.com/Fictsu/Fictsu/database"
	handlers "github.com/Fictsu/Fictsu/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	client_callback_URL := os.Getenv("CLIENT_CALLBACK_URL")

	goth.UseProviders(
		google.New(
			client_id,
			client_secret,
			client_callback_URL,
			"email", "profile",
		),
	)
	
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		HttpOnly: true, 	// Prevent JavaScript access to the cookie
    	Secure:   true, 	// Set to true in production (requires HTTPS)
		MaxAge:   1209600, 	// 2 weeks in seconds
	}

	db.Connection()
	defer db.CloseConnection()

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/", handlers.GetAllFictions)
	api.GET("/:fiction_id", handlers.GetFiction)
	api.GET("/:fiction_id/:chapter_id", handlers.GetChapter)
	api.GET("/auth/:provider", handlers.GetOpenAuthorization)
	api.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		handlers.Callback(ctx, store)
	})

	api.POST("/", handlers.CreateFiction)
	api.POST("/:fiction_id", handlers.CreateChapter)

	api.PUT("/:fiction_id", handlers.EditFiction)
	api.PUT("/:fiction_id/:chapter_id", handlers.EditChapter)

	api.DELETE("/:fiction_id", handlers.DeleteFiction)
	api.DELETE("/:fiction_id/:chapter_id", handlers.DeleteChapter)

	router.Run()
}
