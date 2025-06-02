package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type UserHandlers struct {
	repo repositories.UserRepoInterface
}

func NewUser(repo repositories.UserRepoInterface) *UserHandlers {
	return &UserHandlers{repo: repo}
}

func (h *UserHandlers) FetchAllUsersHandler(ctx *gin.Context) {
	res := models.NewResponse(ctx)
	search := ctx.Query("search") // Ambil dari query string (?search=...)

	users, err := h.repo.GetAllUsers(ctx.Request.Context(), search)
	if err != nil {
		log.Println(http.StatusInternalServerError, err.Error())
		res.InternalServerError("Internal Server Error", "Failed to fetch users")
		return
	}

	res.Success("Fetch users success", users)
}

func (h *UserHandlers) FetchOneUserByAdminHandler(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	userID := ctx.Param("id")
	if userID == "" {
		responder.BadRequest("User ID is required in URL", nil)
		return
	}

	res, err := h.repo.GetOneUserByAdmin(ctx.Request.Context(), userID)
	if err != nil {
		log.Printf("[Handler][FetchOneUserByAdminHandler] error fetching user ID %s: %v", userID, err)
		responder.InternalServerError("Something went wrong while retrieving user data", nil)
		return
	}

	responder.Success("User fetched successfully", res)
}

func (h *UserHandlers) PatchUserByAdminHandler(ctx *gin.Context) {
	responder := models.NewResponse(ctx)

	userID := ctx.Param("id")
	if userID == "" {
		responder.BadRequest("User ID is required in URL", nil)
		return
	}

	var req models.UpdateUserByAdminReq
	if err := ctx.ShouldBind(&req); err != nil {
		responder.BadRequest("Invalid input", err)
		return
	}

	// Handle image upload
	if req.Image != nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), req.Image.Filename)
		savePath := filepath.Join("public/profile-images", filename)

		if err := ctx.SaveUploadedFile(req.Image, savePath); err != nil {
			responder.InternalServerError("Failed to save uploaded image", err)
			return
		}

		req.Image.Filename = filename
	}

	req.ID = userID

	// Call repository
	res, err := h.repo.UpdateUserByAdmin(ctx.Request.Context(), req)
	if err != nil {
		responder.InternalServerError("Failed to update user", err)
		return
	}

	responder.Success("User updated successfully", res)
}
