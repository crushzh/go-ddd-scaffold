package database

import (
	"time"

	"go-ddd-scaffold/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserModel is the database model for users (auth only)
type UserModel struct {
	ID        uint   `gorm:"primarykey"`
	Username  string `gorm:"uniqueIndex;size:50;not null"`
	Password  string `gorm:"size:255;not null"`
	Role      string `gorm:"size:20;default:admin"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

// EnsureDefaultAdmin creates the default admin user if no users exist
func EnsureDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&UserModel{}).Count(&count)
	if count > 0 {
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("failed to hash default password: %v", err)
		return
	}

	admin := &UserModel{
		Username: "admin",
		Password: string(hashed),
		Role:     "admin",
	}
	if err := db.Create(admin).Error; err != nil {
		logger.Errorf("failed to create default admin: %v", err)
		return
	}
	logger.Info("default admin created: admin / admin123")
}
