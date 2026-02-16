package service

import (
	"go-ddd-scaffold/internal/application/dto"
	"go-ddd-scaffold/internal/domain/example"
)

// ExampleAppService orchestrates example domain logic
type ExampleAppService struct {
	repo example.Repository
}

// NewExampleAppService creates a new application service
func NewExampleAppService(repo example.Repository) *ExampleAppService {
	return &ExampleAppService{repo: repo}
}

// Create creates a new example
func (s *ExampleAppService) Create(req *dto.CreateExampleRequest) (*dto.ExampleResponse, error) {
	entity := example.NewExample(req.Name, req.Description)

	if err := s.repo.Save(entity); err != nil {
		return nil, err
	}

	return dto.FromExample(entity), nil
}

// GetByID returns an example by ID
func (s *ExampleAppService) GetByID(id uint) (*dto.ExampleResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return dto.FromExample(entity), nil
}

// List returns paginated examples
func (s *ExampleAppService) List(req *dto.QueryExampleRequest) ([]*dto.ExampleResponse, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	entities, total, err := s.repo.List(req.Page, req.PageSize, req.Keyword, example.Status(req.Status))
	if err != nil {
		return nil, 0, err
	}

	return dto.FromExampleList(entities), total, nil
}

// Update updates an example
func (s *ExampleAppService) Update(id uint, req *dto.UpdateExampleRequest) (*dto.ExampleResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Call domain methods
	if req.Name != nil || req.Description != nil {
		name := ""
		desc := ""
		if req.Name != nil {
			name = *req.Name
		}
		if req.Description != nil {
			desc = *req.Description
		}
		entity.UpdateInfo(name, desc)
	}

	if req.Status != nil {
		switch example.Status(*req.Status) {
		case example.StatusActive:
			entity.Activate()
		case example.StatusInactive:
			entity.Deactivate()
		}
	}

	if err := s.repo.Save(entity); err != nil {
		return nil, err
	}

	return dto.FromExample(entity), nil
}

// Delete deletes an example
func (s *ExampleAppService) Delete(id uint) error {
	return s.repo.Delete(id)
}
