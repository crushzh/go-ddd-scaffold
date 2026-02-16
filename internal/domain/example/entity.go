package example

import "time"

// Example is the aggregate root
type Example struct {
	ID          uint
	Name        string
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Status is a value object
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

// IsValid checks if status is valid
func (s Status) IsValid() bool {
	return s == StatusActive || s == StatusInactive
}

// String returns string representation
func (s Status) String() string {
	return string(s)
}

// NewExample creates a new Example (factory method)
func NewExample(name, description string) *Example {
	return &Example{
		Name:        name,
		Description: description,
		Status:      StatusActive,
	}
}

// Activate sets status to active
func (e *Example) Activate() {
	e.Status = StatusActive
}

// Deactivate sets status to inactive
func (e *Example) Deactivate() {
	e.Status = StatusInactive
}

// UpdateInfo updates basic info
func (e *Example) UpdateInfo(name, description string) {
	if name != "" {
		e.Name = name
	}
	if description != "" {
		e.Description = description
	}
}
