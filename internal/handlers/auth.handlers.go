package handlers

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
)

type AuthHandlers struct {
	repo repositories.AuthRepoInterface
}

func NewAuthHandlers(repo repositories.AuthRepoInterface) *AuthHandlers {
	return &AuthHandlers{repo: repo}
}

func (h *AuthHandlers) Login(ctx *gin.Context) {

}

func (h *AuthHandlers) Register(ctx *gin.Context) {

}

func (h *AuthHandlers) Logout(ctx *gin.Context) {

}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
