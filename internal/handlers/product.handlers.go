package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	fp "path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/pkg"
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

func (h *ProductHandlers) AddProduct(ctx *gin.Context) {
	response := models.NewResponse(ctx)
	var formBody models.ProductRequest
	if err := ctx.ShouldBind(&formBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println("Multipart form error:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get uploaded files",
			"error":   err.Error(),
		})
		return
	}

	// Handle multiple images upload
	imageFiles := form.File["images"]
	if len(imageFiles) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "At least one image is required",
		})
		return
	}

	var imagesname []string
	for _, file := range imageFiles {
		// Validate image file
		if !isImage(file) {
			log.Printf("File %s is not an image", file.Filename)
			continue
		}

		filename, _, err := fileHandling(ctx, file)
		if err != nil {
			log.Printf("Failed to upload image: %v", err)
			continue
		}

		imagesname = append(imagesname, filename)
	}

	if len(imagesname) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "No valid images were uploaded",
		})
		return
	}
	errQuery := h.repo.AddProduct(ctx.Request.Context(), &formBody, imagesname)
	if errQuery != nil {
		response.InternalServerError("internal server errors", errQuery)
		return
	}
	response.Success("Add products success", gin.H{
		"name": formBody.Name,
	})
}

func isImage(file *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	fileType := file.Header.Get("Content-Type")
	for _, t := range allowedTypes {
		if fileType == t {
			return true
		}
	}
	return false
}

func fileHandling(ctx *gin.Context, file *multipart.FileHeader) (filename, filepath string, err error) {
	// responder := models.NewResponse(ctx)
	ext := fp.Ext(file.Filename)
	filename = fmt.Sprintf("%d_product%s", time.Now().UnixNano(), ext)
	filepath = fp.Join("public", "product-image", filename)

	if err := os.MkdirAll(fp.Dir(filepath), 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %v", err)
	}

	return filename, filepath, nil
}

func (h *ProductHandlers) UpdateProduct(ctx *gin.Context) {
	response := models.NewResponse(ctx)

	// Get product ID from URL parameter
	productID, ok := ctx.Params.Get("id")
	log.Println("PRODUCTID", productID)
	if !ok {
		log.Println("[DEBUG]", ok)
		response.BadRequest("Invalid product ID", nil)
		return
	}

	// Bind form data
	var formBody models.ProductRequest
	if err := ctx.ShouldBind(&formBody); err != nil {
		response.BadRequest("Invalid request data", gin.H{"error": err.Error()})
		return
	}

	var imagesName []string

	// Handle image upload if provided
	form, err := ctx.MultipartForm()
	if err == nil { // If multipart form exists (images might be uploaded)
		imageFiles := form.File["images"]

		// Process each uploaded file
		for _, file := range imageFiles {
			// Validate image file
			if !isImage(file) {
				log.Printf("File %s is not an image", file.Filename)
				continue
			}

			filename, _, err := fileHandling(ctx, file)
			if err != nil {
				log.Printf("Failed to upload image: %v", err)
				continue
			}

			imagesName = append(imagesName, filename)
		}
	} else {
		log.Println("No images uploaded or form not multipart")
	}
	log.Println("[DEBUG IMAGES]", formBody)

	// Call repository to update product
	err = h.repo.UpdateProduct(ctx.Request.Context(), productID, &formBody, imagesName)
	if err != nil {
		log.Printf("Update product error: %v", err)
		response.InternalServerError("Failed to update product", err)
		return
	}

	response.Success("Product updated successfully", gin.H{
		"product_id": productID,
		"name":       formBody.Name,
	})
}

func (h *ProductHandlers) UpdateLikeProduct(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	rawClaims, exists := ctx.Get("payloads")
	var err error

	if !exists {
		responder.Unauthorized("authentication required", err)
		return
	}

	claims, ok := rawClaims.(*pkg.Claims)
	if !ok || claims == nil || claims.Uuid == "" {
		responder.Unauthorized("invalid authentication claims", "Missing or invalid JWT claims")
		return
	}
	userId := claims.Uuid

	productId, ok := ctx.Params.Get("id")
	if !ok {
		responder.BadRequest("params needed Error", err)
		return
	}

	isExist, err := h.repo.ToggleLike(ctx.Request.Context(), userId, productId)
	if err != nil {
		responder.InternalServerError("params needed Error", err)
		return
	}

	responder.Success("success", isExist)

}
func (h *ProductHandlers) GetLikeProducts(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	rawClaims, exists := ctx.Get("payloads")
	var err error

	if !exists {
		responder.Unauthorized("authentication required", err)
		return
	}

	claims, ok := rawClaims.(*pkg.Claims)
	if !ok || claims == nil || claims.Uuid == "" {
		responder.Unauthorized("invalid authentication claims", "Missing or invalid JWT claims")
		return
	}
	userId := claims.Uuid

	productId, ok := ctx.Params.Get("id")
	if !ok {
		responder.BadRequest("params needed Error", err)
		return
	}

	isExist, err := h.repo.GetLikeStatus(ctx.Request.Context(), userId, productId)
	if err != nil {
		responder.InternalServerError("params needed Error", err)
		return
	}

	responder.Success("success", isExist)
}

func (h *ProductHandlers) FetchAllProductsAdminHandler(ctx *gin.Context) {
	var params models.ProductQueryParams
	response := models.NewResponse(ctx)

	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest("params invalid", err.Error())
		return
	}
	res, err := h.repo.GetAllProductsAdmin(ctx.Request.Context(), &params)
	if err != nil {
		response.InternalServerError("internal server errors", err.Error())
		return
	}

	response.Success("get products success", res)
}
