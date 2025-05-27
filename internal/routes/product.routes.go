package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func product(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {
	route := r.Group("/product")
	repo := repositories.NewProduct(db)
	handlers := handlers.NewProduct(repo)
	route.GET("", handlers.FetchAllProductsHandler)
	route.GET("/:id", handlers.FetchDetailProductHandler)
	route.PATCH("/:id", handlers.UpdateProduct)
	route.POST("/:id/product-like", mdw.VerifyToken, handlers.UpdateLikeProduct)
	// route.GET("/:id/product-like", handlers.GetLikeProduct)
	// routeImg := route.Group("/image")
	// routeImg.DELETE("/:id", handlers.RemoveImage)
	// route.POST("", handlers.AddProduct)

	route.POST("", mdw.VerifyToken, mdw.AccsessGate("admin"), handlers.AddProduct)

}
