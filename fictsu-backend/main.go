package main

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/providers/google"

	db 	"github.com/Fictsu/Fictsu/database"
	configs "github.com/Fictsu/Fictsu/configs"
	handlers "github.com/Fictsu/Fictsu/handlers"
)

func main() {
	configs.LoadEnv()

	store := sessions.NewCookieStore([]byte(configs.SessionKey))
	store.Options = &sessions.Options{
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}

	goth.UseProviders(
		google.New(
			configs.ClientID,
			configs.ClientSecret,
			configs.ClientCallbackURL,
			"email", "profile",
		),
	)

	db.Connection()
	defer db.CloseConnection()

	router := gin.Default()

	API := router.Group("/api")

	// GET
	API.GET("", handlers.GetAllFictions)
	API.GET("/user", func(ctx *gin.Context) {
		handlers.GetUserProfile(ctx, store)
	})
	API.GET("/allusers", handlers.GetAllUsers)
	API.GET("f/:fiction_id", handlers.GetFiction)
	API.GET("f/:fiction_id/:chapter_id", handlers.GetChapter)
	API.GET("/auth/:provider", handlers.GetOpenAuthorization)
	API.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		handlers.AuthorizedCallback(ctx, store)
	})

	// POST
	API.POST("/c", func(ctx *gin.Context) {
		handlers.CreateFiction(ctx, store)
	})
	API.POST("f/:fiction_id/c", func(ctx *gin.Context) {
		handlers.CreateChapter(ctx, store)
	})
	API.POST("f/:fiction_id/fav", func(ctx *gin.Context) {
		handlers.AddFavoriteFiction(ctx, store)
	})

	// PUT
	API.PUT("f/:fiction_id/u", func(ctx *gin.Context) {
		handlers.EditFiction(ctx, store)
	})
	API.PUT("f/:fiction_id/:chapter_id/u", func(ctx *gin.Context) {
		handlers.EditChapter(ctx, store)
	})

	// DELETE
	API.DELETE("f/:fiction_id/d", func(ctx *gin.Context) {
		handlers.DeleteFiction(ctx, store)
	})
	API.DELETE("f/:fiction_id/fav/rmv", func(ctx *gin.Context) {
		handlers.RemoveFavoriteFiction(ctx, store)
	})
	API.DELETE("f/:fiction_id/:chapter_id/d", func(ctx *gin.Context) {
		handlers.DeleteChapter(ctx, store)
	})

	// OpenAI
	AI := API.Group("/ai")

	// POST
	AI.POST("/t", handlers.OpenAIGetText)
	AI.POST("/tti", handlers.OpenAIGetTextToImage)

	router.Run()
}
