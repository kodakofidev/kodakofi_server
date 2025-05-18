package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	middlewares := middlewares.InitMiddleware()
	rg := router.Group("/api")
	auth(rg, db)
	profile(rg, db, middlewares)
	product(rg, db)
	order(rg, db, middlewares)
	return router
}
