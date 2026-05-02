package main

import (
	"log"
	"web-hosting/internal/modules/auth"
	"web-hosting/internal/modules/role"
	"web-hosting/internal/modules/user"
	"web-hosting/internal/package/env"
	"web-hosting/internal/providers"

	_ "github.com/common-nighthawk/go-figure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	_ "web-hosting/docs" // Ganti dengan path modul di go.mod kamu

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func run(server *gin.Engine) {
	// server.Static("/assets", "./assets")
	port := env.GetWithDefault[string]("GO_PORT", "8080")

	var serve string
	if env.GetWithDefault[string]("GO_APP", "localhost") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

// @title           TIM 1 API
// @version         1.0
// @description     apakek.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost
// @BasePath  /api

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Format: Bearer <access token>
func main() {
	injector := do.New()

	server := gin.Default()
	server.Use(cors.Default())

	providers.RegisterProviders(injector)

	user.RegisterRoutes(server, injector)
	auth.RegisterRoutes(server, injector)
	role.RegisterRoutes(server, injector)

	server.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	run(server)
}
