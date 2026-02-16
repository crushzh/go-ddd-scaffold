package database

import (
	"time"

	"go-ddd-scaffold/internal/domain/example"
)

// ExampleModel is the GORM model (infrastructure layer)
type ExampleModel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null;index"`
	Description string `gorm:"size:500"`
	Status      string `gorm:"size:20;default:active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName overrides the table name
func (ExampleModel) TableName() string {
	return "examples"
}

// ToDomain converts to domain entity
func (m *ExampleModel) ToDomain() *example.Example {
	return &example.Example{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Status:      example.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromDomain converts from domain entity
func FromDomain(e *example.Example) *ExampleModel {
	return &ExampleModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Status:      string(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
