package middlewares

import (
	"net/http"
	"strings"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			res := utils.BuildResponseFailed("Failed Auth", "Authorization header missing", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			res := utils.BuildResponseFailed("Failed Auth", "Invalid authorization header", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		authHeader = authHeader[len("Bearer "):]
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			res := utils.BuildResponseFailed("Failed Auth", "Invalid token", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !token.Valid {
			res := utils.BuildResponseFailed("Failed Auth", "Invalid token", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			res := utils.BuildResponseFailed("Failed Auth", "Invalid token", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		roleName, err := jwtService.GetRoleNameByToken(authHeader)
		if err != nil {
			res := utils.BuildResponseFailed("Failed Auth", "Invalid token", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set("user_id", userId)
		ctx.Set("role_name", roleName)
		ctx.Set("token", authHeader)
		ctx.Next()
	}
}
