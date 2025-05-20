package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func admin(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {
	route := r.Group("/admin")
	repo := repositories.NewOrder(db)
	handlers := handlers.NewOrder(repo)
	route.GET("sales", mdw.VerifyToken, mdw.AccsessGate("admin"), handlers.FetchTotalSales)
	route.GET("income", mdw.VerifyToken, mdw.AccsessGate("admin"), handlers.FetchIncome)
}
