package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	rg := router.Group("/api")
	auth(rg, db)
	profile(rg, db)
	orders(rg, db)
	return router
}
