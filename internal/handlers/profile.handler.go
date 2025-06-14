package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/internal/utils"
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
		responder.Unauthorized("authentication required", err)
		return
	}

	claims, ok := rawClaims.(*pkg.Claims)
	if !ok || claims == nil {
		responder.Unauthorized("invalid authentication claims", err)
		return
	}

	user, err := h.repo.GetProfile(ctx.Request.Context(), claims.Uuid)
	if err != nil {
		responder.NotFound("Profile Not Found", err)
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

	var profileImageFilename string
	if formBody.ProfileImage != nil {
		filename, _, err := utils.FileNameProfile(ctx, formBody.ProfileImage, userId)
		if err != nil {
			responder.InternalServerError("Error", "Internal server error")
			return
		}

		log.Println("[DEBUG] FILE NAME CHECK", filename)
		profileImageFilename = filename
	}

	// log.Println("[DEBUG] IMAGE URL", profileImageURL)

	result, err := h.repo.EditProfile(ctx.Request.Context(), userId, formBody, profileImageFilename)
	if err != nil {
		responder.InternalServerError("Failed to edit profile", err)
		return
	}

	responder.Success("Profile Updated succesfully!", result)
}
