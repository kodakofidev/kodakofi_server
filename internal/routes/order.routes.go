package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func order(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {
	route := r.Group("/order")
	repo := repositories.NewOrder(db)
	handlers := handlers.NewOrder(repo)
	route.POST("", mdw.VerifyToken, mdw.AccsessGate("user"), handlers.PostOrderHandler)
	route.GET("", mdw.VerifyToken, mdw.AccsessGate("user"), handlers.GetHistoryOrders)
}
