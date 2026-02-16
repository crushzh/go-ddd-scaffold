package router

import (
	"go-ddd-scaffold/internal/container"
	"go-ddd-scaffold/internal/interfaces/http/handler"
	"go-ddd-scaffold/internal/interfaces/http/middleware"
	"go-ddd-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
)

// Setup initializes the HTTP router
func Setup(c *container.Container) *gin.Engine {
	if c.Config.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(
		middleware.Recovery(),
		middleware.CORS(),
		middleware.RequestID(),
		middleware.Logger(),
	)

	// Health check
	r.GET("/health", func(ctx *gin.Context) {
		response.OK(ctx)
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth (public)
		auth := v1.Group("/auth")
		{
			authHandler := handler.NewAuthHandler(&c.Config.JWT)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", handler.AuthMiddleware(&c.Config.JWT), authHandler.RefreshToken)
		}

		// Authenticated routes
		authorized := v1.Group("")
		authorized.Use(handler.AuthMiddleware(&c.Config.JWT))
		{
			// Example module
			exampleHandler := handler.NewExampleHandler(c.ExampleService)
			examples := authorized.Group("/examples")
			{
				examples.GET("", exampleHandler.List)
				examples.POST("", exampleHandler.Create)
				examples.GET("/:id", exampleHandler.Get)
				examples.PUT("/:id", exampleHandler.Update)
				examples.DELETE("/:id", exampleHandler.Delete)
			}

			// GEN:ROUTE_REGISTER - Code generator appends routes here, do not remove
		}
	}

	return r
}
