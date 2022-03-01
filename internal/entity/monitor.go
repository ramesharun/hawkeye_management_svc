package entity

import (
	"time"
)

// Apichecks represents the API Monitor registered
type Apichecks struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Monitor_ID string    `json:"monitor_id"`
	OrgID      string    `json:"org_id"`
	Tenant     string    `json:"tenant"`
	IsDeleted  bool      `json:"is_deleted"`
	UpdatedAt  time.Time `json:"last_ts"`
}

type TotalMonitorCount struct {
	Count string `json:"total_monitors_count"`
}

type TotalMonitorCountForTenant struct {
	TenantName string `json:"tenant_name"`
	Count      string `json:"total_monitors_count"`
}

type TotalMonitorCountForOrg struct {
	TenantName string `json:"org_name"`
	Count      string `json:"total_monitors_count"`
}
