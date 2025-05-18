package handlers

import (
	"log"

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
	var params models.ProductQueryParams
	response := models.NewResponse(ctx)
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest("params invalid", err.Error())
		return
	}
	log.Println("[debug query params]", params)
	res, err := h.repo.GetAllProducts(ctx.Request.Context(), &params)
	if err != nil {
		response.InternalServerError("internal server errors", err.Error())
		return
	}
	response.Success("get products success", res)
}

func (h *ProductHandlers) FetchDetailProductHandler(ctx *gin.Context) {
	response := models.NewResponse(ctx)
	var err error
	if response == nil {
		response.InternalServerError("internal server error", err.Error())
		return
	}

	id := ctx.Param("id")
	detail, err := h.repo.GetDetailProduct(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound("product not available", err.Error())
		return
	}

	// Fetch recommended products
	recommended, err := h.repo.GetRecommendation(ctx.Request.Context(), 9)
	if err != nil {
		response.InternalServerError("failed to get recommended products", err.Error())
		return
	}

	payload := gin.H{
		"detail":      detail,
		"recommended": recommended,
	}

	response.Success("get product detail with recommendation success", payload)
}
