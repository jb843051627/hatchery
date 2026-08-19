package model

import "time"

const (
	AlertLevelInfo     = "info"
	AlertLevelWarning  = "warning"
	AlertLevelCritical = "critical"

	AlertStatusActive        = "active"
	AlertStatusAcknowledged  = "acknowledged"
	AlertStatusResolved      = "resolved"
)

type Alert struct {
	ID            int64     `json:"id"`
	IncubatorID   int64     `json:"incubator_id"`
	BatchID       int64     `json:"batch_id"`
	Level         string    `json:"level"`
	Message       string    `json:"message"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func ValidAlertLevel(s string) bool {
	switch s {
	case AlertLevelInfo, AlertLevelWarning, AlertLevelCritical:
		return true
	}
	return false
}
