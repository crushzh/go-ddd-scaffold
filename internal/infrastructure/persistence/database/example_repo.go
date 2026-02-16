package database

import (
	"go-ddd-scaffold/internal/domain/example"

	"gorm.io/gorm"
)

// ExampleRepository implements example.Repository
type ExampleRepository struct {
	db *gorm.DB
}

// NewExampleRepository creates a new repository
func NewExampleRepository(database *DB) example.Repository {
	return &ExampleRepository{db: database.GormDB()}
}

// FindByID finds by ID
func (r *ExampleRepository) FindByID(id uint) (*example.Example, error) {
	var model ExampleModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// List returns paginated results
func (r *ExampleRepository) List(page, pageSize int, keyword string, status example.Status) ([]*example.Example, int64, error) {
	var models []ExampleModel
	var total int64

	query := r.db.Model(&ExampleModel{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", string(status))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*example.Example, len(models))
	for i := range models {
		entities[i] = models[i].ToDomain()
	}

	return entities, total, nil
}

// Save creates or updates
func (r *ExampleRepository) Save(entity *example.Example) error {
	model := FromDomain(entity)
	if model.ID == 0 {
		if err := r.db.Create(model).Error; err != nil {
			return err
		}
		entity.ID = model.ID
		entity.CreatedAt = model.CreatedAt
		entity.UpdatedAt = model.UpdatedAt
		return nil
	}
	return r.db.Save(model).Error
}

// Delete deletes by ID
func (r *ExampleRepository) Delete(id uint) error {
	return r.db.Delete(&ExampleModel{}, id).Error
}
