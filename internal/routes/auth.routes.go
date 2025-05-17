package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func auth(r *gin.RouterGroup, db *pgxpool.Pool) {
	repo := repositories.NewAuth(db)
	handlers := handlers.NewAuthHandlers(repo)

	auth := r.Group("/auth")
	auth.POST("", handlers.Login)
	auth.POST("/new", handlers.Register)
	auth.POST("/verify", handlers.VerifyEmail)
	auth.POST("/otp", handlers.SendOTP)
}
