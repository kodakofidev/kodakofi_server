package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
<<<<<<< HEAD
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func orders(r *gin.RouterGroup, db *pgxpool.Pool) {
	repo := repositories.NewOrder(db)
	handlers := handlers.NewOrder(repo)

	orders := r.Group("/orders")
	orders.GET("", handlers.GetHistoryOrders)
}
=======
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func order(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {
	route := r.Group("/order")
	repo := repositories.NewOrder(db)
	handlers := handlers.NewOrder(repo)
	route.POST("", mdw.VerifyToken, mdw.AccsessGate("user"), handlers.PostOrderHandler)
}
>>>>>>> fdce5597d055ddee6b11ceed27a2e3b1af773002
