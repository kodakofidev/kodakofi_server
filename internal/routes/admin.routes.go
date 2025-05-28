package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/middlewares"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

func admin(r *gin.RouterGroup, db *pgxpool.Pool, mdw *middlewares.Middleware) {
	route := r.Group("admin")

	repoOrder := repositories.NewOrder(db)
	repoUser := repositories.NewUser(db)
	repoProduct := repositories.NewProduct(db)

	handlersOrder := handlers.NewOrder(repoOrder)
	handlersUser := handlers.NewUser(repoUser)
	handlersProduct := handlers.NewProduct(repoProduct)

	route.GET("data-sales", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.FetchDataSalesHandler)
	route.PATCH("status", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.UpdateOrderStatus)
	route.GET("users", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersUser.FetchAllUsersHandler)
	// route.PATCH("user", mdw.VerifyToken, mdw.AccsessGate("admin"), handlers.)
	// route.POST("user", mdw.VerifyToken, mdw.AccsessGate("admin"), handlers.)
	route.GET("products", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersProduct.FetchAllProductsAdminHandler)
}
