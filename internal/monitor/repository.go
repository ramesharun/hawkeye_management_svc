package monitor

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/ramesharun/hawkeye-management-svc/internal/entity"
	"github.com/ramesharun/hawkeye-management-svc/pkg/dbcontext"
	"github.com/ramesharun/hawkeye-management-svc/pkg/log"
)

// Repository encapsulates the logic to access monitors from the data source.
type Repository interface {
	// Get returns the monitor with the specified row ID.
	Get(ctx context.Context, id string) (entity.Apichecks, error)
	// GetByOrg returns the monitor with the specified org ID.
	GetByOrg(ctx context.Context, org_id string, offset, limit int) ([]entity.Apichecks, error)

	GetByTenant(ctx context.Context, tenant string, offset, limit int) ([]entity.Apichecks, error)
	// Count returns the number of monitors.
	Count(ctx context.Context) (int, error)
	// Query returns the list of monitors with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Apichecks, error)
	// Create saves a new monitor in the storage.
	Create(ctx context.Context, monitor entity.Apichecks) error
	// Update updates the monitor with given ID in the storage.
	Update(ctx context.Context, monitor entity.Apichecks) error
	// Delete removes the monitor with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists monitors in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new monitor repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the monitor with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Apichecks, error) {
	var monitor entity.Apichecks
	err := r.db.With(ctx).Select("*").From("hawkeye.apichecks").Model(id, &monitor)
	return monitor, err
}

// Create saves a new monitor record in the database.
// It returns the ID of the newly inserted monitor record.
func (r repository) Create(ctx context.Context, monitor entity.Apichecks) error {
	return r.db.With(ctx).Model(&monitor).Insert()
}

// Update saves the changes to an monitor in the database.
func (r repository) Update(ctx context.Context, monitor entity.Apichecks) error {
	return r.db.With(ctx).Model(&monitor).Update()
}

// Delete deletes an monitor with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	monitor, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&monitor).Delete()
}

// Count returns the number of the monitor records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("hawkeye.apichecks").Row(&count)
	return count, err
}

// Query retrieves the monitor records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Apichecks, error) {
	var monitors []entity.Apichecks
	err := r.db.With(ctx).
		Select("*").From("hawkeye.apichecks").
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&monitors)
	return monitors, err
}

// Query retrieves the monitor records with the specified offset and limit from the database.
func (r repository) GetByTenant(ctx context.Context, tenant string, offset, limit int) ([]entity.Apichecks, error) {
	var monitors []entity.Apichecks
	err := r.db.With(ctx).
		Select("*").From("hawkeye.apichecks").
		Where(dbx.HashExp{"tenant": tenant}).
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&monitors)
	return monitors, err
}

// Query retrieves the monitor records with the specified offset and limit from the database.
func (r repository) GetByOrg(ctx context.Context, org_id string, offset, limit int) ([]entity.Apichecks, error) {
	var monitors []entity.Apichecks
	err := r.db.With(ctx).
		Select("*").From("hawkeye.apichecks").
		Where(dbx.HashExp{"org_id": org_id}).
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&monitors)
	return monitors, err
}
