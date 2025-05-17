package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

func (m *Middleware) AccsessGate(allowedRole ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		responder := models.NewResponse(ctx)
		var err error

		payloads, exits := ctx.Get("payloads")
		if !exits {
			responder.Unauthorized("Please login first!", err.Error())
			return
		}
		userPayload, ok := payloads.(*pkg.Claims)
		if !ok {
			responder.Unauthorized("Your login identity is malformed, please login again!", err.Error())
			return
		}
		if !slices.Contains(allowedRole, userPayload.Role) {
			responder.Forbidden("You do not have permission to access", err.Error())
			return
		}
		ctx.Next()
	}
}
