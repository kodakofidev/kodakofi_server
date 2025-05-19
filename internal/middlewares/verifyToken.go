package middlewares

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

func (a *Middleware) VerifyToken(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	responder := models.NewResponse(ctx)

	if bearerToken == "" {
		responder.Unauthorized("Token Not found", any(nil))
		return
	}

	if !strings.HasPrefix(bearerToken, "Bearer ") {
		responder.Unauthorized("Token Not valid", any(nil))
		return
	}
	token := strings.Split(bearerToken, " ")[1]
	if token == "" {
		responder.Unauthorized("Unauthorized", any(nil))
		return
	}
	payloads := &pkg.Claims{}
	if err := payloads.VerifyToken(token); err != (pkg.JWTErr{}) {
		log.Println(err.Err.Error())
		if err.Type == "Token" {
			responder.Unauthorized("Unauthorized", err.Err.Error())
			return
		}
		responder.InternalServerError("Internal server error", err.Err.Error())
		return
	}
	ctx.Set("payloads", payloads)
	ctx.Next()
}
