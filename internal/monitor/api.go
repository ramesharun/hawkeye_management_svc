package monitor

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ramesharun/hawkeye-management-svc/internal/errors"
	"github.com/ramesharun/hawkeye-management-svc/pkg/log"
	"github.com/ramesharun/hawkeye-management-svc/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/monitors/<id>", res.get)
	r.Get("/monitors", res.query)
	r.Get("/monitorscount", res.count)
	r.Get("/orgmonitorscount/<org_id>", res.countbyorg)
	r.Get("/tenantmonitorscount/<tenant>", res.countbytenant)
	r.Get("/monitors/org/<org_id>", res.getbyorg)
	r.Get("/monitors/tenant/<tenant>", res.getbytenant)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/monitors", res.create)
	r.Put("/monitors/<id>", res.update)
	r.Delete("/monitors/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	monitor, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(monitor)
}

func (r resource) getbyorg(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	monitors, err := r.service.GetByOrg(ctx, c.Param("org_id"), pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = monitors
	return c.Write(pages)
}

func (r resource) getbytenant(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	monitors, err := r.service.GetByTenant(ctx, c.Param("tenant"), pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = monitors
	return c.Write(pages)
}

func (r resource) count(c *routing.Context) error {
	monitor, err := r.service.Count(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(monitor)
}

func (r resource) countbyorg(c *routing.Context) error {
	monitor, err := r.service.CountByOrg(c.Request.Context(), c.Param("org_id"))
	if err != nil {
		return err
	}

	return c.Write(monitor)
}

func (r resource) countbytenant(c *routing.Context) error {
	monitor, err := r.service.CountByTenant(c.Request.Context(), c.Param("tenant"))
	if err != nil {
		return err
	}

	return c.Write(monitor)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	monitors, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = monitors
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateMonitorRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	monitor, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(monitor, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateMonitorRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	monitor, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(monitor)
}

func (r resource) delete(c *routing.Context) error {
	monitor, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(monitor)
}
