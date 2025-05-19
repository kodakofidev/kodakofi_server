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

func (h *ProductHandlers) AddProduct(ctx *gin.Context) {
	var formBody models.ProductRequest
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Println("Binding error:", err.Error())
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

	var imagename []string
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

		imagename = append(imagename, filename)
	}

	if len(imagename) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "No valid images were uploaded",
		})
		return
	}
	log.Println(imagename)
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
	log.Println("[DEBUG ext]", ext)
	filename = fmt.Sprintf("%d_product%s", time.Now().UnixNano(), ext)
	log.Println("[DEBUG FILE NAME]", filename)
	filepath = fp.Join("public", "img", "product", filename)

	if err := os.MkdirAll(fp.Dir(filepath), 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %v", err)
	}

	return filename, filepath, nil
}
