package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

func (a *Middleware) VerifyToken(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "Unauthorized",
		})
		return
	}

	if !strings.HasPrefix(bearerToken, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "Unauthorized",
		})
		return
	}
	token := strings.Split(bearerToken, " ")[1]
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "Unauthorized",
		})
		return
	}
	payloads := &pkg.Claims{}
	if err := payloads.VerifyToken(token); err != (pkg.JWTErr{}) {
		log.Println(err.Err.Error())
		if err.Type == "Token" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": err.Err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
		return
	}
	ctx.Set("payloads", payloads)
	ctx.Next()
}
