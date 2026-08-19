package model

import "time"

const (
	SchedulePhasePreheat  = "preheat"
	SchedulePhaseIncubate = "incubate"
	SchedulePhaseHatch    = "hatch"
	SchedulePhaseCleanup  = "cleanup"

	ScheduleStatusPending   = "pending"
	ScheduleStatusActive    = "active"
	ScheduleStatusCompleted = "completed"
	ScheduleStatusCancelled = "cancelled"
)

type Schedule struct {
	ID             int64     `json:"id"`
	IncubatorID    int64     `json:"incubator_id"`
	BatchID        int64     `json:"batch_id"`
	Phase          string    `json:"phase"`
	TargetTemp     float64   `json:"target_temp"`
	TargetHumidity float64   `json:"target_humidity"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func ValidSchedulePhase(s string) bool {
	switch s {
	case SchedulePhasePreheat, SchedulePhaseIncubate, SchedulePhaseHatch, SchedulePhaseCleanup:
		return true
	}
	return false
}

func ValidScheduleStatus(s string) bool {
	switch s {
	case ScheduleStatusPending, ScheduleStatusActive, ScheduleStatusCompleted, ScheduleStatusCancelled:
		return true
	}
	return false
}
