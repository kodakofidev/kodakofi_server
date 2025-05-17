package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type ProductHandlers struct {
	repo repositories.ProductRepoInterface
}

func NewProduct(repo repositories.ProductRepoInterface) *ProductHandlers {
	return &ProductHandlers{repo: repo}
}

func (h *ProductHandlers) FetchAllProductsHandler(ctx *gin.Context) {
	var params *models.ProductQueryParams
	response := models.NewResponse(ctx)
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest("params invalid", err.Error())
	}
	res, err := h.repo.GetAllProducts(ctx.Request.Context(), params)
	if err != nil {
		response.InternalServerError("internal server errors", err.Error())
	}
	response.Success("get products success", res)
}

func (h *ProductHandlers) FetchDetailProductHandler(ctx *gin.Context) {

}
