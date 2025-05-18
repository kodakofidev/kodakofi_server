package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func product(r *gin.RouterGroup, db *pgxpool.Pool) {
	route := r.Group("/product")
	repo := repositories.NewProduct(db)
	handlers := handlers.NewProduct(repo)
	route.GET("", handlers.FetchAllProductsHandler)
	route.GET("/:id", handlers.FetchDetailProductHandler)
}
