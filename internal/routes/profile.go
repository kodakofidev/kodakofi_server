package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func profile(r *gin.Engine, db *pgxpool.Pool) {
	route := r.Group("/api/profile")
	repo := repositories.NewProfile(db)
	handlers := handlers.NewProfileHandlers(repo)

	route.GET("", handlers.FetchProfileHandler)
	route.POST("/edit")
}
