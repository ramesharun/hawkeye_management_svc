package monitor

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ramesharun/hawkeye-management-svc/internal/entity"
	"github.com/ramesharun/hawkeye-management-svc/pkg/log"
)

// Service encapsulates usecase logic for monitors.
type Service interface {
	Get(ctx context.Context, id string) (Monitor, error)
	GetByOrg(ctx context.Context, org_id string, offset, limit int) ([]Monitor, error)
	GetByTenant(ctx context.Context, tenant string, offset, limit int) ([]Monitor, error)
	Query(ctx context.Context, offset, limit int) ([]Monitor, error)
	Count(ctx context.Context) (int, error)
	CountByOrg(ctx context.Context, org_id string) (int, error)
	CountByTenant(ctx context.Context, tenant string) (int, error)
	Create(ctx context.Context, input CreateMonitorRequest) (Monitor, error)
	Update(ctx context.Context, id string, input UpdateMonitorRequest) (Monitor, error)
	Delete(ctx context.Context, id string) (Monitor, error)
}

// Monitor represents the data about an monitor.
type Monitor struct {
	entity.Apichecks
}

// CreateMonitorRequest represents an monitor creation request.
type CreateMonitorRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateMonitorRequest fields.
func (m CreateMonitorRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateMonitorRequest represents an monitor update request.
type UpdateMonitorRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateMonitorRequest fields.
func (m UpdateMonitorRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new monitor service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the monitor with the specified the monitor ID.
func (s service) Get(ctx context.Context, id string) (Monitor, error) {
	monitor, err := s.repo.Get(ctx, id)
	if err != nil {
		return Monitor{}, err
	}
	return Monitor{monitor}, nil
}

// Create creates a new monitor.
func (s service) Create(ctx context.Context, req CreateMonitorRequest) (Monitor, error) {
	if err := req.Validate(); err != nil {
		return Monitor{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.Apichecks{
		ID:        id,
		Name:      req.Name,
		UpdatedAt: now,
	})
	if err != nil {
		return Monitor{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the monitor with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateMonitorRequest) (Monitor, error) {
	if err := req.Validate(); err != nil {
		return Monitor{}, err
	}

	monitor, err := s.Get(ctx, id)
	if err != nil {
		return monitor, err
	}
	monitor.Name = req.Name
	monitor.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, monitor.Apichecks); err != nil {
		return monitor, err
	}
	return monitor, nil
}

// Delete deletes the monitor with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Monitor, error) {
	monitor, err := s.Get(ctx, id)
	if err != nil {
		return Monitor{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Monitor{}, err
	}
	return monitor, nil
}

// Count returns the number of monitors.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Count returns the number of monitors.
func (s service) CountByOrg(ctx context.Context, org_id string) (int, error) {
	return s.repo.CountByOrg(ctx, org_id)
}

// Count returns the number of monitors.
func (s service) CountByTenant(ctx context.Context, tenant string) (int, error) {
	return s.repo.CountByTenant(ctx, tenant)
}

// Query returns the monitors with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Monitor, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []Monitor{}
	for _, item := range items {
		result = append(result, Monitor{item})
	}
	return result, nil
}

func (s service) GetByTenant(ctx context.Context, tenant string, offset, limit int) ([]Monitor, error) {
	items, err := s.repo.GetByTenant(ctx, tenant, offset, limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []Monitor{}
	for _, item := range items {
		result = append(result, Monitor{item})
	}
	return result, nil
}

func (s service) GetByOrg(ctx context.Context, tenant string, offset, limit int) ([]Monitor, error) {
	items, err := s.repo.GetByOrg(ctx, tenant, offset, limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []Monitor{}
	for _, item := range items {
		result = append(result, Monitor{item})
	}
	return result, nil
}
