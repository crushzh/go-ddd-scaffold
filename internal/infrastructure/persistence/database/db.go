package database

import (
	"fmt"
	"time"

	"go-ddd-scaffold/pkg/config"
	"go-ddd-scaffold/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB manages database connections
type DB struct {
	db *gorm.DB
}

// NewDB creates a database connection
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "sqlite":
		dialector = sqlite.Open(cfg.Path)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	logLevel := gormlogger.Warn
	if cfg.Type == "sqlite" {
		logLevel = gormlogger.Silent
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// SQLite WAL mode
	if cfg.Type == "sqlite" {
		db.Exec("PRAGMA journal_mode=WAL")
		db.Exec("PRAGMA busy_timeout=5000")
	}

	// Connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	logger.Infof("database connected: %s", cfg.Type)
	return &DB{db: db}, nil
}

// GormDB returns the underlying GORM instance
func (d *DB) GormDB() *gorm.DB {
	return d.db
}

// AutoMigrate runs auto-migration
func (d *DB) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

// Close closes the database connection
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
