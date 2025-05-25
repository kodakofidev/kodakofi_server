package handlers

import (
	"log"
	"net/http"

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
