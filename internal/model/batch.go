package model

import "time"

const (
	BatchStatusPending    = "pending"
	BatchStatusIncubating = "incubating"
	BatchStatusHatching   = "hatching"
	BatchStatusCompleted  = "completed"
	BatchStatusFailed     = "failed"
)

type Batch struct {
	ID                int64     `json:"id"`
	IncubatorID       int64     `json:"incubator_id"`
	EggCount          int       `json:"egg_count"`
	Species           string    `json:"species"`
	StartDate         time.Time `json:"start_date"`
	ExpectedHatchDate time.Time `json:"expected_hatch_date"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

func ValidBatchStatus(s string) bool {
	switch s {
	case BatchStatusPending, BatchStatusIncubating, BatchStatusHatching,
		BatchStatusCompleted, BatchStatusFailed:
		return true
	}
	return false
}
