package main

import (
	"log"
	"net/http"
	"web-hosting/internal/package/env"

	_ "github.com/common-nighthawk/go-figure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func run(server *gin.Engine) {
	server.Static("/assets", "./assets")

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
	_ = injector

	server := gin.Default()
	server.Use(cors.Default())

	apiRoutes := server.Group("/api")
	{
		apiRoutes.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello",
			})
		})
	}

	run(server)
}
