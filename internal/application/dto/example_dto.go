package dto

import (
	"time"

	"go-ddd-scaffold/internal/domain/example"
)

// ExampleResponse is the example response DTO
type ExampleResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FromExample converts from domain entity
func FromExample(e *example.Example) *ExampleResponse {
	return &ExampleResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Status:      string(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// FromExampleList converts from domain entity list
func FromExampleList(items []*example.Example) []*ExampleResponse {
	result := make([]*ExampleResponse, len(items))
	for i, item := range items {
		result[i] = FromExample(item)
	}
	return result
}

// CreateExampleRequest is the create request DTO
type CreateExampleRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateExampleRequest is the update request DTO
type UpdateExampleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// QueryExampleRequest is the query request DTO
type QueryExampleRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   string `form:"status" json:"status"`
}
