package example

// Repository defines the example repository interface
type Repository interface {
	// FindByID finds by ID
	FindByID(id uint) (*Example, error)

	// List returns paginated results
	List(page, pageSize int, keyword string, status Status) ([]*Example, int64, error)

	// Save creates or updates
	Save(entity *Example) error

	// Delete deletes by ID
	Delete(id uint) error
}
