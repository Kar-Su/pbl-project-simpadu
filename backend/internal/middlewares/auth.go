package middlewares

import (
	"strings"
	"web-hosting/internal/modules/auth/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Authorization header missing"})
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header"})
			return
		}

		authHeader = authHeader[len("Bearer "):]
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		roleName, err := jwtService.GetRoleNameByToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("user_id", userId)
		ctx.Set("role_name", roleName)
		ctx.Set("token", authHeader)
		ctx.Next()
	}
}
