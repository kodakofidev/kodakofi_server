package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func orders(r *gin.RouterGroup, db *pgxpool.Pool) {
	repo := repositories.NewOrder(db)
	handlers := handlers.NewOrder(repo)

	orders := r.Group("/orders")
	orders.GET("", handlers.GetHistoryOrders)
}