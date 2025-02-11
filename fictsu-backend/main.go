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
    store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}

	goth.UseProviders(
		google.New(
			client_id,
			client_secret,
			client_callback_URL,
			"email", "profile",
		),
	)

	db.Connection()
	defer db.CloseConnection()

	router := gin.Default()

	api := router.Group("/api")

	api.GET("", handlers.GetAllFictions)
	api.GET("/user", func(ctx *gin.Context) {
		handlers.GetUserProfile(ctx, store)
	})
	api.GET("/allusers", handlers.GetAllUsers)
	api.GET("/:fiction_id", handlers.GetFiction)
	api.GET("/:fiction_id/:chapter_id", handlers.GetChapter)
	api.GET("/auth/:provider", handlers.GetOpenAuthorization)
	api.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		handlers.AuthorizedCallback(ctx, store)
	})

	api.POST("", func(ctx *gin.Context) {
		handlers.CreateFiction(ctx, store)
	})
	api.POST("/:fiction_id", handlers.CreateChapter)

	api.PUT("/:fiction_id", handlers.EditFiction)
	api.PUT("/:fiction_id/:chapter_id", handlers.EditChapter)

	api.DELETE("/:fiction_id", handlers.DeleteFiction)
	api.DELETE("/:fiction_id/:chapter_id", handlers.DeleteChapter)

	router.Run()
}
