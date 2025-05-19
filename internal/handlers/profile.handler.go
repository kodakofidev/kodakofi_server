package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

type ProfileHandlers struct {
	repo repositories.ProfileRepoInterface
}

func NewProfileHandlers(repo repositories.ProfileRepoInterface) *ProfileHandlers {
	return &ProfileHandlers{repo: repo}
}

func (h *ProfileHandlers) FetchProfileHandler(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	rawClaims, exists := ctx.Get("payloads")
	var err error

	if !exists {
		responder.Unauthorized("authentication required", any(err))
		return
	}

	claims, ok := rawClaims.(*pkg.Claims)
	if !ok || claims == nil {
		responder.Unauthorized("invalid authentication claims", err)
		return
	}

	user, err := h.repo.GetProfile(ctx.Request.Context(), claims.Uuid)
	if err != nil {
		responder.NotFound("Profile Not Found", err.Error())
		return
	}

	responder.Success("Succes", user)
}

func (h *ProfileHandlers) EditProfileHandler(ctx *gin.Context) {
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
	// userId := "77abbee4-11d8-4db2-9f05-769363454b61"

	var formBody models.ProfileForm
	if err := ctx.ShouldBind(&formBody); err != nil {
		responder.BadRequest("Binding Error", err)
		return
	}

	var profileImageURL string
	if formBody.ProfileImage != nil {
		filename, filepath, err := h.handleFileUpload(ctx, formBody.ProfileImage, userId)
		if err != nil {
			responder.InternalServerError("Error", "Internal server error")
			return
		}

		log.Println("[DEBUG] FILE NAME CHECK", filename)
		profileImageURL = filepath
	}

	// log.Println("[DEBUG] IMAGE URL", profileImageURL)

	result, err := h.repo.EditProfile(ctx.Request.Context(), userId, formBody, profileImageURL)
	if err != nil {
		responder.InternalServerError("Failed to edit profile", err)
		return
	}

	responder.Success("Profile Updated succesfully!", result)
}

// UPLOAD IMAGE HANDLER
func (h *ProfileHandlers) handleFileUpload(ctx *gin.Context, file *multipart.FileHeader, userID string) (filename, filePath string, err error) {
	//First delete any existing profile image for this very user
	oldFiles, err := filepath.Glob(filepath.Join("public", "profile-images", "*_"+userID+"_profile*"))
	if err != nil {
		return "", "", fmt.Errorf("failed to check for existing files: %w", err)
	}

	for _, oldFile := range oldFiles {
		if err := os.Remove(oldFile); err != nil {
			log.Printf("Warning: failed to delete old file %s: %v", oldFile, err)
		}
	}

	//GENERATE FILE NAME AND ADDING PATH

	ext := filepath.Ext(file.Filename)
	filename = fmt.Sprintf("%d_%s_profile%s", time.Now().UnixNano(), userID, ext)
	filePath = filepath.Join("public", "profile-images", filename)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, filename, nil
}
