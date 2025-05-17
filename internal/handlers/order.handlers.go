package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type OrderHandlers struct {
	repo repositories.OrderRepoInterface
}

func NewOrder(repo repositories.OrderRepoInterface) *OrderHandlers {
	return &OrderHandlers{repo: repo}
}

func (h *OrderHandlers) PostOrderHandler(ctx *gin.Context) {

}
