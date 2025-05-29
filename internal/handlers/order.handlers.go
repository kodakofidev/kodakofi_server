package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/pkg"
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

// handlers get history order
func (h *OrderHandlers) GetHistoryOrders(ctx *gin.Context) {
	claims, _ := ctx.Get("payloads")
	userClaims := claims.(*pkg.Claims)

	response := models.NewResponse(ctx)

	// tangkap query
	pageQ := ctx.Query("page")
	statusQ := ctx.Query("status")
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
		offset = pageQInt*4 - 4
	}

	log.Println("offset", offset)
	log.Println("statusQ", statusQ)

	result, err := h.repo.GetHistoryOrders(ctx, offset, statusQ, userClaims.Uuid)
	if err != nil {
		response.InternalServerError("Internal Server Error", "a server error occured")
		return
	}
	println(len(result))
	if len(result) == 0 {
		response.NotFound("Not Found", "history order is empty")
		return
	}

	response.Success("success", result)
}

func (h *OrderHandlers) FetchDetailOrdersHandler(ctx *gin.Context) {
	res := models.NewResponse(ctx)

	// Ambil user_id dari JWT payload
	payloads, exists := ctx.Get("payloads")
	if !exists {
		res.Unauthorized("Please login first", nil)
		return
	}
	userClaims, ok := payloads.(*pkg.Claims)
	if !ok {
		res.Unauthorized("Malformed login identity", nil)
		return
	}
	userID := userClaims.Uuid

	// Ambil transaction_code dari query param
	transactionCode := ctx.Param("transaction_code")
	if transactionCode == "" {
		res.BadRequest("Transaction code is required", nil)
		return
	}

	// Ambil data dari repository
	orderDetail, err := h.repo.GetDetailOrder(ctx, userID, transactionCode)
	if err != nil {
		res.InternalServerError("Failed to fetch detail order", nil)
		return
	}

	res.Success("Fetch order details successed", orderDetail)
}

func (h *OrderHandlers) FetchDetailOrderAdminHandler(ctx *gin.Context) {
	res := models.NewResponse(ctx)

	transactionCode := ctx.Param("transaction_code")
	if transactionCode == "" {
		res.BadRequest("Transaction code is required", nil)
		return
	}

	orderDetail, err := h.repo.GetDetailOrderAdmin(ctx, transactionCode)
	if err != nil {
		res.InternalServerError("Failed to fetch detail order", nil)
		return
	}

	res.Success("Fetch order details successed", orderDetail)
}

func (h *OrderHandlers) FetchDataSalesHandler(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	responder := models.NewResponse(ctx)

	if startDateStr == "" || endDateStr == "" {
		log.Println("Bad Request:", "start_date and end_date query parameters are required")
		responder.BadRequest("Bad Request", "date parameters are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Println("Bad Request:", err)
		responder.BadRequest("Bad Request", "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Println("Bad Request:", err)
		responder.BadRequest("Bad Request", "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	data, err := h.repo.GetDataSales(ctx, startDate, endDate)
	if err != nil {
		log.Println("Internal Server Error:", err)
		responder.InternalServerError("Failed to fetch sales data", err.Error())
		return
	}

	responder.Success("Sales data fetched successfully", data)
}

func (h *OrderHandlers) UpdateOrderStatus(ctx *gin.Context) {
	response := models.NewResponse(ctx)

	var req models.UpdateOrderStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Bad Request:", err)
		response.BadRequest("Bad Request", "invalid request body")
		return
	}

	res, err := h.repo.UpdateStatusOrder(ctx.Request.Context(), req.OrderID, req.StatusID)
	if err != nil {
		log.Println("Internal Server Error:", err)
		response.InternalServerError("Internal Server Error:", "Failed to update order status")
		return
	}

	response.Success("Order status updated successfully", res)
}
