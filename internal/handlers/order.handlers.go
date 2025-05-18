package handlers

import (
	"errors"
	"log"
	"strconv"

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

}

// handlers get history order
func (h *OrderHandlers) GetHistoryOrders(ctx *gin.Context) {
	// claims, _ := ctx.Get("payloads")
	// userClaims := claims.(*pkg.Claims)

	response := models.NewResponse(ctx)

	userId := "416e8dea-843e-4de9-85cb-ab1dcf4f378f"
	
	// tangkap query
	pageQ := ctx.Query("page");
	statusQ := ctx.Query("status");
	var offset int
	var pageQInt int
	if pageQ != "" {
		pageQNum, err := strconv.Atoi(pageQ)
		if err != nil {
			response.InternalServerError("a server error occured", err.Error())
			return
		}
		pageQInt += pageQNum
	}

	if pageQInt == 1 {
		offset = 0
	} else if pageQInt == 0 {
		offset = -1
	} else {
		offset = pageQInt * 4 - 4
	}

	log.Println("offset", offset)
	log.Println("statusQ", statusQ)

	result, err := h.repo.GetHistoryOrders(ctx, offset, statusQ, userId)
	if err != nil {
		response.InternalServerError("a server error occured", err.Error())
		return
	}

	if len(result) == 0 {
		response.NotFound("history order not found", errors.New("history order is empty"))
		return
	}

	response.Success("success", result)
}
