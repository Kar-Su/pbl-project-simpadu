package middlewares

import (
	"net/http"
	"slices"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleName := ctx.MustGet("role_name").(string)

		found := slices.Contains(allowedRoles, roleName)

		if !found {
			res := utils.BuildResponseFailed("Role anda tidak diizinkan", "Forbidden", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}
		ctx.Next()
	}
}
