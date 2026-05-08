package main

import (
	"log"
	"web-hosting/internal/modules/akademik"
	"web-hosting/internal/modules/auth"
	"web-hosting/internal/modules/jurusan"
	"web-hosting/internal/modules/kurikulum"
	"web-hosting/internal/modules/mk"
	"web-hosting/internal/modules/prodi"
	"web-hosting/internal/modules/role"
	"web-hosting/internal/modules/user"
	"web-hosting/internal/package/env"
	"web-hosting/internal/providers"

	_ "github.com/common-nighthawk/go-figure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	"web-hosting/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var SWAGGER_PATH string

func run(server *gin.Engine) {
	// server.Static("/assets", "./assets")
	port := env.GetWithDefault[string]("GO_PORT", "8080")

	var serve string
	app := env.GetWithDefault[string]("GO_APP", "localhost")
	if app == "localhost" {
		serve = "0.0.0.0:" + port
		docs.SwaggerInfo.Host = "localhost"
		SWAGGER_PATH = "/api/swagger/*any"
		server.Use(cors.Default())
	} else {
		serve = ":" + port
		docs.SwaggerInfo.Host = app
		SWAGGER_PATH = "/swagger/*any"
	}
	log.Printf("API docs url: %v", SWAGGER_PATH)
	server.GET(SWAGGER_PATH, ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("server is running on %s", app)
	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %s", err)
	}

}

// @title           TIM 1 API
// @version         1.0
// @description     **BARANG SIAPA YANG UPDATE/HAPUS DATA SEMBARANGAN DAPAT KECUPAN DARI HAFIZ**.
// @description     **BUAT DATA BARU JIKA INGIN HAPUS/UPDATE**.
// @description
// @description     **DUMMY DATA LINK:**
// @description     **https://github.com/Kar-Su/pbl-project-simpadu/tree/main/backend/internal/database/seeders/json**
// @description
// @description 	**TUTORIAL LINK:***
// @description     **https://github.com/Kar-Su/pbl-project-simpadu**
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Format: Bearer <access token>
func main() {
	injector := do.New()

	server := gin.Default()

	providers.RegisterProviders(injector)

	user.RegisterRoutes(server, injector)
	auth.RegisterRoutes(server, injector)
	role.RegisterRoutes(server, injector)
	jurusan.RegisterRoutes(server, injector)
	prodi.RegisterRoutes(server, injector)
	mk.RegisterRoutes(server, injector)
	akademik.RegisterRoutes(server, injector)
	kurikulum.RegisterRoutes(server, injector)

	run(server)
}
