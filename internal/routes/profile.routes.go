package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func profile(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {

	repo := repositories.NewProfile(db)
	handlers := handlers.NewProfileHandlers(repo)

	profile := r.Group("/profile")

	profile.GET("", mdw.VerifyToken, mdw.AccsessGate("user"), handlers.FetchProfileHandler)
	profile.PATCH("/edit", mdw.VerifyToken, mdw.AccsessGate("user"), handlers.EditProfileHandler)
	// profile.PATCH("/edit", handlers.EditProfileHandler)

}
