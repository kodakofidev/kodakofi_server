package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func profile(r *gin.RouterGroup, db *pgxpool.Pool) {

	repo := repositories.NewProfile(db)
	handlers := handlers.NewProfileHandlers(repo)

	profile := r.Group("/profile")

	profile.GET("", handlers.FetchProfileHandler)
	profile.POST("/edit")
}
