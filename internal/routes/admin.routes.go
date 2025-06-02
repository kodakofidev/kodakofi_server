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

	repoOrder := repositories.NewOrder(db)
	repoUser := repositories.NewUser(db)
	repoProduct := repositories.NewProduct(db)

	handlersOrder := handlers.NewOrder(repoOrder)
	handlersUser := handlers.NewUser(repoUser)
	handlersProduct := handlers.NewProduct(repoProduct)

	route.GET("/data-sales", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.FetchDataSalesHandler)
	route.GET("/users", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersUser.FetchAllUsersHandler)
	route.GET("/users/:id", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersUser.FetchOneUserByAdminHandler)
	route.PATCH("/users/:id", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersUser.PatchUserByAdminHandler)
	// route.POST("/user", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersUser.)
	route.GET("/products", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersProduct.FetchAllProductsAdminHandler)
	route.GET("/orders", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.FetchHistoryOrdersAdminHandler)
	route.GET("/orders/:transaction_code", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.FetchDetailOrderAdminHandler)
	route.GET("/orders/status", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.FetchOrderStatusesHandler)
	route.PATCH("/orders/status", mdw.VerifyToken, mdw.AccsessGate("admin"), handlersOrder.UpdateOrderStatus)
}
