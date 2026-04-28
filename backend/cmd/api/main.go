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

func main() {
	injector := do.New()

	server := gin.Default()
	server.Use(cors.Default())

	providers.RegisterProviders(injector)

	user.RegisterRoutes(server, injector)
	auth.RegisterRoutes(server, injector)
	role.RegisterRoutes(server, injector)
	run(server)
}
