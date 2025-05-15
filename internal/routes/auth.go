package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func auth(r *gin.Engine, db *pgxpool.Pool) {
	route := r.Group("/auth")
	repo := repositories.NewAuth(db)
	handlers := handlers.NewAuthHandlers(repo)

	route.POST("", handlers.Login)
	route.POST("/new", handlers.Register)
}
