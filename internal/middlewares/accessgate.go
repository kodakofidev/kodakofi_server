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

		payloads, exits := ctx.Get("payloads")
		if !exits {
			responder.Unauthorized("Please login first!", any(nil))
			return
		}
		userPayload, ok := payloads.(*pkg.Claims)
		if !ok {
			responder.Unauthorized("Your login identity is malformed, please login again!", any(nil))
			return
		}
		if !slices.Contains(allowedRole, userPayload.Role) {
			responder.Forbidden("You do not have permission to access", any(nil))
			return
		}
		ctx.Next()
	}
}
