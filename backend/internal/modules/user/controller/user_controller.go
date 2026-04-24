package controller

import (
	"web-hosting/internal/modules/user/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserController interface {
	GetUserAdmin(ctx *gin.Context)
	GetUserNonAdmin(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	GetUserByRole(ctx *gin.Context)
	UpdateAdmin(ctx *gin.Context)
	UpdateNonAdmin(ctx *gin.Context)
	DeleteAdmin(ctx *gin.Context)
	DeleteNonAdmin(ctx *gin.Context)
}

type userController struct {
	userService    service.UserService
	userValidation *validator.Validate
	db             *gorm.DB
}
