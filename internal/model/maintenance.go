package model

import "time"

const (
	MaintenanceStatusPending   = "pending"
	MaintenanceStatusInProgress = "in_progress"
	MaintenanceStatusCompleted  = "completed"
)

type MaintenanceTask struct {
	ID           int64      `json:"id"`
	IncubatorID  int64      `json:"incubator_id"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	ScheduledFor time.Time  `json:"scheduled_for"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
