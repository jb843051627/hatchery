package model

import "time"

const (
	IncubatorStatusActive     = "active"
	IncubatorStatusInactive   = "inactive"
	IncubatorStatusMaintenance = "maintenance"
)

type Incubator struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	InstalledAt time.Time `json:"installed_at"`
	CreatedAt   time.Time `json:"created_at"`
}
