package router

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"go-ddd-scaffold/internal/container"
	"go-ddd-scaffold/internal/interfaces/http/handler"
	"go-ddd-scaffold/internal/interfaces/http/middleware"
	"go-ddd-scaffold/internal/web"
	"go-ddd-scaffold/pkg/logger"
	"go-ddd-scaffold/pkg/response"

	_ "go-ddd-scaffold/docs/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Swagger API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	// ====== Frontend static files ======
	registerFrontendRoutes(r)

	return r
}

// registerFrontendRoutes registers frontend static file routes (SPA support)
func registerFrontendRoutes(r *gin.Engine) {
	// Method 1: go:embed
	distFS, err := web.GetDistFS()
	if err == nil {
		registerEmbeddedFrontend(r, distFS)
		logger.Info("frontend files embedded in binary")
		return
	}

	// Method 2: external dist directory
	for _, dir := range []string{"dist", "web/dist", "static"} {
		absDir, _ := filepath.Abs(dir)
		if isDir(absDir) {
			r.Static("/assets", filepath.Join(absDir, "assets"))
			r.StaticFile("/", filepath.Join(absDir, "index.html"))
			r.NoRoute(spaFallback(absDir))
			logger.Infof("frontend directory: %s", absDir)
			return
		}
	}

	logger.Warn("frontend files not found (ignore in dev mode)")
}

func registerEmbeddedFrontend(r *gin.Engine, distFS fs.FS) {
	httpFS := http.FS(distFS)

	// Static assets (long-term cache)
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.FileFromFS(c.Request.URL.Path, httpFS)
	})

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API and Swagger routes return 404
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			response.NotFound(c, "API not found")
			return
		}

		// Try to find the file
		if f, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
			f.Close()
			c.FileFromFS(path, httpFS)
			return
		}

		// Fall back to index.html
		if indexFile, err := fs.ReadFile(distFS, "index.html"); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexFile)
			return
		}

		response.NotFound(c, "page not found")
	})
}

func spaFallback(staticDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			response.NotFound(c, "API not found")
			return
		}
		c.File(filepath.Join(staticDir, "index.html"))
	}
}

func isDir(path string) bool {
	entries, err := filepath.Glob(filepath.Join(path, "*"))
	return err == nil && len(entries) > 0
}
