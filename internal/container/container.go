package container

import (
	"go-ddd-scaffold/internal/application/service"
	"go-ddd-scaffold/internal/infrastructure/persistence/database"
	"go-ddd-scaffold/pkg/config"
	"go-ddd-scaffold/pkg/logger"
)

// Container manages dependency injection
// In DDD architecture, Container assembles cross-layer dependencies
type Container struct {
	Config *config.Config
	DB     *database.DB

	// Application services
	ExampleService *service.ExampleAppService
	// GEN:SERVICE_REGISTER - Code generator appends services here, do not remove
}

// New creates and initializes the container
func New(cfg *config.Config) (*Container, error) {
	c := &Container{Config: cfg}

	// 1. Init infrastructure
	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		return nil, err
	}
	c.DB = db

	// 2. Auto-migrate
	if err := db.AutoMigrate(
		&database.UserModel{},
		&database.ExampleModel{},
		// GEN:MODEL_MIGRATE - Code generator appends models here, do not remove
	); err != nil {
		return nil, err
	}

	// Seed default admin user
	database.EnsureDefaultAdmin(db.GormDB())

	// 3. Create repositories (infra -> domain interface)
	exampleRepo := database.NewExampleRepository(db)

	// 4. Create application services (inject repos)
	c.ExampleService = service.NewExampleAppService(exampleRepo)
	// GEN:SERVICE_INIT - Code generator appends initialization here, do not remove

	logger.Info("DI container initialized")
	return c, nil
}

// Close releases all resources
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
