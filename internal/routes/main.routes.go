package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	middlewares := middlewares.InitMiddleware()
	router.Use(middlewares.CORSMiddleware)
	rg := router.Group("/api")
	auth(rg, db)
	profile(rg, db, middlewares)
	product(rg, db, middlewares)
	order(rg, db, middlewares)
	admin(rg, db, middlewares)
	return router
}
