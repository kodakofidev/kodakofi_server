package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) CORSMiddleware(ctx *gin.Context) {
	// Setup whitelist origin
	whitelistOrigin := []string{"http://localhost:5173"} // Add more origins as needed
	origin := ctx.GetHeader("Origin")

	// Always set CORS headers
	if origin != "" && slices.Contains(whitelistOrigin, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	} else {
		ctx.Header("Access-Control-Allow-Origin", "*") // Allow all origins during development
	}

	ctx.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, PATCH, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Max-Age", "86400") // 24 hours

	// Handle preflight requests
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
