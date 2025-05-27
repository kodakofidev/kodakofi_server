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
	responder := models.NewResponse(ctx)
	productID := ctx.Param("id")

	// Bind form data
	var updateData models.ProductRequest
	if err := ctx.ShouldBind(&updateData); err != nil {
		responder.BadRequest("Invalid request data", gin.H{"error": err.Error()})
		return
	}

	log.Println("[DEBUG PRODUCT]", productID)

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

	log.Println("[DEBUG]", imagesname)
	if len(imagesname) != 0{

	}

	// // Get current images from database
	// currentImages, err := h.repo.GetListImageProduct(ctx.Request.Context(), productID)
	// if err != nil {
	// 	responder.InternalServerError("Failed to get current images", nil)
	// 	return
	// }

	// // Process file uploads
	// var newImageNames []string
	// form, err := ctx.MultipartForm()
	// if err == nil { // Jika ada form file
	// 	imageFiles := form.File["images"]

	// 	// Validasi maksimal 3 gambar
	// 	if len(imageFiles) > 3 {
	// 		responder.BadRequest("Maximum 3 images allowed", nil)
	// 		return
	// 	}

	// 	for _, file := range imageFiles {
	// 		if !isImage(file) {
	// 			responder.BadRequest("Invalid image file type", nil)
	// 			return
	// 		}

	// 		filename, _, err := fileHandling(ctx, file)
	// 		if err != nil {
	// 			responder.InternalServerError("Failed to save image", nil)
	// 			return
	// 		}
	// 		newImageNames = append(newImageNames, filename)
	// 	}
	// }
	// shouldUpdateImages := false
	// if len(newImageNames) > 0 {
	// 	shouldUpdateImages = true
	// } else {
	// 	// Cek apakah ada request untuk menghapus gambar (misal: keep_images kosong)
	// 	if updateData.KeepImages != nil && len(updateData.KeepImages) < len(currentImages) {
	// 		shouldUpdateImages = true
	// 	}
	// }

	// // Update product data
	// err = h.repo.UpdateProduct(ctx.Request.Context(), productID, &updateData, newImageNames, shouldUpdateImages, currentImages) // Kirim current images ke repository

	// if err != nil {
	// 	// Cleanup: hapus gambar baru jika update gagal
	// 	for _, img := range newImageNames {
	// 		os.Remove(filepath.Join("public", "product-image", img))
	// 	}
	// 	responder.InternalServerError("Failed to update product", nil)
	// 	return
	// }

	// // Hapus gambar lama jika ada gambar baru
	// if shouldUpdateImages && len(newImageNames) > 0 {
	// 	for _, oldImg := range currentImages {
	// 		os.Remove(filepath.Join("public", "product-image", oldImg))
	// 	}
	// }

	// responder.Success("Product updated successfully", gin.H{
	// 	"product_id":     productID,
	// 	"images_updated": shouldUpdateImages,
	// })
}
