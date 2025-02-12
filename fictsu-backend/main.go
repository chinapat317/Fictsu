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

	// GET
	api.GET("", handlers.GetAllFictions)
	api.GET("/user", func(ctx *gin.Context) {
		handlers.GetUserProfile(ctx, store)
	})
	api.GET("/allusers", handlers.GetAllUsers)
	api.GET("f/:fiction_id", handlers.GetFiction)
	api.GET("f/:fiction_id/:chapter_id", handlers.GetChapter)
	api.GET("/auth/:provider", handlers.GetOpenAuthorization)
	api.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		handlers.AuthorizedCallback(ctx, store)
	})

	// POST
	api.POST("/c", func(ctx *gin.Context) {
		handlers.CreateFiction(ctx, store)
	})
	api.POST("f/:fiction_id/c", func(ctx *gin.Context) {
		handlers.CreateChapter(ctx, store)
	})
	api.POST("f/:fiction_id/fav", func(ctx *gin.Context) {
		handlers.AddFavoriteFiction(ctx, store)
	})

	// PUT
	api.PUT("f/:fiction_id/u", func(ctx *gin.Context) {
		handlers.EditFiction(ctx, store)
	})
	api.PUT("f/:fiction_id/:chapter_id/u", func(ctx *gin.Context) {
		handlers.EditChapter(ctx, store)
	})

	// DELETE
	api.DELETE("f/:fiction_id/d", func(ctx *gin.Context) {
		handlers.DeleteFiction(ctx, store)
	})
	api.DELETE("f/:fiction_id/fav/rmv", func(ctx *gin.Context) {
		handlers.RemoveFavoriteFiction(ctx, store)
	})
	api.DELETE("f/:fiction_id/:chapter_id/d", func(ctx *gin.Context) {
		handlers.DeleteChapter(ctx, store)
	})

	router.Run()
}
