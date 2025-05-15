package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type ProductHandlers struct {
	repo repositories.ProductRepoInterface
}

func NewProduct(repo repositories.ProductRepoInterface) *ProductHandlers {
	return &ProductHandlers{repo: repo}
}

func (h *ProductHandlers) FetchAllProductsHandler(ctx *gin.Context) {

}

func (h *ProductHandlers) FetchDetailProductHandler(ctx *gin.Context) {

}
