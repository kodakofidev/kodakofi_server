package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type OrderHandlers struct {
	repo repositories.OrderRepoInterface
}

func NewOrder(repo repositories.OrderRepoInterface) *OrderHandlers {
	return &OrderHandlers{repo: repo}
}

func (h *OrderHandlers) PostOrderHandler(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	order := models.CreateOrderRequest{}

	if err := ctx.ShouldBindJSON(&order); err != nil {
		responder.BadRequest("Invalid request payload", err.Error())
		return
	}

	createOrder, err := h.repo.CreateOrder(ctx, &order)
	if err != nil {
		responder.InternalServerError("Failed to create order", err.Error())
		return
	}

	responder.Created("Order created successfully", createOrder)
}
